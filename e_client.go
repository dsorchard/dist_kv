package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	address string
}

func NewHttpClient(address string) *HttpClient {
	return &HttpClient{
		address: address,
	}
}

func (c *HttpClient) Put(key string, value string) error {
	url := fmt.Sprintf("http://%s/put/%s/%s", c.address, key, value)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error forwarding request")
	}
	return nil
}

func (c *HttpClient) Get(key string) (string, error) {
	url := fmt.Sprintf("http://%s/get/%s", c.address, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("key not found")
	} else if resp.StatusCode != http.StatusOK {
		return "", errors.New("error retrieving key from responsible node")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
