package main

import (
	"sync"
)

type KVStore struct {
	data sync.Map
}

func NewKeyValueStore() *KVStore {
	return &KVStore{
		data: sync.Map{},
	}
}

func (s *KVStore) Set(key string, value string) {
	s.data.Store(key, value)
}

func (s *KVStore) Get(key string) (string, bool) {
	value, ok := s.data.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (s *KVStore) Delete(key string) {
	s.data.Delete(key)
}
