package main

import (
	"flag"
	"fmt"
)

func main() {

	var gossipPort, httpPort int
	flag.IntVar(&gossipPort, "gossip", 8000, "port number for gossip protocol")
	flag.IntVar(&httpPort, "http", 8001, "port number for http server")
	flag.Parse()

	bootstrap := "127.0.0.1:8000"
	gossipHost := fmt.Sprintf("127.0.0.1:%d", gossipPort)

	app := NewDistKVServer(gossipHost, "kv-store", bootstrap, httpPort)
	app.setupHandlers()

	select {}
}
