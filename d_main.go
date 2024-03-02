package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	config := loadConfig()

	var externalPort, internalPort int
	flag.IntVar(&externalPort, "http", config.ExternalPort, "port number for external request")
	flag.IntVar(&internalPort, "gossip", config.InternalPort, "port number for internal protocol communication")
	flag.Parse()

	config.Host = GetLocalIP()
	config.ExternalPort = externalPort
	config.InternalPort = internalPort

	cluster, err := NewCluster(
		NewNode(fmt.Sprintf("%s:%d", config.Host, config.InternalPort)),
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
	api.Run(fmt.Sprintf(":%d", config.ExternalPort))
}
