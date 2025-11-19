package types_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBytesValue_String(t *testing.T) {
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
			expected: "aGVsbG8=",
		},
		{
			name:     "binary data",
			value:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expected: "AAECA/8=",
		},
		{
			name:     "unicode text",
			value:    []byte("こんにちは"),
			expected: "44GT44KT44Gr44Gh44Gv",
		},
		{
			name:     "long text",
			value:    []byte("this is a very long string that contains multiple words and should be handled correctly by the BytesValue implementation"),
			expected: "dGhpcyBpcyBhIHZlcnkgbG9uZyBzdHJpbmcgdGhhdCBjb250YWlucyBtdWx0aXBsZSB3b3JkcyBhbmQgc2hvdWxkIGJlIGhhbmRsZWQgY29ycmVjdGx5IGJ5IHRoZSBCeXRlc1ZhbHVlIGltcGxlbWVudGF0aW9u",
		},
		{
			name:     "special characters",
			value:    []byte("!@#$%^&*()"),
			expected: "IUAjJCVeJiooKQ==",
		},
		{
			name:     "emoji",
			value:    []byte("👋"),
			expected: "8J+Riw==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesValue{Value: tt.value}
			if got := b.String(); got != tt.expected {
				t.Errorf("BytesValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBytesValue_String_Nil(t *testing.T) {
	var b *types.BytesValue
	got := b.String()
	expected := ""
	if got != expected {
		t.Errorf("BytesValue.String() nil = %v, want %v", got, expected)
	}
}

func TestBytesValue_Set(t *testing.T) {
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
			name:     "simple base64",
			input:    "aGVsbG8=",
			expected: []byte("hello"),
		},
		{
			name:     "binary data base64",
			input:    "AAECA/8=",
			expected: []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
		},
		{
			name:     "unicode base64",
			input:    "44GT44KT44Gr44Gh44Gv",
			expected: []byte("こんにちは"),
		},
		{
			name:     "special characters base64",
			input:    "IUAjJCVeJiooKQ==",
			expected: []byte("!@#$%^&*()"),
		},
		{
			name:     "emoji base64",
			input:    "8J+Riw==",
			expected: []byte("👋"),
		},
		{
			name:     "base64 with padding",
			input:    "dGhpcyBpcyBhIHRlc3Q=",
			expected: []byte("this is a test"),
		},
		{
			name:        "base64 without padding",
			input:       "dGhpcyBpcyBhIHRlc3Q",
			expected:    []byte("this is a test"),
			expectError: true,
		},
		{
			name:        "invalid base64",
			input:       "invalid base64!!!",
			expectError: true,
		},
		{
			name:        "incorrect padding",
			input:       "aGVsbG8",
			expectError: true,
		},
		{
			name:        "non-base64 characters",
			input:       "aGVsbG8!@#",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesValue{}
			err := b.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("BytesValue.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("BytesValue.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if !bytes.Equal(b.Value, tt.expected) {
				t.Errorf("BytesValue.Set() = %v, want %v", b.Value, tt.expected)
			}
		})
	}
}

func TestBytesValue_Type(t *testing.T) {
	b := &types.BytesValue{}
	if got := b.Type(); got != "bytesBase64" {
		t.Errorf("BytesValue.Type() = %v, want %v", got, "bytesBase64")
	}
}

func TestBytesValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*types.BytesValue)(nil)
}

func TestBytesValue_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	b := &types.BytesValue{}

	fs.Var(b, "bytes", "Test bytes value")

	// Test valid input
	err := fs.Parse([]string{"--bytes", "aGVsbG8gd29ybGQ="})
	if err != nil {
		t.Fatalf("Failed to parse bytes: %v", err)
	}

	expected := []byte("hello world")
	if !bytes.Equal(b.Value, expected) {
		t.Errorf("Parsed bytes = %v, want %v", b.Value, expected)
	}

	// Test String() method after Set()
	if got := b.String(); got != "aGVsbG8gd29ybGQ=" {
		t.Errorf("BytesValue.String() after Set() = %v, want %v", got, "aGVsbG8gd29ybGQ=")
	}
}

func TestBytesValue_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	b := &types.BytesValue{}

	fs.Var(b, "bytes", "Test bytes value")

	// Test invalid input
	err := fs.Parse([]string{"--bytes", "invalid base64"})
	if err == nil {
		t.Fatalf("Expected error for invalid base64 input, got nil")
	}
}

