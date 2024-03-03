package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

type Client interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	PutShard(shardId int, shard map[string]string) error
}

type HttpClient struct {
	ring   HashRing
	quorum int
	api    *HttpAPIServer
}

func NewHttpClient(ring HashRing, api *HttpAPIServer) Client {
	return &HttpClient{
		ring:   ring,
		api:    api,
		quorum: (ring.ReplicationFactor() / 2) + 1,
	}
}

func (c *HttpClient) Put(key string, value string) error {
	localNodeAddress := c.api.GetAddress()

	replicas := c.ring.ResolveNodes(key, c.quorum)
	results := make([]string, 0)

	for _, replica := range replicas {
		if replica == localNodeAddress {
			shardId := c.api.distKV.ring.ResolvePartitionID(key)
			c.api.distKV.store.Set(shardId, key, value)
			results = append(results, "ok")
			continue
		}
		url := fmt.Sprintf("http://%s/store/%s/%s", replica, key, value)
		resp, err := http.Post(url, "", nil)
		if err != nil {
			// TODO: we should not throw error. But here we are doing it for easy debugging.
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			results = append(results, "ok")
		}
	}
	if len(results) < c.quorum {
		return fmt.Errorf("not enough replicas have been updated")
	}
	return nil
}

func (c *HttpClient) Get(key string) (string, error) {
	replicas := c.ring.ResolveNodes(key, c.quorum)
	results := make([]string, 0)
	for _, replica := range replicas {
		if replica == c.api.GetAddress() {
			shardId := c.api.distKV.ring.ResolvePartitionID(key)
			if value, ok := c.api.distKV.store.Get(shardId, key); !ok {
				//TODO: need to handle empty value better
				results = append(results, "")
			} else {
				results = append(results, value)
			}
			continue
		}
		url := fmt.Sprintf("http://%s/store/%s", replica, key)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			results = append(results, string(body))
		}

	}
	if len(results) == c.quorum {
		// pick the latest timestamp
		// value is of the form "timestamp:value", so we can just sort the strings.
		sort.Slice(results, func(i, j int) bool {
			return results[i] > results[j]
		})
		return results[0], nil
	}

	return "", nil
}

// PutShard sends a shard to the replicas
// TODO: analyze if this will result in duplicates
func (c *HttpClient) PutShard(shardId int, shard map[string]string) error {
	members := c.ring.ResolveNodesForPartition(shardId, c.ring.ReplicationFactor())
	results := make([]string, 0)
	for _, member := range members {
		url := fmt.Sprintf("http://%s/shard/%d", member, shardId)

		jsonData, err := json.Marshal(shard)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			results = append(results, "ok")
		}
	}
	if len(results) < c.quorum {
		return fmt.Errorf("not enough replicas have been updated")
	}
	return nil
}
