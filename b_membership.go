package main

import (
	"github.com/charmbracelet/log"
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
	config.Name = GetLocalIP() + "@port:" + strconv.Itoa(gossipPort)
	config.BindAddr = GetLocalIP()
	config.BindPort = gossipPort
	config.LogOutput = NewMemberlistLogger()

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

// -----------------------Logger----------------------------------

type MemberlistLogger struct {
	Logger *log.Logger
}

func NewMemberlistLogger() MemberlistLogger {
	return MemberlistLogger{
		Logger: log.WithPrefix("memberlist"),
	}
}

func (l MemberlistLogger) Write(p []byte) (n int, err error) {
	l.Logger.Debugf(string(p))
	return len(p), nil
}
