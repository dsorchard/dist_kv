package main

import (
	"flag"
)

const (
	localIp = "127.0.0.1"
)

func main() {
	config := loadConfig()

	var gossipPort, httpPort int
	flag.IntVar(&gossipPort, "gossip", config.InternalPort, "port number for gossip protocol")
	flag.IntVar(&httpPort, "http", config.ExternalPort, "port number for http server")
	flag.Parse()

	config.Host = GetLocalIP()
	config.InternalPort = gossipPort
	config.ExternalPort = httpPort

	kv := NewDistKV(config)
	kv.Bootstrap()
}

func GetLocalIP() string {
	return localIp
}
