package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
)

type DistKVServer struct {
	config *configuration
	kv     *KVStore
	ring   *HashRing
	node   *GossipNode
}

func NewDistKVServer(config *configuration) *DistKVServer {
	node, membershipChangeCh, err := NewNode(config.InternalPort)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	kv := NewKVStore()

	ring := NewRing(config.PartitionCount, config.KeyReplicationCount)

	distKV := DistKVServer{
		config: config,
		node:   node,
		kv:     kv,
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

	api := NewAPI(d)
	api.Run(fmt.Sprintf(":%d", d.config.ExternalPort))
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
	for partitionId, partitionContent := range d.kv.Data {
		newOwner := d.ring.GetPartitionOwner(partitionId)
		if newOwner != d.config.Host {
			delete(d.kv.Data, partitionId)

			// send to newOwner
			//TODO: do later
			log.Printf("Redistributing partition %d to %s", partitionId, newOwner)
			log.Printf("Partition content: %v", partitionContent)
		}
	}
}
