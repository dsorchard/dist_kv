package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"strconv"
)

type GossipNode struct {
	*memberlist.Memberlist
}

func NewNode(gossipPort int) (*GossipNode, chan memberlist.NodeEvent, error) {
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

	return &GossipNode{Memberlist: list}, membershipChangeCh, nil
}

func (c *GossipNode) Join(existing []string) error {
	_, err := c.Memberlist.Join(existing)
	return err
}

func (c *GossipNode) NodeHttpAddress() string {
	return fmt.Sprintf("%s:%d", GetLocalIP(), c.Memberlist.LocalNode().Port+1)
}
