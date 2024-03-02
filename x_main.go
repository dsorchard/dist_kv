package main

import (
	"dist_kv"
	"flag"
)

func main() {
	config := dist_kv.loadConfig()

	var gossipPort, httpPort int
	flag.IntVar(&gossipPort, "gossip", config.InternalPort, "port number for gossip protocol")
	flag.IntVar(&httpPort, "http", config.ExternalPort, "port number for http server")
	flag.Parse()

	config.Host = dist_kv.GetLocalIP()
	config.InternalPort = gossipPort
	config.ExternalPort = httpPort

	kv := dist_kv.NewDistKV(config)
	kv.Bootstrap()
}
