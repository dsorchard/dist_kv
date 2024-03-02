package main

import (
	"github.com/buraksezer/consistent"
	"github.com/hashicorp/memberlist"
)

type Node struct {
	Addr string
}

func NewNode(addr string) *Node {
	return &Node{
		Addr: addr,
	}
}

type Cluster struct {
	*memberlist.Memberlist
	LocalNode *Node
	store     *KeyValueStore
	Ring      *consistent.Consistent
}

func NewCluster(localNode *Node, store *KeyValueStore) (*Cluster, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = localNode.Addr
	config.BindAddr = localNode.Addr

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	return &Cluster{
		Memberlist: list,
		LocalNode:  localNode,
		store:      store,
	}, nil
}

func (c *Cluster) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)
	return err
}

func (c *Cluster) NotifyMsg(msg []byte) {
	// Handle incoming messages for data replication
}
