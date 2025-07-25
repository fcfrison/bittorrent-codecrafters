package main

import (
	"crypto/sha1"
	"errors"
)

func PieceHashes(hashes []byte) ([][]byte, error) {
	isMultipleOf20 := len(hashes)%20 == 0
	if !isMultipleOf20 || len(hashes) == 0 {
		return nil, errors.New("error: the pieces field is incompatible with SHA1 expected length")
	}
	nrPieces := len(hashes) / 20
	listOfHashes := make([][]byte, nrPieces)
	for i := range nrPieces {
		j := i * 20
		listItem := hashes[j : j+20]
		listOfHashes[i] = listItem
	}
	return listOfHashes, nil
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
func calculateInfoHash(info map[string]any) ([20]byte, error) {
	bencodedInfo, err := EncodeDictionary(info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(bencodedInfo), nil
}
func fmtHandshakeMsg(infoHash *[20]byte, peerId *[20]byte) []byte {
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
