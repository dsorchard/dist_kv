package main

import (
	"sync"
)

type KVStore struct {
	Data map[int]*sync.Map // Use a pointer to sync.Map
}

func NewKVStore() *KVStore {
	return &KVStore{
		Data: make(map[int]*sync.Map), // Initialize the map
	}
}

func (s *KVStore) Set(shard int, key string, value string) {
	if _, ok := s.Data[shard]; !ok {
		s.Data[shard] = &sync.Map{}
	}
	s.Data[shard].Store(key, value)
}

func (s *KVStore) Get(shard int, key string) (string, bool) {
	if _, ok := s.Data[shard]; !ok {
		return "", false
	}
	value, ok := s.Data[shard].Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}
