package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var httpLogger = log.WithPrefix("http")

type HttpAPIServer struct {
	router   *mux.Router
	distKV   *DistKVServer
	httpPort int
}

func NewAPI(distKV *DistKVServer, httpPort int) *HttpAPIServer {
	api := &HttpAPIServer{
		router:   mux.NewRouter(),
		distKV:   distKV,
		httpPort: httpPort,
	}
	api.router.HandleFunc("/kv/{key}/{value}", api.setKvHandler).Methods("POST")
	api.router.HandleFunc("/kv/{key}", api.getKvHandler).Methods("GET")

	api.router.HandleFunc("/store/{key}/{value}", api.setStoreHandler).Methods("POST")
	api.router.HandleFunc("/store/{key}", api.getStoreHandler).Methods("GET")

	api.router.HandleFunc("/shards/{key}", api.setShardHandler).Methods("POST")
	api.router.HandleFunc("/shards", api.getShardHandler).Methods("GET")
	return api
}

func (api *HttpAPIServer) Run() {
	addr := fmt.Sprintf(":%d", api.httpPort)
	_ = http.ListenAndServe(addr, api.router)
}

func (api *HttpAPIServer) GetAddress() string {
	return fmt.Sprintf("%s:%d", GetLocalIP(), api.httpPort)
}

func (api *HttpAPIServer) getKvHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)["key"]
	routeNodeAddress := api.distKV.ring.ResolveNode(key)

	value, err := api.distKV.client.Get(key)
	if err != nil {
		httpLogger.Errorf("Error forwarding request to %s: %v", routeNodeAddress, err)
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(value))
	return
}

func (api *HttpAPIServer) setKvHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]

	routeNodeAddress := api.distKV.ring.ResolveNode(key)
	err := api.distKV.client.Put(key, value)
	if err != nil {
		httpLogger.Errorf("Error forwarding request to %s: %v", routeNodeAddress, err)
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *HttpAPIServer) getShardHandler(w http.ResponseWriter, r *http.Request) {

	shards := api.distKV.store.GetShards()
	for shardId, shard := range shards {
		httpLogger.Warnf("Shard %d", shardId)
		shard.Range(func(key, value interface{}) bool {
			httpLogger.Warnf("Key: %s, Value: %s", key, value)
			return true
		})
	}
	w.WriteHeader(http.StatusOK)
}

func (api *HttpAPIServer) setShardHandler(w http.ResponseWriter, r *http.Request) {
	shardIdStr := mux.Vars(r)["key"]
	shardId, _ := strconv.Atoi(shardIdStr)

	var shard map[string]string
	if err := json.NewDecoder(r.Body).Decode(&shard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	api.distKV.store.SetShard(shardId, shard)
	w.WriteHeader(http.StatusOK)
}

func (api *HttpAPIServer) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	shardId := api.distKV.ring.ResolvePartitionID(key)
	value, ok := api.distKV.store.Get(shardId, key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte(value))
}

func (api *HttpAPIServer) setStoreHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]
	shardId := api.distKV.ring.ResolvePartitionID(key)
	api.distKV.store.Set(shardId, key, value)
	w.WriteHeader(http.StatusOK)
}
