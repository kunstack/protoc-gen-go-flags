package types_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBytesHexValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    []byte
		expected string
	}{
		{
			name:     "empty bytes",
			value:    []byte{},
			expected: "",
		},
		{
			name:     "nil bytes",
			value:    nil,
			expected: "",
		},
		{
			name:     "simple string",
			value:    []byte("hello"),
			expected: "68656C6C6F",
		},
		{
			name:     "binary data",
			value:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expected: "00010203FF",
		},
		{
			name:     "unicode text",
			value:    []byte("こんにちは"),
			expected: "E38193E38293E381ABE381A1E381AF",
		},
		{
			name:     "special characters",
			value:    []byte("!@#$%^&*()"),
			expected: "21402324255E262A2829",
		},
		{
			name:     "emoji",
			value:    []byte("👋"),
			expected: "F09F918B",
		},
		{
			name:     "single byte",
			value:    []byte{0x42},
			expected: "42",
		},
		{
			name:     "all zeros",
			value:    []byte{0x00, 0x00, 0x00},
			expected: "000000",
		},
		{
			name:     "all FF",
			value:    []byte{0xFF, 0xFF, 0xFF},
			expected: "FFFFFF",
		},
		{
			name:     "mixed case letters",
			value:    []byte("HelloWorld"),
			expected: "48656C6C6F576F726C64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := wrapperspb.Bytes(tt.value)
			b := types.BytesHex(value)
			if got := b.String(); got != tt.expected {
				t.Errorf("BytesHexValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBytesHexValue_String_Nil(t *testing.T) {
	var b *types.BytesHexValue
	got := b.String()
	expected := ""
	if got != expected {
		t.Errorf("BytesHexValue.String() nil = %v, want %v", got, expected)
	}
}

func TestBytesHexValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []byte
		expectError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "whitespace",
			input:    "   ",
			expected: []byte{},
		},
		{
			name:     "simple hex",
			input:    "68656C6C6F",
			expected: []byte("hello"),
		},
		{
			name:     "binary data hex",
			input:    "00010203FF",
			expected: []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
		},
		{
			name:     "unicode hex",
			input:    "E38193E38293E381ABE381A1E381AF",
			expected: []byte("こんにちは"),
		},
		{
			name:     "special characters hex",
			input:    "21402324255E262A2829",
			expected: []byte("!@#$%^&*()"),
		},
		{
			name:     "emoji hex",
			input:    "F09F918B",
			expected: []byte("👋"),
		},
		{
			name:     "single byte hex",
			input:    "42",
			expected: []byte{0x42},
		},
		{
			name:     "all zeros hex",
			input:    "000000",
			expected: []byte{0x00, 0x00, 0x00},
		},
		{
			name:     "all FF hex",
			input:    "FFFFFF",
			expected: []byte{0xFF, 0xFF, 0xFF},
		},
		{
			name:     "lowercase hex",
			input:    "deadbeef",
			expected: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:     "uppercase hex",
			input:    "DEADBEEF",
			expected: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:     "mixed case hex",
			input:    "DeadBeef",
			expected: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:     "hex with spaces",
			input:    "  68656C6C6F  ",
			expected: []byte("hello"),
		},
		{
			name:        "invalid hex - odd length",
			input:       "123",
			expectError: true,
		},
		{
			name:        "invalid hex - non-hex characters",
			input:       "hello world",
			expectError: true,
		},
		{
			name:        "invalid hex - contains letters beyond F",
			input:       "GH12",
			expectError: true,
		},
		{
			name:        "invalid hex - contains symbols",
			input:       "12!@",
			expectError: true,
		},
		{
			name:        "invalid hex - empty after trim",
			input:       "   ",
			expectError: false, // Should not error, just set empty bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesHexValue{}
			err := b.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("BytesHexValue.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("BytesHexValue.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if !bytes.Equal(b.Value, tt.expected) {
				t.Errorf("BytesHexValue.Set() = %v, want %v", b.Value, tt.expected)
			}
		})
	}
}

