package main

import "sync"

type StorageEngine struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewStorageEngine() *StorageEngine {
	return &StorageEngine{
		store: make(map[string]string),
	}
}

func (kv *StorageEngine) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
}

func (kv *StorageEngine) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.store[key]
	return value, exists
}
