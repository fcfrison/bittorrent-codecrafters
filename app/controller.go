package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var peerId string = "12345678901234567890"

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
func peersCommand(input string, client *BitTorrentTrackerClient) {
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
	client.SetUrl(string(parsedData["announce"].([]byte)))
	var peer_id [20]byte
	copy(peer_id[:], peerId[:])
	info := parsedData["info"].(map[string]any)
	info_hash, err := calculateInfoHash(info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	params := NewDiscoverPeersParamsStruct(info_hash, peer_id, 0, 0, info["length"].(int), 1)
	resp, err := client.DiscoverPeers(params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decodedBencode, err := decodeBencode(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decodedValue := decodedBencode.(map[string]any)
	if decodedValue["peers"] == nil {
		fmt.Println("error: no peer was announced by the tracker server")
		os.Exit(1)
	}
	peers, err := parsePeers(decodedValue["peers"].([]byte))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, val := range peers {
		fmt.Printf("%s\n", val.StrRepr())
	}

}
func fmtHandshakeMsg(infoHash [20]byte, peerId [20]byte) []byte {
	hndShkReq := make([]byte, 68)
	hndShkReq[0] = 0x13
	for i, val := range []byte("BitTorrent protocol") {
		hndShkReq[i+1] = val
	}
	for i := 0; i < 20; i++ {
		hndShkReq[i+28] = infoHash[i]
		peerId[i+48] = peerId[i]

	}
	return hndShkReq
}
func parseHandshakeMsg(msg [68]byte, infoHash [20]byte) bool {
	if msg[0] != 0x13 || string(msg[1:21]) == "BitTorrent protocol" {
		return false
	}
	for _, val := range msg[20:28] {
		if !(val == 0x00) {
			return false
		}
	}
	for i, val := range infoHash {
		if val != msg[28+i] {
			return false
		}
	}
	return true
}
func handshakeCmd() {

}
