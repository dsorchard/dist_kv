package main

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/uber/ringpop-go"
	"github.com/uber/ringpop-go/discovery/statichosts"
	"github.com/uber/ringpop-go/swim"
	"github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/json"
)

type DistKVServer struct {
	rp         *ringpop.Ringpop
	store      *StorageEngine
	ch         *tchannel.Channel
	httpServer *HttpAPIServer
}

func NewDistKVServer(hostPort, serviceName, cluster string, httpPort int) *DistKVServer {
	ch, err := tchannel.NewChannel(serviceName, nil)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	rp, err := ringpop.New("dist-kv",
		ringpop.Channel(ch),
		ringpop.Identity(hostPort),
	)
	if err != nil {
		log.Fatalf("failed to create ringpop: %v", err)
	}

	bootstrapOpts := &swim.BootstrapOptions{
		DiscoverProvider: statichosts.New(cluster),
	}
	if _, err := rp.Bootstrap(bootstrapOpts); err != nil {
		log.Fatalf("failed to bootstrap: %v", err)
	}

	app := &DistKVServer{
		rp:    rp,
		store: NewStorageEngine(),
		ch:    ch,
	}

	httpServer := NewAPI(app, httpPort)
	app.httpServer = httpServer

	return app
}

func (app *DistKVServer) setupHandlers() {
	err := json.Register(app.ch, json.Handlers{
		"set": app.handleSet,
		"get": app.handleGet,
	}, func(ctx context.Context, err error) {
		log.Fatal("error handling request:", err)
	})
	if err != nil {
		return
	}

	app.httpServer.Run()
}

func (app *DistKVServer) handleSet(ctx json.Context, req *SetRequest) (res *SetResponse, err error) {
	app.store.Set(req.Key, req.Value)
	res.Status = "success"
	return res, nil
}

func (app *DistKVServer) handleGet(ctx json.Context, req *GetRequest) (res *GetResponse, err error) {
	value, exists := app.store.Get(req.Key)
	if !exists {
		res.Status = "key not found"
		return res, nil
	}
	res.Status = "success"
	res.Value = value
	return res, nil
}

func (app *DistKVServer) Set(key, value string) {
	dest, err := app.rp.Lookup(key)
	if err != nil {
		panic(err)
	}

	local, err := app.rp.WhoAmI()
	if err != nil {
		panic(err)
	}

	if dest == local {
		app.store.Set(key, value)
		return
	}

	req := &SetRequest{Key: key, Value: value}
	res := &SetResponse{}

	client := json.NewClient(app.ch, dest, nil)
	err = client.Call(nil, "set", req, res)
	if err != nil {
		panic(err)
	}

	if res.Status != "success" {
		panic("failed to set key on remote node")
	}

	return
}

func (app *DistKVServer) Get(key string) (string, bool) {
	dest, err := app.rp.Lookup(key)
	if err != nil {
		panic(err)
	}

	local, err := app.rp.WhoAmI()
	if err != nil {
		panic(err)
	}

	if dest == local {
		return app.store.Get(key)
	}

	req := &GetRequest{Key: key}
	res := &GetResponse{}

	client := json.NewClient(app.ch, dest, nil)
	err = client.Call(nil, "get", req, res)
	if err != nil {
		panic(err)
	}

	if res.Status != "success" {
		panic("failed to get key from remote node")
	}

	return res.Value, true
}

//----------------------Req/Res--------------------------------------

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResponse struct {
	Status string `json:"status"`
}

type GetRequest struct {
	Key string `json:"key"`
}

type GetResponse struct {
	Status string `json:"status"`
	Value  string `json:"value,omitempty"`
}
