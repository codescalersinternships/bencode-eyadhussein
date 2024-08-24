package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validDictEncoded = "d3:bar4:spam3:fooi42ee"
	validDictDecoded = map[string]any{
		"bar": "spam",
		"foo": 42,
	}

	validNestedDictEncoded = "d3:bar4:spam3:food3:bar4:spam3:fooi42eee"
	validNestedDictDecoded = map[string]any{
		"bar": "spam",
		"foo": map[string]any{
			"bar": "spam",
			"foo": 42,
		},
	}
)

func TestDecode(t *testing.T) {

	validDecodeTests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "simple string",
			input:    "4:spam",
			expected: "spam",
		},
		{
			name:     "integer",
			input:    "i42e",
			expected: 42,
		},
		{
			name:     "negative integer",
			input:    "i-42e",
			expected: -42,
		},
		{
			name:     "list",
			input:    "l4:spami42ee",
			expected: []any{"spam", 42},
		},
		{
			name:     "dictionary",
			input:    validDictEncoded,
			expected: validDictDecoded,
		},
		{
			name:     "nested dictionary",
			input:    validNestedDictEncoded,
			expected: validNestedDictDecoded,
		},
	}

	for _, tt := range validDecodeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}

	invalidDecodeTests := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "invalid encoded string",
			input: "invalid",
			err:   invalidFormat,
		},
		{
			name: "empty input",
			err:  invalidFormat,
		},
		{
			name:  "missing colon",
			input: "4spam",
			err:   missingColonErr,
		},
		{
			name:  "invalid character in string length",
			input: "x:spam",
			err:   invalidCharErr,
		},
		{
			name:  "string length exceeds available data",
			input: "10:spamx",
			err:   strLenExceedsDataErr,
		},
	}

	for _, tt := range invalidDecodeTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			assert.Equal(t, tt.err, err)
		})
	}

}

func TestEncode(t *testing.T) {
	validEncodeTests := []struct {
		name  string
		input any
	}{
		{
			name:  "simple string",
			input: "spam",
		},
		{
			name:  "integer",
			input: 42,
		},
		{
			name:  "negative integer",
			input: -42,
		},
		{
			name:  "list",
			input: []any{"spam", 42},
		},
		{
			name:  "dictionary",
			input: map[string]any{"foo": "bar", "baz": 123},
		},
		{
			name:  "nested dictionary",
			input: map[string]any{"outer": map[string]any{"inner": "value"}},
		},
	}

	for _, tt := range validEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := Encode(tt.input)
			assert.NoError(t, err)

			decoded, err := Decode(encoded)
			assert.NoError(t, err)

			assert.Equal(t, tt.input, decoded)
		})
	}

	invalidEncodeTests := []struct {
		name  string
		input any
		err   error
	}{
		{
			name:  "unsupported type",
			input: make(chan int),
			err:   unsupportedType,
		},
	}

	for _, tt := range invalidEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encode(tt.input)
			assert.Equal(t, tt.err, err)
		})
	}
}
