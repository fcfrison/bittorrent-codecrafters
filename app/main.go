package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var _ = json.Marshal

func main() {
	command := os.Args[1]
	arg2 := os.Args[2]
	switch command {
	case "decode":
		decodeCommand([]byte(arg2))
	case "info":
		infoCommand(arg2)
	case "encode":
		encodeCommand([]byte(arg2))
	case "peers":
		peersCommand(arg2)
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
