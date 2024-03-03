package main

import (
	"github.com/hashicorp/memberlist"
	"strconv"
)

type Membership interface {
	Join(existing []string) error
	MembershipChangeCh() chan memberlist.NodeEvent
}

// GossipMembership is a membership implementation using hashicorp/memberlist
// It could be EtcdMembership as well as done in JunoDB
type GossipMembership struct {
	membershipList     *memberlist.Memberlist
	hostName           string
	gossipPort         int
	membershipChangeCh chan memberlist.NodeEvent
}

func NewGossipMembership(gossipPort int) (Membership, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = GetLocalIP() + ":" + strconv.Itoa(gossipPort)
	config.BindPort = gossipPort
	membershipChangeCh := make(chan memberlist.NodeEvent, 16)
	config.Events = &memberlist.ChannelEventDelegate{
		Ch: membershipChangeCh,
	}

	membershipList, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	return &GossipMembership{
		membershipList:     membershipList,
		membershipChangeCh: membershipChangeCh,
	}, nil
}

func (c *GossipMembership) Join(existing []string) error {
	_, err := c.membershipList.Join(existing)
	return err
}

func (c *GossipMembership) MembershipChangeCh() chan memberlist.NodeEvent {
	return c.membershipChangeCh
}
