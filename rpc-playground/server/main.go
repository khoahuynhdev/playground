package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type (
	Args       struct{}
	TimeServer int64
)

func (t *TimeServer) GiveServerTime(args *Args, reply *int64) error {
	*reply = time.Now().Unix()
	return nil
}

func main() {
	// create a new RPC server
	timeServer := new(TimeServer)

	// register RPC server
	rpc.Register(timeServer)

	// NOTE: can I register just a function, no method from object or receiver func?
	// short: kinda can't, ref: https://go.dev/src/net/rpc/server.go
	rpc.HandleHTTP()
	// listen for request on port 1234
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	http.Serve(l, nil)
}

// NOTE: no response from rpc server? check here
// https://stackoverflow.com/questions/17452199/no-response-from-gorilla-rpc-json-rpc-service
