package main

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

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
