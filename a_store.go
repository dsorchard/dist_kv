package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type StorageEngine interface {
	Set(shard int, key string, value string)
	Get(shard int, key string) (string, bool)
	GetShards() map[int]*sync.Map
	GetShard(shard int) map[string]string
	DeleteShard(shard int)
	SetShard(shardId int, shard map[string]string)
}

type MemStorageEngine struct {
	Shards map[int]*sync.Map // Each Shard is a Virtual Node.
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
	//TODO: remove this Multi Versioning Hack
	if !strings.Contains(value, ":") {
		value = fmt.Sprintf("%d:%s", time.Now().Nanosecond(), value)
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

func (s *MemStorageEngine) GetShards() map[int]*sync.Map {
	return s.Shards
}

// GetShard returns a map of key-value pairs for the given shard
// In the case of RocksDB etc., it could be a snapshot file.
func (s *MemStorageEngine) GetShard(shard int) map[string]string {
	if _, ok := s.Shards[shard]; !ok {
		return nil
	}
	shardMap := make(map[string]string)
	s.Shards[shard].Range(func(key, value interface{}) bool {
		shardMap[key.(string)] = value.(string)
		return true
	})
	return shardMap
}

func (s *MemStorageEngine) DeleteShard(shard int) {
	delete(s.Shards, shard)
}

func (s *MemStorageEngine) SetShard(shardId int, shard map[string]string) {
	for key, value := range shard {
		s.Set(shardId, key, value)
	}
}
