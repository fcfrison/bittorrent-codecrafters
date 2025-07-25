package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var _ = json.Marshal

func main() {
	if len(os.Args) < 3 {
		fmt.Println("error: the ")
	}
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
		client := NewBitTorrentTrackerClient()
		peersCommand(arg2, client)
	case "handshake":
		errMsg := "error: here is how you should execute your program: handshake sample.torrent <peer_ip>:<peer_port> "
		if len(os.Args) < 4 {
			fmt.Println(errMsg)
			os.Exit(1)
		}
		peerInfo := strings.Split(os.Args[3], ":")
		if len(peerInfo) != 2 {
			fmt.Println(errMsg)
			os.Exit(1)
		}
		peer_ip, peer_port := peerInfo[0], peerInfo[1]
		peer_port_n, err := strconv.Atoi(peer_port)
		if err != nil {
			fmt.Println("error: peer_port isn't in a numeric format")
			os.Exit(1)
		}
		clientConfig := NewClientConfig(peer_ip, peer_port_n)
		tcpClient, err := NewBitTorrentTcpClient(clientConfig)
		if err != nil {
			fmt.Println(err)
		}
		handshakeCommand(arg2, tcpClient)
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
