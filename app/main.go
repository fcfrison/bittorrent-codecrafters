package main

import (
	"crypto/sha1"
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
		pieceHashes, err := PieceHashes(info["pieces"].([]byte))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Tracker URL: %s\n", parsedData["announce"])
		fmt.Printf("Length: %d\n", info["length"])
		fmt.Printf("Info Hash: %x\n", sha1.Sum(bencodedInfo))
		fmt.Printf("Piece Length: %d\n", info["piece length"])
		fmt.Println("Piece Hashes:")
		for _, val := range pieceHashes {
			fmt.Printf("%x\n", val)
		}
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
