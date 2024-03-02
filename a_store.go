package main

import (
	"sync"
)

type KeyValueStore struct {
	data sync.Map
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data: sync.Map{},
	}
}

func (s *KeyValueStore) Set(key string, value string) {
	s.data.Store(key, value)
}

func (s *KeyValueStore) Get(key string) (string, bool) {
	value, ok := s.data.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (s *KeyValueStore) Delete(key string) {
	s.data.Delete(key)
}
