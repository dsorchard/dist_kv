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
	go distKV.HandleMembershipChange(membershipChangeCh)

	return &distKV
}

func (d *DistKV) Bootstrap() {
	err := d.node.Join(d.config.BootstrapNodes)
	if err != nil {
		log.Fatalf("Failed to join cluster: %v", err)
	}

	api := NewAPI(d.node)
	api.Run(fmt.Sprintf(":%d", d.config.ExternalPort))
}

func (d *DistKV) HandleMembershipChange(membershipChangeCh chan memberlist.NodeEvent) {
	for {
		select {
		case event := <-membershipChangeCh:
			switch event.Event {
			case memberlist.NodeJoin:
				d.ring.AddNode(event.Node.Name)
				log.Printf("Node joined: %s", event.Node.Name)
			case memberlist.NodeLeave:
				d.ring.RemoveNode(event.Node.Name)
				log.Printf("Node left: %s", event.Node.Name)
			default:
				log.Fatalf("Unknown event: %v", event.Event)
			}
		}
	}
}
