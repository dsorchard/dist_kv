package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/hashicorp/memberlist"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

var distKvLogger = log.WithPrefix("dist-kv")
var configLogger = log.WithPrefix("config")

type DistKVServer struct {
	config *configuration
	store  StorageEngine
	ring   HashRing
	node   Membership
	api    *HttpAPIServer
	client Client
}

func NewDistKVServer(config *configuration) *DistKVServer {
	httpAddress := fmt.Sprintf("%s:%d", GetLocalIP(), config.ExternalPort)
	node, err := NewGossipMembership(config.InternalPort, httpAddress)
	if err != nil {
		distKvLogger.Fatalf("Failed to create node: %v", err)
	}

	kv := NewMemStorageEngine()

	ring := NewBoundedLoadConsistentHashRing(config.PartitionCount, config.KeyReplicationCount)

	distKV := DistKVServer{
		config: config,
		node:   node,
		store:  kv,
		ring:   ring,
	}
	go distKV.handleMembershipChange(node.MembershipChangeCh())

	api := NewAPI(&distKV, config.ExternalPort)
	distKV.api = api

	distKV.client = NewHttpClient(ring, api)

	return &distKV
}

func (d *DistKVServer) Bootstrap() {
	// Join the gossip network
	err := d.node.Join(d.config.BootstrapNodes)
	if err != nil {
		distKvLogger.Fatalf("Failed to join distKV: %v", err)
	}

	// Start the HTTP API server
	d.api.Run()
}

func (d *DistKVServer) handleMembershipChange(membershipChangeCh chan memberlist.NodeEvent) {
	for {
		select {
		case event := <-membershipChangeCh:
			httpAddress := event.Node.Name
			switch event.Event {
			case memberlist.NodeJoin:
				d.ring.AddNode(httpAddress)
				d.redistributePartitions()
				distKvLogger.Infof("Node joined: %s", httpAddress)
			case memberlist.NodeLeave:
				d.ring.RemoveNode(httpAddress)
				distKvLogger.Infof("Node left: %s", httpAddress)
			default:
				distKvLogger.Fatalf("Unknown event: %v", event.Event)
			}
		}
	}
}

func (d *DistKVServer) redistributePartitions() {
	for partitionId := range d.store.GetShards() {
		newOwner := d.ring.ResolvePartitionOwnerNode(partitionId)
		if newOwner != d.api.GetAddress() {
			// send to newOwner
			err := d.client.PutShard(partitionId, d.store.GetShard(partitionId))
			if err != nil {
				distKvLogger.Fatalf("Failed to redistribute partition %d to %s: %v", partitionId, newOwner, err)
				continue
			}
			distKvLogger.Infof("Redistributing partition %d to %s", partitionId, newOwner)

			// delete from old owner
			d.store.DeleteShard(partitionId)
		}
	}
}

// -----------------Config -------------------
const (
	localIp = "127.0.0.1"
)

type configuration struct {
	Host                string
	ExternalPort        int      `yaml:"ExternalPort"`
	InternalPort        int      `yaml:"InternalPort"`
	BootstrapNodes      []string `yaml:"BootstrapNodes"`
	PartitionCount      int      `yaml:"PartitionCount"`
	KeyReplicationCount int      `yaml:"KeyReplicationCount"`
}

func loadConfig() *configuration {
	configLogger.Info("loading configurations from config.yml")

	config := &configuration{
		Host:                "0.0.0.0",
		InternalPort:        8000,
		ExternalPort:        8001,
		BootstrapNodes:      []string{},
		PartitionCount:      30,
		KeyReplicationCount: 3,
	}

	data, err := os.ReadFile("z_config.yml")
	if err != nil {
		configLogger.Fatal("Cannot load config.yml")
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		configLogger.Fatal("Failed to unmarshal config.yml")
	}

	for i, addr := range config.BootstrapNodes {
		if strings.HasPrefix(addr, ":") {
			config.BootstrapNodes[i] = GetLocalIP() + addr
		}
	}

	return config
}

func GetLocalIP() string {
	return localIp
}
