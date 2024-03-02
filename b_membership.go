package main

import (
	"github.com/hashicorp/memberlist"
	"strconv"
)

type Node struct {
	*memberlist.Memberlist
}

func NewNode(gossipPort int) (*Node, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = GetLocalIP() + ":" + strconv.Itoa(gossipPort)
	config.BindPort = gossipPort

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	return &Node{
		Memberlist: list,
	}, nil
}

func (c *Node) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)
	return err
}

func (c *Node) NotifyMsg(msg []byte) {
	// Handle incoming messages for data replication
}
