package main

import (
	"go-RPC/tcp-server/remote"
	"log"
	"net"
	"net/rpc"
)

func main() {
	compose := new(remote.Compose)

	rpc.Register(compose)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// rpc.Accept(listener)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("listen error: %v", err)
		}

		go rpc.ServeConn(conn)
	}
}
