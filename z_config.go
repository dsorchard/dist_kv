package main

import (
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

var logger = logging.MustGetLogger("swimring")

type configuration struct {
	Host         string
	ExternalPort int `yaml:"ExternalPort"`
	InternalPort int `yaml:"InternalPort"`

	//VirtualNodeSize  int `yaml:"VirtualNodeSize"`
	//KVSReplicaPoints int `yaml:"KVSReplicaPoints"`
	BootstrapNodes []string `yaml:"BootstrapNodes"`
}

func loadConfig() *configuration {
	logger.Info("Loading configurations from config.yml")

	config := &configuration{
		Host:         "0.0.0.0",
		ExternalPort: 7000,
		InternalPort: 7001,
		//VirtualNodeSize:    5,
		//KVSReplicaPoints:   3,
		BootstrapNodes: []string{},
	}

	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logger.Warning("Cannot load config.yml")
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		logger.Error("Fail to unmarshal config.yml")
	}

	for i, addr := range config.BootstrapNodes {
		if strings.HasPrefix(addr, ":") {
			config.BootstrapNodes[i] = GetLocalIP() + addr
		}
	}

	return config
}
