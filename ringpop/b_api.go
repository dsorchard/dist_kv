package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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

	return api
}

func (api *HttpAPIServer) Run() {
	addr := fmt.Sprintf(":%d", api.httpPort)
	_ = http.ListenAndServe(addr, api.router)
}

func (api *HttpAPIServer) getKvHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)["key"]
	value, ok := api.distKV.Get(key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte(value))
	return
}

func (api *HttpAPIServer) setKvHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]
	api.distKV.Set(key, value)
	w.WriteHeader(http.StatusOK)
}
