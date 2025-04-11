package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type args struct {
	Foo, Bar string
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string

	e := client.Call("Compose.Details", &args{"Foo!", "Bar!"}, &reply)
	if e != nil {
		log.Fatalf("Something went wrong: %v", e.Error())
	}

	fmt.Printf("The 'reply' pointer value has been changed to: %s", reply)
}
