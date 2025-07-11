package main

import "errors"

func PieceHashes(hashes []byte) ([][]byte, error) {
	isMultipleOf20 := len(hashes)%20 == 0
	if !isMultipleOf20 || len(hashes) == 0 {
		return nil, errors.New("error: the pieces field isn't incompatible with SHA1 expected length")
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
