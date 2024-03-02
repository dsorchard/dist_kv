package main

import (
	"fmt"
	"log"
)

type DistKV struct {
	config *configuration
	node   *Node
	kv     *KVStore
	ring   *HashRing
}

func NewDistKV(config *configuration) *DistKV {
	node, err := NewNode(config.InternalPort)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	kv := NewKeyValueStore()

	ring := NewRing()

	return &DistKV{
		node: node,
		kv:   kv,
		ring: ring,
	}
}

func (d *DistKV) Bootstrap() {
	err := d.node.Join(d.config.BootstrapNodes)
	if err != nil {
		log.Fatalf("Failed to join cluster: %v", err)
	}

	api := NewAPI(d.node)
	api.Run(fmt.Sprintf(":%d", d.config.ExternalPort))
}
