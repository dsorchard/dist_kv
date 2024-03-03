package main

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

type HashRing interface {
	AddNode(node string)
	RemoveNode(node string)
	ResolveNode(key string) string
	ResolveNodes(key string, count int) []string
	ResolvePartitionID(key string) int
	ResolvePartitionOwnerNode(partitionID int) string
	ReplicationFactor() int
	ResolveNodesForPartition(partitionID int, count int) []string
}

type BoundedLoadConsistentHashRing struct {
	ring              *consistent.Consistent
	replicationFactor int
}

func NewBoundedLoadConsistentHashRing(virtualNodeCount, replicationFactorForEachKey int) HashRing {
	cfg := consistent.Config{
		PartitionCount:    virtualNodeCount,            // virtual node count
		ReplicationFactor: replicationFactorForEachKey, // number of replicas for each key
		Load:              1.25,                        // server can be 25% > average before the system redistributes
		Hasher:            hasher{},
	}
	return &BoundedLoadConsistentHashRing{
		ring:              consistent.New(nil, cfg),
		replicationFactor: replicationFactorForEachKey,
	}
}

func (r *BoundedLoadConsistentHashRing) AddNode(node string) {
	r.ring.Add(member(node))
}

func (r *BoundedLoadConsistentHashRing) RemoveNode(node string) {
	r.ring.Remove(node)
}

func (r *BoundedLoadConsistentHashRing) ResolveNode(key string) string {
	return r.ring.LocateKey([]byte(key)).String()
}

// ResolveNodes returns the closest N nodes to the key in the ring.
// TODO: Will be used for replication.
func (r *BoundedLoadConsistentHashRing) ResolveNodes(key string, count int) []string {
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

func (r *BoundedLoadConsistentHashRing) ResolvePartitionID(key string) int {
	return r.ring.FindPartitionID([]byte(key))
}

func (r *BoundedLoadConsistentHashRing) ResolvePartitionOwnerNode(partitionID int) string {
	return r.ring.GetPartitionOwner(partitionID).String()
}

func (r *BoundedLoadConsistentHashRing) ReplicationFactor() int {
	return r.replicationFactor
}

func (r *BoundedLoadConsistentHashRing) ResolveNodesForPartition(partitionID int, count int) []string {
	members, err := r.ring.GetClosestNForPartition(partitionID, count)
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
