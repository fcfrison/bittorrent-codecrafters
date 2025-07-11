package main

import (
	"crypto/sha1"
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
		decoded, err := decodeBencode([]byte(arg2))
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
		parsedData, err := parseTorrentFile(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		info := parsedData["info"].(map[string]any)
		bencodedInfo, err := EncodeDictionary(info)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		hash := sha1.Sum(bencodedInfo)
		fmt.Printf("Tracker URL: %s\n", parsedData["announce"])
		fmt.Printf("Length: %d\n", info["length"])
		fmt.Printf("Info Hash: %x\n", hash)
	case "encode":
		decoded, err := decodeBencode([]byte(arg2))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var encoded []byte
		switch i := decoded.(type) {
		case int:
			encoded, err = EncodeInteger(i)
		case []byte:
			encoded, err = EncodeString(i)
		case map[string]any:
			encoded, err = EncodeDictionary(i)
		case []any:
			encoded, err = EncodeList(i)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(encoded))
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
func parseTorrentFile(data []byte) (map[string]any, error) {
	if findOutBencodeType(rune(data[0])) != BencodeDict {
		return nil, errors.New("parse error: the .torrent file couldn't be parsed")
	}
	parsedInfo, err := decodeBencode(data)
	if err != nil {
		return nil, err
	}
	return parsedInfo.(map[string]any), nil
}
