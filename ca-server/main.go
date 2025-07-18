package main

import (
	"ca-server/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