func TestBytesHexValue_Type(t *testing.T) {
	b := &types.BytesHexValue{}
	if got := b.Type(); got != "bytesHex" {
		t.Errorf("BytesHexValue.Type() = %v, want %v", got, "bytesHex")
	}
}

func TestBytesHexValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*types.BytesHexValue)(nil)
}

func TestBytesHexValue_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	b := &types.BytesHexValue{}

	fs.Var(b, "bytes", "Test bytes hex value")

	// Test valid input
	err := fs.Parse([]string{"--bytes", "68656C6C6F20776F726C64"})
	if err != nil {
		t.Fatalf("Failed to parse bytes: %v", err)
	}

	expected := []byte("hello world")
	if !bytes.Equal(b.Value, expected) {
		t.Errorf("Parsed bytes = %v, want %v", b.Value, expected)
	}

	// Test String() method after Set()
	if got := b.String(); got != "68656C6C6F20776F726C64" {
		t.Errorf("BytesHexValue.String() after Set() = %v, want %v", got, "68656C6C6F20776F726C64")
	}
}

func TestBytesHexValue_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	b := &types.BytesHexValue{}

	fs.Var(b, "bytes", "Test bytes hex value")

	// Test invalid input
	err := fs.Parse([]string{"--bytes", "invalid hex"})
	if err == nil {
		t.Fatalf("Expected error for invalid hex input, got nil")
	}
}

func TestBytesHex_Constructor(t *testing.T) {
	wrap := &wrapperspb.BytesValue{Value: []byte("hello world")}
	b := types.BytesHex(wrap)

	// Test that the returned BytesHexValue has correct fields
	if (*wrapperspb.BytesValue)(b) != wrap {
		t.Errorf("BytesHex() wrapper = %p, want %p", (*wrapperspb.BytesValue)(b), wrap)
	}

	// Test that the value is correct
	if !bytes.Equal(b.Value, []byte("hello world")) {
		t.Errorf("BytesHex().Value = %v, want %v", b.Value, []byte("hello world"))
	}

	// Test String() method on constructed bytes
	expected := "68656C6C6F20776F726C64"
	if got := b.String(); got != expected {
		t.Errorf("BytesHex().String() = %v, want %v", got, expected)
	}
}

func TestBytesHex_ConstructorNil(t *testing.T) {
	b := types.BytesHex(nil)

	// Test that it works with nil wrapper
	if b != nil {
		t.Errorf("BytesHex(nil) = %p, want nil", b)
	}
}

func TestBytesHexValue_StringAfterSet(t *testing.T) {
	b := &types.BytesHexValue{}

	// Set a value
	err := b.Set("746573742076616C7565")
	if err != nil {
		t.Fatalf("Failed to set bytes hex value: %v", err)
	}

	// Test String() method
	if got := b.String(); got != "746573742076616C7565" {
		t.Errorf("BytesHexValue.String() after Set() = %v, want %v", got, "746573742076616C7565")
	}
}

func TestBytesHexValue_MultipleSets(t *testing.T) {
	b := &types.BytesHexValue{}

	// Set first value
	err := b.Set("6669727374")
	if err != nil {
		t.Fatalf("Failed to set first bytes hex value: %v", err)
	}
	if !bytes.Equal(b.Value, []byte("first")) {
		t.Errorf("First set failed: got %v, want %v", b.Value, []byte("first"))
	}

	// Set second value
	err = b.Set("7365636F6E64")
	if err != nil {
		t.Fatalf("Failed to set second bytes hex value: %v", err)
	}
	if !bytes.Equal(b.Value, []byte("second")) {
		t.Errorf("Second set failed: got %v, want %v", b.Value, []byte("second"))
	}

	// Test final String() value
	if got := b.String(); got != "7365636F6E64" {
		t.Errorf("Final String() = %v, want %v", got, "7365636F6E64")
	}
}

