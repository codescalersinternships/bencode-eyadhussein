// Package bencode encode and decodes data in bencode format that is used by the BitTorrent peer-to-peer file sharing protocol.
package bencode

import (
	"errors"
	"strconv"
	"strings"
)

var (
	missingColonErr      = errors.New("invalid format missing ':' separator")
	invalidCharErr       = errors.New("invalid character in string length")
	strLenExceedsDataErr = errors.New("string length exceeds available data")
	invalidFormat        = errors.New("invalid format")
	unsupportedType      = errors.New("unsupported type")
)

// Decode takes a bencoded string and returns the corresponding structure
func Decode(benc string) (any, error) {
	if len(benc) == 0 {
		return nil, invalidFormat
	}

	switch benc[0] {
	case 'i':
		if len(benc) < 2 || benc[len(benc)-1] != 'e' {
			return nil, invalidFormat
		}
		return decodeInt(benc)
	case 'l':
		if len(benc) < 2 || benc[len(benc)-1] != 'e' {
			return nil, invalidFormat
		}
		return decodeList(benc)
	case 'd':
		if len(benc) < 2 || benc[len(benc)-1] != 'e' {
			return nil, invalidFormat
		}
		return decodeDict(benc)
	default:
		return decodeString(benc)
	}
}

// Encode takes a structure that can be encoded in bencode and returns the bencoded string
func Encode(benc any) (string, error) {
	switch input := benc.(type) {
	case int:
		if input == 0 {
			return "i0e", nil
		}
		return "i" + strconv.Itoa(input) + "e", nil

	case string:
		return strconv.Itoa(len(input)) + ":" + input, nil

	case []any:
		var result string
		for _, item := range input {
			encoded, err := Encode(item)
			if err != nil {
				return "", err
			}
			result += encoded
		}
		return "l" + result + "e", nil

	case map[string]any:
		var result string
		for key, value := range input {
			encodedKey, err := Encode(key)
			if err != nil {
				return "", err
			}
			encodedValue, err := Encode(value)
			if err != nil {
				return "", err
			}
			result += encodedKey + encodedValue
		}
		return "d" + result + "e", nil

	default:
		return "", unsupportedType
	}
}

func decodeInt(benc string) (int, error) {
	end := strings.IndexByte(benc[1:], 'e')
	intStr := benc[1 : end+1]
	return strconv.Atoi(intStr)
}

func decodeList(benc string) ([]any, error) {
	benc = benc[1:]
	var result []any

	for len(benc) > 0 && benc[0] != 'e' {
		item, remaining, err := decodeItem(benc)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
		benc = remaining
	}

	return result, nil
}
func decodeDict(benc string) (map[string]any, error) {

	benc = benc[1:]
	result := make(map[string]any)

	for len(benc) > 0 && benc[0] != 'e' {
		key, remaining, err := decodeItem(benc)
		if err != nil {
			return nil, err
		}

		value, remaining, err := decodeItem(remaining)
		if err != nil {
			return nil, err
		}

		result[key.(string)] = value
		benc = remaining
	}

	return result, nil
}

func decodeItem(benc string) (any, string, error) {
	item, err := Decode(benc)
	if err != nil {
		return nil, "", err
	}

	encoded, err := Encode(item)
	if err != nil {
		return nil, "", err
	}

	return item, benc[len(encoded):], nil
}

func decodeString(benc string) (string, error) {
	lenContentPair := strings.SplitN(benc, ":", 2)
	if len(lenContentPair) != 2 {
		return "", missingColonErr
	}

	lengthStr := lenContentPair[0]
	for _, ch := range lengthStr {
		if ch < '0' || ch > '9' {
			return "", invalidCharErr
		}
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	content := lenContentPair[1]
	if length > len(content) {
		return "", strLenExceedsDataErr
	}

	return content[:length], nil
}
