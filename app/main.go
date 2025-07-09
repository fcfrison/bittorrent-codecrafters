package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]
	arg2 := os.Args[2]
	switch command {
	case "decode":
		decoded, err := decodeBencode(arg2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	case "info":
		data, err := os.ReadFile(arg2)
		if err != nil {
			fmt.Println("error: file couldn't be open")
			os.Exit(1)
		}
		parsedData, err := parseTorrentFile(string(data))
		if err != nil {
			fmt.Println(err)
		}
		info := parsedData["info"].(map[string]any)
		fmt.Printf("Tracker URL: %s\n", parsedData["announce"])
		fmt.Printf("Length: %d", info["length"])
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
func parseTorrentFile(data string) (map[string]any, error) {
	if findOutBencodeType(rune(data[0])) != BencodeDict {
		return nil, errors.New("parse error: the .torrent file couldn't be parsed")
	}
	parsedInfo, err := decodeBencode(data)
	if err != nil {
		return nil, err
	}
	return parsedInfo.(map[string]any), nil
}
