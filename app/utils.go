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
