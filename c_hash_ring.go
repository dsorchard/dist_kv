package main

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

type HashRing struct {
	ring *consistent.Consistent
}

func NewRing() *HashRing {
	cfg := consistent.Config{
		PartitionCount:    7,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}
	return &HashRing{
		ring: consistent.New(nil, cfg),
	}
}

func (r *HashRing) AddNode(node string) {
	r.ring.Add(member(node))
}

func (r *HashRing) RemoveNode(node string) {
	r.ring.Remove(node)
}

func (r *HashRing) GetNode(key string) string {
	return r.ring.LocateKey([]byte(key)).String()
}

func (r *HashRing) GetNodes(key string, count int) []string {
	members, err := r.ring.GetClosestN([]byte(key), count)
	if err != nil {
		return nil
	}
	nodes := make([]string, len(members))
	for i, m := range members {
		nodes[i] = m.String()
	}
	return nodes
}

//------------------------ Sub Classes ---------------------------------

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

type member string

func (m member) String() string {
	return string(m)
}
