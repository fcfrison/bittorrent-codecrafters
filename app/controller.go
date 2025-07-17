package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
)

func decodeCommand(input []byte) {
	decoded, err := decodeBencode(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonOutput, err := json.Marshal(decoded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}
func calculateInfoHash(info map[string]any) ([20]byte, error) {
	bencodedInfo, err := EncodeDictionary(info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(bencodedInfo), nil
}
func infoCommand(input string) {
	data, err := os.ReadFile(input)
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
	infohash, err := calculateInfoHash(info)
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
	fmt.Printf("Info Hash: %x\n", infohash)
	fmt.Printf("Piece Length: %d\n", info["piece length"])
	fmt.Println("Piece Hashes:")
	for _, val := range pieceHashes {
		fmt.Printf("%x\n", val)
	}
}
func encodeCommand(input []byte) {
	decoded, err := decodeBencode(input)
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
}
func peersCommand(input string) {
	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("error: file couldn't be open")
		os.Exit(1)
	}
	parsedData, err := parseTorrentFile(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func generateDiscoverPeersParamsStruct(parsedData map[string]any, port int, peer_id [20]byte,
	uploaded int, downloaded int, left int, compact int) (*DiscoverPeersParams, error) {
	info := parsedData["info"].(map[string]any)
	info_hash, err := calculateInfoHash(info)
	if err != nil {
		return nil, err
	}
	return &DiscoverPeersParams{
		info_hash:  info_hash,
		peer_id:    peer_id,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       left,
		compact:    compact,
	}, nil

}
