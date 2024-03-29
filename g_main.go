package main

import (
	"flag"
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

	kv := NewDistKVServer(config)
	kv.Bootstrap()
}

// -----------------Utils -------------------
