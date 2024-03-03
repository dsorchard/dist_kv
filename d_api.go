package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	router   *mux.Router
	distKV   *DistKVServer
	httpPort int
}

func NewAPI(distKV *DistKVServer, httpPort int) *API {
	api := &API{
		router:   mux.NewRouter(),
		distKV:   distKV,
		httpPort: httpPort,
	}
	api.router.HandleFunc("/put/{key}/{value}", api.setHandler).Methods("POST")
	api.router.HandleFunc("/get/{key}", api.getHandler).Methods("GET")
	return api
}

func (api *API) Run() {
	addr := fmt.Sprintf(":%d", api.httpPort)
	_ = http.ListenAndServe(addr, api.router)
}

func (api *API) GetAddress() string {
	return fmt.Sprintf("%s:%d", GetLocalIP(), api.httpPort)
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)["key"]
	routeNodeAddress := api.distKV.ring.ResolveNode(key)

	localNodeAddress := api.GetAddress()
	if routeNodeAddress == localNodeAddress {
		shardId := api.distKV.ring.ResolvePartitionID(key)
		value, ok := api.distKV.kv.Get(shardId, key)
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(value))
	} else {
		client := NewClient(routeNodeAddress)
		value, err := client.Get(key)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", routeNodeAddress, err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(value))
		return
	}
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]

	routeNodeAddress := api.distKV.ring.ResolveNode(key)
	localNodeAddress := api.GetAddress()
	if routeNodeAddress == localNodeAddress {
		shardId := api.distKV.ring.ResolvePartitionID(key)
		api.distKV.kv.Set(shardId, key, value)
	} else {
		client := NewClient(routeNodeAddress)
		err := client.Put(key, value)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", routeNodeAddress, err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
