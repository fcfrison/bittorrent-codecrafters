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

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (any, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	} else if bencodedString[0] == 'i' {
		return decodeInteger(bencodedString)
	}
	return "", fmt.Errorf("decode error: it wasn't possible to decode the value")

}
func decodeString(bencodedString string) (string, error) {
	var firstColonIndex int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}
func decodeInteger(bencodedString string) (int, error) {
	strLen := len(bencodedString)
	if strLen < 3 {
		return 0, errors.New("decode error: string length must be >= 3")
	}
	if bencodedString[strLen-1] != 'e' {
		return 0, errors.New("decode error: the string must end with 'i'")
	}
	if bencodedString[1] == '0' && strLen > 3 {
		return 0, errors.New("decode error: integers starting with 0 can't have more than one digit")
	}
	var isNegative bool = bencodedString[1] == '-'
	if isNegative && bencodedString[2] == '0' {
		return 0, errors.New("decode error: the combo -0 is a invalid one")
	}
	if isNegative {
		return strconv.Atoi(bencodedString[2 : strLen-1])
	}
	return strconv.Atoi(bencodedString[1 : strLen-1])

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
