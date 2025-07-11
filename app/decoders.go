package main

import (
	"errors"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

const (
	BencodeInteger = iota
	BencodeList
	BencodeDict
	BencodeString
	BencodeInvalid
)

func decodeBencode(bencodedString []byte) (any, error) {
	if len(bencodedString) == 0 {
		return 0, errors.New("decode error: the string is empty")
	}
	pos := 0
	var value any
	var err error
	for pos < len(bencodedString) {
		value, err = _decodeBencode(&pos, bencodedString)
		if err != nil {
			return value, err
		}
		pos++
	}
	return value, err
}
func _decodeBencode(pos *int, bencodedString []byte) (any, error) {
	switch findOutBencodeType(rune(bencodedString[*pos])) {
	case BencodeString:
		return decodeString(pos, bencodedString)
	case BencodeInteger:
		return decodeInteger(pos, bencodedString)
	case BencodeList:
		return decodeList(pos, bencodedString)
	case BencodeDict:
		return decodeDictionary(pos, bencodedString)
	default:
		break
	}
	return "", errors.New("decode error: invalid character")
}
func decodeDictionary(pos *int, bencodedString []byte) (map[string]any, error) {
	var newDict map[string]any = make(map[string]any)
	if len(bencodedString[*pos:]) < 2 {
		return newDict, errors.New("decode error: invalid dictionary")
	}
	*pos++
	var previous []byte
	for *pos < len(bencodedString) {
		if bencodedString[*pos] == 'e' {
			return newDict, nil
		}
		if !(findOutBencodeType(rune(bencodedString[*pos])) == BencodeString) {
			return newDict, errors.New("decode error: the key must be a string")
		}
		key, err := decodeString(pos, bencodedString)
		if err != nil {
			return newDict, err
		}
		if previous == nil {
			previous = key
		}
		current := key
		if string(previous) > string(current) {
			return newDict, errors.New("decode error:the dictionary keys aren't lexicographically sorted")
		}
		*pos++
		if *pos >= len(bencodedString) {
			return newDict, errors.New("decode error: malformed dict")
		}
		value, err := _decodeBencode(pos, bencodedString)
		if err != nil {
			return newDict, err
		}
		newDict[string(key)] = value
		*pos++
		previous = current
	}
	return newDict, errors.New("decode error: malformed dictionary")

}
func decodeList(pos *int, bencodedString []byte) ([]any, error) {
	newList := make([]any, 0)
	if len(bencodedString[*pos:]) < 2 {
		return newList, errors.New("decode error: list length must be greater than 2")
	}
	*pos++
	for *pos < len(bencodedString) {
		if bencodedString[*pos] == 'e' {
			return newList, nil
		}
		el, err := _decodeBencode(pos, bencodedString)
		if err != nil {
			return newList, err
		}
		newList = append(newList, el)
		*pos++
	}
	return newList, errors.New("decode error: list malformed")
}
func decodeString(pos *int, bencodedString []byte) ([]byte, error) {
	if len(bencodedString) == 0 || len(bencodedString[*pos:]) == 0 {
		return nil, errors.New("decode error: ")
	}
	sizeFirstPos := *pos
	for *pos < len(bencodedString) {
		if bencodedString[*pos] != ':' {
			*pos++
			continue
		}
		if *pos == 0 || *pos+1 == len(bencodedString) {
			break
		}
		length, err := strconv.Atoi(string(bencodedString[sizeFirstPos:*pos]))
		*pos++
		if err != nil || length > len(bencodedString)-(*pos) {
			break
		}
		start := *pos
		end := *pos + length
		*pos = end - 1
		return bencodedString[start:end], nil
	}
	return nil, errors.New("decode error: byte string couldn't be parsed")
}

func decodeInteger(pos *int, bencodedString []byte) (int, error) {
	bencodedStringLength := len(bencodedString)
	if err := validateMinimumIntegerLength(pos, bencodedString); err != nil {
		return 0, err
	}
	*pos++
	isNegative := bencodedString[*pos] == '-'
	if isNegative && !isValidNegativeNumber(pos, bencodedString) {
		return 0, errors.New("decode error: string composition isn't correct")
	}
	if isNegative {
		*pos++
	}
	if !isNegative && bencodedString[*pos] == '0' && bencodedString[*pos+1] != 'e' {
		return 0, errors.New("decode error: the composition i0xe isn't allowed")
	}
	start := *pos
	for *pos < bencodedStringLength {
		if bencodedString[*pos] == 'e' {
			break
		}
		*pos++
	}
	if *pos >= bencodedStringLength {
		return 0, errors.New("decode error: \"i\" termination character couldn't be found")
	}
	intValue, err := strconv.Atoi(string(bencodedString[start:*pos]))
	if err != nil {
		return 0, errors.New("decode error: malformed integer value")
	}
	if isNegative {
		intValue = (-1) * intValue
	}
	return intValue, err
}
func validateMinimumIntegerLength(pos *int, bencodedString []byte) error {
	if len(bencodedString) == 0 || len(bencodedString[*pos:]) < 3 {
		return errors.New("decode error: integer value couldn't be decoded")
	}
	return nil
}
func isValidNegativeNumber(pos *int, bencodedString []byte) bool {
	remaining := bencodedString[*pos:]
	if len(remaining) < 3 {
		return false
	}
	if remaining[1] == '0' {
		return false
	}
	return true
}

func findOutBencodeType(char rune) int {
	if char == 'i' {
		return BencodeInteger
	} else if char == 'l' {
		return BencodeList
	} else if char == 'd' {
		return BencodeDict
	} else if unicode.IsDigit(char) {
		return BencodeString
	} else {
		return BencodeInvalid
	}
}