func TestBytesHexValue_RoundTrip(t *testing.T) {
	tests := [][]byte{
		{},
		[]byte("hello"),
		[]byte{0x00, 0x01, 0x02, 0xFF},
		[]byte("こんにちは"),
		[]byte("!@#$%^&*()"),
		[]byte("👋"),
		make([]byte, 1000),
	}

	for i, original := range tests {
		t.Run(fmt.Sprintf("roundtrip_%d", i), func(t *testing.T) {
			b := &types.BytesHexValue{}

			// Convert to hex and set
			hexStr := hex.EncodeToString(original)
			err := b.Set(hexStr)
			if err != nil {
				t.Fatalf("Failed to set bytes hex value: %v", err)
			}

			// Convert back and compare
			if !bytes.Equal(b.Value, original) {
				t.Errorf("Round trip failed: got %v, want %v", b.Value, original)
			}
		})
	}
}

func TestBytesHexValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{"empty", "", []byte{}},
		{"whitespace", "   ", []byte{}},
		{"single byte", "41", []byte("A")},
		{"null byte", "00", []byte{0x00}},
		{"max byte", "FF", []byte{0xFF}},
		{"mixed ASCII", "48656C6C6F20576F726C64", []byte("Hello World")},
		{"binary with null", "00666F6F", []byte{0x00, 'f', 'o', 'o'}},
		{"all ASCII", "4142434445464748494A4B4C4D4E4F505152535455565758595A", []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
		{"common pattern", "DEADBEEF", []byte{0xDE, 0xAD, 0xBE, 0xEF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesHexValue{}
			err := b.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set bytes hex value %q: %v", tt.input, err)
			}
			if !bytes.Equal(b.Value, tt.want) {
				t.Errorf("BytesHexValue edge case: input=%q, got=%v, want=%v", tt.input, b.Value, tt.want)
			}
		})
	}
}

func TestBytesHexValue_LongData(t *testing.T) {
	b := &types.BytesHexValue{}

	// Test with a large byte array
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Convert to hex and set
	hexStr := hex.EncodeToString(largeData)
	err := b.Set(hexStr)
	if err != nil {
		t.Fatalf("Failed to set large bytes hex value: %v", err)
	}

	if !bytes.Equal(b.Value, largeData) {
		t.Errorf("Large data round trip failed")
	}

	// Test String() method produces valid hex
	if got := b.String(); got != strings.ToUpper(hexStr) {
		t.Errorf("Large data String() mismatch")
	}
}

func TestBytesHexValue_CaseSensitivity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		value    []byte
	}{
		{
			name:     "lowercase input",
			input:    "deadbeef",
			expected: "DEADBEEF",
			value:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:     "uppercase input",
			input:    "DEADBEEF",
			expected: "DEADBEEF",
			value:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:     "mixed case input",
			input:    "DeadBeef",
			expected: "DEADBEEF",
			value:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesHexValue{}
			err := b.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set bytes hex value: %v", err)
			}

			if !bytes.Equal(b.Value, tt.value) {
				t.Errorf("Value mismatch: got %v, want %v", b.Value, tt.value)
			}

			if got := b.String(); got != tt.expected {
				t.Errorf("String() mismatch: got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBytesHexValue_InvalidInputs(t *testing.T) {
	tests := []string{
		"xyz",    // Invalid hex characters
		"123",    // Odd length
		"12 34",  // Contains spaces (should be trimmed)
		"12!@",   // Contains symbols
		"ghij",   // Letters beyond f
		"0x1234", // Hex prefix not supported
		"123h",   // Trailing non-hex
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			b := &types.BytesHexValue{}
			err := b.Set(input)

			// Only expect error for truly invalid hex (not just whitespace)
			if input == "   " {
				if err != nil {
					t.Errorf("Whitespace should not cause error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid hex input %q, got nil", input)
				}
			}
		})
	}
}
