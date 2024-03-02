package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
	distKV *DistKV
}

func NewAPI(distKV *DistKV) *API {
	api := &API{
		router: mux.NewRouter(),
		distKV: distKV,
	}
	api.setupRoutes()
	return api
}

func (api *API) setupRoutes() {
	api.router.HandleFunc("/put/{key}/{value}", api.setHandler).Methods("POST")
	api.router.HandleFunc("/get/{key}", api.getHandler).Methods("GET")
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {

	key := mux.Vars(r)["key"]
	routeNodeAddress := api.distKV.ring.GetNode(key)

	localNodeAddress := api.distKV.node.NodeHttpAddress()
	if routeNodeAddress == localNodeAddress {
		value, ok := api.distKV.kv.Get(key)
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(value))
	} else {
		url := fmt.Sprintf("http://%s/get/%s", routeNodeAddress, key)
		resp, err := http.Get(url)
		if err != nil {
			// Log the error and return a server error response
			log.Printf("Error forwarding request to %s: %v", routeNodeAddress, err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Check if the responsible node found the key
		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		} else if resp.StatusCode != http.StatusOK {
			// Handle unexpected status code
			http.Error(w, "Error retrieving key from responsible node", http.StatusInternalServerError)
			return
		}

		// Forward the response body (value) to the client
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Error reading key value", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(body)
	}
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]

	routeNodeAddress := api.distKV.ring.GetNode(key)
	localNodeAddress := api.distKV.node.NodeHttpAddress()
	if routeNodeAddress == localNodeAddress {
		api.distKV.kv.Set(key, value)
		w.WriteHeader(http.StatusOK)
	} else {
		url := fmt.Sprintf("http://%s/put/%s/%s", routeNodeAddress, key, value)
		resp, err := http.Post(url, "", nil)
		if err != nil {
			// Log the error and return a server error response
			log.Printf("Error forwarding request to %s: %v", routeNodeAddress, err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Check if the responsible node found the key
		if resp.StatusCode != http.StatusOK {
			// Handle unexpected status code
			http.Error(w, "Error setting key on responsible node", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (api *API) Run(addr string) {
	_ = http.ListenAndServe(addr, api.router)
}