func TestBytes_Constructor(t *testing.T) {
	wrap := &wrapperspb.BytesValue{Value: []byte("hello world")}
	b := types.Bytes(wrap)

	// Test that the returned BytesValue has correct fields
	if (*wrapperspb.BytesValue)(b) != wrap {
		t.Errorf("Bytes() wrapper = %p, want %p", (*wrapperspb.BytesValue)(b), wrap)
	}

	// Test that the value is correct
	if !bytes.Equal(b.Value, []byte("hello world")) {
		t.Errorf("Bytes().Value = %v, want %v", b.Value, []byte("hello world"))
	}

	// Test String() method on constructed bytes
	expected := "aGVsbG8gd29ybGQ="
	if got := b.String(); got != expected {
		t.Errorf("Bytes().String() = %v, want %v", got, expected)
	}
}

func TestBytes_ConstructorNil(t *testing.T) {
	b := types.Bytes(nil)

	// Test that it works with nil wrapper
	if b != nil {
		t.Errorf("Bytes(nil) = %p, want nil", b)
	}
}

func TestBytesValue_StringAfterSet(t *testing.T) {
	b := &types.BytesValue{}

	// Set a value
	err := b.Set("dGVzdCB2YWx1ZQ==")
	if err != nil {
		t.Fatalf("Failed to set bytes value: %v", err)
	}

	// Test String() method
	if got := b.String(); got != "dGVzdCB2YWx1ZQ==" {
		t.Errorf("BytesValue.String() after Set() = %v, want %v", got, "dGVzdCB2YWx1ZQ==")
	}
}

func TestBytesValue_MultipleSets(t *testing.T) {
	b := &types.BytesValue{}

	// Set first value
	err := b.Set("Zmlyc3Q=")
	if err != nil {
		t.Fatalf("Failed to set first bytes value: %v", err)
	}
	if !bytes.Equal(b.Value, []byte("first")) {
		t.Errorf("First set failed: got %v, want %v", b.Value, []byte("first"))
	}

	// Set second value
	err = b.Set("c2Vjb25k")
	if err != nil {
		t.Fatalf("Failed to set second bytes value: %v", err)
	}
	if !bytes.Equal(b.Value, []byte("second")) {
		t.Errorf("Second set failed: got %v, want %v", b.Value, []byte("second"))
	}

	// Test final String() value
	if got := b.String(); got != "c2Vjb25k" {
		t.Errorf("Final String() = %v, want %v", got, "c2Vjb25k")
	}
}

func TestBytesValue_RoundTrip(t *testing.T) {
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
			b := &types.BytesValue{}

			// Convert to base64 and set
			base64Str := base64.StdEncoding.EncodeToString(original)
			err := b.Set(base64Str)
			if err != nil {
				t.Fatalf("Failed to set bytes value: %v", err)
			}

			// Convert back and compare
			if !bytes.Equal(b.Value, original) {
				t.Errorf("Round trip failed: got %v, want %v", b.Value, original)
			}
		})
	}
}

func TestBytesValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{"empty", "", []byte{}},
		{"whitespace", "   ", []byte{}},
		{"single byte", "YQ==", []byte("a")},
		{"null bytes", "AA==", []byte{0x00}},
		{"max byte", "/w==", []byte{0xFF}},
		{"mixed ASCII", "SGVsbG8gV29ybGQ=", []byte("Hello World")},
		{"binary with null", "AGZvbw==", []byte{0x00, 'f', 'o', 'o'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &types.BytesValue{}
			err := b.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set bytes value %q: %v", tt.input, err)
			}
			if !bytes.Equal(b.Value, tt.want) {
				t.Errorf("BytesValue edge case: input=%q, got=%v, want=%v", tt.input, b.Value, tt.want)
			}
		})
	}
}

func TestBytesValue_LongData(t *testing.T) {
	b := &types.BytesValue{}

	// Test with a large byte array
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Convert to base64 and set
	base64Str := base64.StdEncoding.EncodeToString(largeData)
	err := b.Set(base64Str)
	if err != nil {
		t.Fatalf("Failed to set large bytes value: %v", err)
	}

	if !bytes.Equal(b.Value, largeData) {
		t.Errorf("Large data round trip failed")
	}

	// Test String() method produces valid base64
	if got := b.String(); got != base64Str {
		t.Errorf("Large data String() mismatch")
	}
}
