package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
)

type DistKV struct {
	config *configuration
	node   *Node
	kv     *KVStore
	ring   *HashRing
}

func NewDistKV(config *configuration) *DistKV {
	node, membershipChangeCh, err := NewNode(config.InternalPort)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	kv := NewKVStore()

	ring := NewRing()

	distKV := DistKV{
		config: config,
		node:   node,
		kv:     kv,
		ring:   ring,
	}
	go distKV.handleMembershipChange(membershipChangeCh)

	return &distKV
}

func (d *DistKV) Bootstrap() {
	err := d.node.Join(d.config.BootstrapNodes)
	if err != nil {
		log.Fatalf("Failed to join distKV: %v", err)
	}

	api := NewAPI(d)
	api.Run(fmt.Sprintf(":%d", d.config.ExternalPort))
}

func (d *DistKV) handleMembershipChange(membershipChangeCh chan memberlist.NodeEvent) {
	for {
		select {
		case event := <-membershipChangeCh:
			httpAddress := fmt.Sprintf("%s:%d", GetLocalIP(), event.Node.Port+1)
			switch event.Event {
			case memberlist.NodeJoin:
				d.ring.AddNode(httpAddress)
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
