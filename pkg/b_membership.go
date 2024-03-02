package main

import (
	"github.com/hashicorp/memberlist"
	"strconv"
)

type Node struct {
	*memberlist.Memberlist
}

func NewNode(gossipPort int) (*Node, chan memberlist.NodeEvent, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = GetLocalIP() + ":" + strconv.Itoa(gossipPort)
	config.BindPort = gossipPort
	membershipChangeCh := make(chan memberlist.NodeEvent, 16)
	config.Events = &memberlist.ChannelEventDelegate{
		Ch: membershipChangeCh,
	}

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, nil, err
	}

	return &Node{Memberlist: list}, membershipChangeCh, nil
}

func (c *Node) Join(seeds []string) error {
	_, err := c.Memberlist.Join(seeds)
	return err
}

func (c *Node) NodeAddress() string {
	return c.Memberlist.LocalNode().Address()
}
