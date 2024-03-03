package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type DistKVServer struct {
	config *configuration
	store  StorageEngine
	ring   HashRing
	node   Membership
}

func NewDistKVServer(config *configuration) *DistKVServer {
	node, err := NewGossipMembership(config.InternalPort)
	membershipChangeCh := node.MembershipChangeCh()
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	kv := NewMemStorageEngine()

	ring := NewBoundedLoadConsistentHashRing(config.PartitionCount, config.KeyReplicationCount)

	distKV := DistKVServer{
		config: config,
		node:   node,
		store:  kv,
		ring:   ring,
	}
	go distKV.handleMembershipChange(membershipChangeCh)

	return &distKV
}

func (d *DistKVServer) Bootstrap() {
	err := d.node.Join(d.config.BootstrapNodes)
	if err != nil {
		log.Fatalf("Failed to join distKV: %v", err)
	}

	api := NewAPI(d, d.config.ExternalPort)
	api.Run()
}

func (d *DistKVServer) handleMembershipChange(membershipChangeCh chan memberlist.NodeEvent) {
	for {
		select {
		case event := <-membershipChangeCh:
			httpAddress := fmt.Sprintf("%s:%d", GetLocalIP(), event.Node.Port+1)
			switch event.Event {
			case memberlist.NodeJoin:
				d.ring.AddNode(httpAddress)
				d.redistributePartitions(httpAddress)
				log.Printf("Node joined: %s", httpAddress)
			case memberlist.NodeLeave:
				d.ring.RemoveNode(httpAddress)
				log.Printf("Node left: %s", httpAddress)
			default:
				log.Fatalf("Unknown event: %v", event.Event)
			}
		}
	}
}

func (d *DistKVServer) redistributePartitions(address string) {
	for partitionId, partitionContent := range d.store.GetShards() {
		newOwner := d.ring.ResolvePartitionOwnerNode(partitionId)
		if newOwner != d.config.Host {
			// send to newOwner
			//TODO: do later
			log.Printf("Redistributing partition %d to %s", partitionId, newOwner)
			log.Printf("Partition content: %v", partitionContent)

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
	log.Printf("Loading configurations from config.yml")

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

func GetLocalIP() string {
	return localIp
}
