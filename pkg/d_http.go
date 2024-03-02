package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	router  *mux.Router
	cluster *Node
}

func NewAPI(cluster *Node) *API {
	api := &API{
		router:  mux.NewRouter(),
		cluster: cluster,
	}
	api.setupRoutes()
	return api
}

func (api *API) setupRoutes() {
	api.router.HandleFunc("/set/{key}/{value}", api.setHandler).Methods("POST")
	api.router.HandleFunc("/get/{key}", api.getHandler).Methods("GET")
	api.router.HandleFunc("/delete/{key}", api.deleteHandler).Methods("DELETE")
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {
	//replicas := api.cluster.ring.GetClosestN(r.Key, rc.sr.config.KVSReplicaPoints)
	//resCh := rc.sendRPCRequests(replicas, GetOp, internalReq)

	// Handle get request
}

func (api *API) deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Handle delete request and replicate data
}

func (api *API) Run(addr string) {
	http.ListenAndServe(addr, api.router)
}
