package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func decodeBencode(bencodedString string) (any, error) {
	if len(bencodedString) == 0 {
		return 0, errors.New("decode error: the string is empty")
	}
	var bencodeRunes []rune = []rune(bencodedString)
	pos := 0
	var value any
	var err error
	for pos < len(bencodeRunes) {
		value, err = _decodeBencode(bencodeRunes, &pos)
		if err != nil {
			return value, err
		}
		pos++
	}
	return value, err
}
func _decodeBencode(bencodeRunes []rune, pos *int) (any, error) {
	switch findOutBencodeType(bencodeRunes[*pos]) {
	case BencodeString:
		return decodeString(pos, bencodeRunes)
	case BencodeInteger:
		return decodeInteger(pos, bencodeRunes)
	case BencodeList:
		return decodeList(pos, bencodeRunes)
	default:
		break
	}
	return "", errors.New("decode error: invalid character")
}
func decodeList(pos *int, bencodeRunes []rune) ([]any, error) {
	newList := make([]any, 0)
	if len(bencodeRunes[*pos:]) < 2 {
		return newList, errors.New("decode error: list length must be greater than 2")
	}
	*pos++
	for *pos < len(bencodeRunes) {
		if bencodeRunes[*pos] == 'e' {
			return newList, nil
		}
		el, err := _decodeBencode(bencodeRunes, pos)
		if err != nil {
			return newList, err
		}
		newList = append(newList, el)
		*pos++
	}
	return newList, errors.New("decode error: list malformed")
}
func decodeString(pos *int, bencodeRunes []rune) (string, error) {
	if len(bencodeRunes) == 0 || len(bencodeRunes[*pos:]) == 0 {
		return "", errors.New("decode error: ")
	}
	sizeFirstPos := *pos
	for *pos < len(bencodeRunes) {
		if bencodeRunes[*pos] != ':' {
			*pos++
			continue
		}
		if *pos == 0 || *pos+1 == len(bencodeRunes) {
			break
		}
		length, err := strconv.Atoi(string(bencodeRunes[sizeFirstPos:*pos]))
		*pos++
		if err != nil || length > len(bencodeRunes)-(*pos) {
			break
		}
		start := *pos
		end := *pos + length
		returnString := string(bencodeRunes[start:end])
		*pos = end - 1
		return returnString, nil
	}
	return "", errors.New("decode error: byte string couldn't be parsed")
}

func decodeInteger(pos *int, bencodeRunes []rune) (int, error) {
	bencodeRunesLength := len(bencodeRunes)
	if err := validateMinimumIntegerLength(pos, bencodeRunes); err != nil {
		return 0, err
	}
	*pos++
	isNegative := bencodeRunes[*pos] == '-'
	if isNegative && !isValidNegativeNumber(pos, bencodeRunes) {
		return 0, errors.New("decode error: string composition isn't correct")
	}
	if isNegative {
		*pos++
	}
	if !isNegative && bencodeRunes[*pos] == '0' && bencodeRunes[*pos+1] != 'e' {
		return 0, errors.New("decode error: the composition i0xe isn't allowed")
	}
	start := *pos
	for *pos < bencodeRunesLength {
		if bencodeRunes[*pos] == 'e' {
			break
		}
		*pos++
	}
	if *pos >= bencodeRunesLength {
		return 0, errors.New("decode error: \"i\" termination character couldn't be found")
	}
	intValue, err := strconv.Atoi(string(bencodeRunes[start:*pos]))
	if err != nil {
		return 0, errors.New("decode error: malformed integer value")
	}
	if isNegative {
		intValue = (-1) * intValue
	}
	return intValue, err
}
func validateMinimumIntegerLength(pos *int, bencodeRunes []rune) error {
	if len(bencodeRunes) == 0 || len(bencodeRunes[*pos:]) < 3 {
		return errors.New("decode error: integer value couldn't be decoded")
	}
	return nil
}
func isValidNegativeNumber(pos *int, bencodeRunes []rune) bool {
	remaining := bencodeRunes[*pos:]
	if len(remaining) < 3 {
		return false
	}
	if remaining[1] == '0' {
		return false
	}
	return true
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345

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

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
