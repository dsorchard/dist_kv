package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	config := loadConfig()

	var httpPort, gossipPort int
	flag.IntVar(&gossipPort, "gossip", config.InternalPort, "port number for gossip protocol")
	flag.IntVar(&httpPort, "http", config.ExternalPort, "port number for http server")
	flag.Parse()

	config.Host = GetLocalIP()
	config.InternalPort = gossipPort
	config.ExternalPort = httpPort

	cluster, err := NewCluster(
		config.InternalPort,
		NewKeyValueStore(),
	)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	err = cluster.Join(config.BootstrapNodes)
	if err != nil {
		log.Fatalf("Failed to join cluster: %v", err)
	}

	api := NewAPI(cluster)
	api.Run(fmt.Sprintf(":%d", httpPort))
}
