package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type configuration struct {
	Host           string
	ExternalPort   int      `yaml:"ExternalPort"`
	InternalPort   int      `yaml:"InternalPort"`
	BootstrapNodes []string `yaml:"BootstrapNodes"`
	//VirtualNodeSize  int `yaml:"VirtualNodeSize"`
	//KVSReplicaPoints int `yaml:"KVSReplicaPoints"`
}

func loadConfig() *configuration {
	log.Printf("Loading configurations from config.yml")

	config := &configuration{
		Host:           "0.0.0.0",
		InternalPort:   8000,
		ExternalPort:   8001,
		BootstrapNodes: []string{},
		//VirtualNodeSize:    5,
		//KVSReplicaPoints:   3,
	}

	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Cannot load config.yml")
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("Fail to unmarshal config.yml")
	}

	for i, addr := range config.BootstrapNodes {
		if strings.HasPrefix(addr, ":") {
			config.BootstrapNodes[i] = GetLocalIP() + addr
		}
	}

	return config
}
