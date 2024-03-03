package main

import (
	"sync"
)

type StorageEngine interface {
	Set(shard int, key string, value string)
	Get(shard int, key string) (string, bool)
	GetShard(shard int) *sync.Map
	DeleteShard(shard int)
	GetShards() map[int]*sync.Map
}

type MemStorageEngine struct {
	Shards map[int]*sync.Map
}

func NewMemStorageEngine() StorageEngine {
	return &MemStorageEngine{
		Shards: make(map[int]*sync.Map),
	}
}

func (s *MemStorageEngine) Set(shard int, key string, value string) {
	if _, ok := s.Shards[shard]; !ok {
		s.Shards[shard] = &sync.Map{}
	}
	s.Shards[shard].Store(key, value)
}

func (s *MemStorageEngine) Get(shard int, key string) (string, bool) {
	if _, ok := s.Shards[shard]; !ok {
		return "", false
	}
	value, ok := s.Shards[shard].Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (s *MemStorageEngine) GetShard(shard int) *sync.Map {
	return s.Shards[shard]
}

func (s *MemStorageEngine) DeleteShard(shard int) {
	delete(s.Shards, shard)
}

func (s *MemStorageEngine) GetShards() map[int]*sync.Map {
	return s.Shards
}
