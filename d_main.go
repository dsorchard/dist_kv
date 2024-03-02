package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	config := loadConfig()

	var externalPort, internalPort int
	flag.IntVar(&externalPort, "export", config.ExternalPort, "port number for external request")
	flag.IntVar(&internalPort, "inport", config.InternalPort, "port number for internal protocol communication")
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

	api := NewAPI(cluster)
	api.Run(fmt.Sprintf(":%d", config.ExternalPort))
}
