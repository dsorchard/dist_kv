package main

import (
	"github.com/buraksezer/consistent"
	"github.com/hashicorp/memberlist"
	"strconv"
)

type Cluster struct {
	*memberlist.Memberlist
	store *KeyValueStore
	Ring  *consistent.Consistent
}

func NewCluster(gossipPort int, store *KeyValueStore) (*Cluster, error) {
	config := memberlist.DefaultLocalConfig()
	config.Name = GetLocalIP() + ":" + strconv.Itoa(gossipPort)
	config.BindPort = gossipPort

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	return &Cluster{
		Memberlist: list,
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
