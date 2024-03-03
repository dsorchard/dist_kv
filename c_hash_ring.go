package main

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

type HashRing struct {
	ring *consistent.Consistent
}

/*
- PartitionCount: This should be significantly higher than the number of servers
to ensure a fine-grained distribution. A common approach is to use a multiple
of the number of servers. For example, with 10 servers, you might use 100 or 200
partitions, depending on your specific requirements for granularity and the overhead of managing more partitions.

- ReplicationFactor: Set this to 3, as you want each key to be replicated
three times and stored on different servers.

- Load: The ideal value for Load depends on your system's characteristics,
but a starting point might be 1.25 or 1.5. This means a server can be
25% or 50% more loaded than the average before the system redistributes
keys to balance the load.
*/

func NewRing(partitionCount, replicationFactor int) *HashRing {
	cfg := consistent.Config{
		PartitionCount:    partitionCount, // micro shards
		ReplicationFactor: replicationFactor,
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
