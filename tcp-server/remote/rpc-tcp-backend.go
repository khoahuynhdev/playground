package remote

import "fmt"

// TCPArgs is structured around the client's provided parameters
// The struct's fields need to be exported too!
type TCPArgs struct {
	Foo string
	Bar string
}

// Compose is our RPC functions return type
type Compose string

func (c *Compose) Details(args *TCPArgs, reply *string) error {
	fmt.Printf("Args received: %+v\n", args)
	*c = "Value from RPC"
	*reply = "Happy New Year!"
	return nil
}
