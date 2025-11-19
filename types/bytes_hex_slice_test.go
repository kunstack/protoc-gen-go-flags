package types_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBytesHexSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:  "single hex string",
			input: "48656C6C6F20576F726C64", // "Hello World" in hex
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
		},
		{
			name:    "invalid hex string",
			input:   "invalid-hex!",
			wantErr: true,
		},
		{
			name:  "binary data",
			input: "000102030405060708090A0B", // binary data in hex
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
		},
		{
			name:    "comma separated should fail",
			input:   "48656c6c6f,576f726c64",
			wantErr: true,
		},
		{
			name:    "quoted comma separated should fail",
			input:   `"48656C6C6F","576F726C64"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesHexSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("BytesHexSliceValue.Set() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((bytesSlice)[i].Value) != string(wantVal.Value) {
						t.Errorf("BytesHexSliceValue.Set()[%d] = %v, want %v", i, (bytesSlice)[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesHexSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.BytesValue
		want  string
	}{
		{
			name:  "empty slice",
			input: make([]*wrapperspb.BytesValue, 0),
			want:  "[]",
		},
		{
			name:  "single hex string",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
			want:  "[48656C6C6F20576F726C64]",
		},
		{
			name:  "multiple values",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
			want:  "[48656C6C6F,576F726C64]", // CSV format output as per implementation
		},
		{
			name:  "binary data",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
			want:  "[000102030405060708090A0B]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := types.BytesHexSlice(&tt.input)
			if got := bs.String(); got != tt.want {
				t.Errorf("BytesHexSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexSliceValue_Type(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)
	if got := bs.Type(); got != "bytesHexSliceValue" {
		t.Errorf("BytesHexSliceValue.Type() = %v, want bytesHexSliceValue", got)
	}
}

func TestBytesHexSliceValue_ImplementsPflagValue(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)
	var _ = bs
}

func TestBytesHexSliceValue_MultipleSets(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	// First set - single value
	err := bs.Set("48656C6C6F") // "Hello" in hex
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set - another single value
	err = bs.Set("576F726C64") // "World" in hex
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Third set - comma separated should fail
	err = bs.Set("54686973204973,412054657374") // "This Is, A Test" in hex - should fail
	if err == nil {
		t.Fatalf("Expected error for comma-separated values, but got none")
	}

	// Verify only the first two values are present
	expected := []string{"Hello", "World"}
	if len(bytesSlice) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(bytesSlice))
	}

	for i, expectedStr := range expected {
		if string(bytesSlice[i].Value) != expectedStr {
			t.Errorf("Expected %s at index %d, got %s", expectedStr, i, bytesSlice[i].Value)
		}
	}
}

func TestBytesHexSliceValue_SpecialCharacters(t *testing.T) {
	// Test with special characters that need hex encoding
	specialBytes := []byte{0, 1, 2, 255, 254, 253}
	encoded := hex.EncodeToString(specialBytes)

	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	err := bs.Set(encoded)
	if err != nil {
		t.Fatalf("Set() with special characters error = %v", err)
	}

	if len(bytesSlice) != 1 {
		t.Fatalf("Expected 1 value, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0].Value) != string(specialBytes) {
		t.Errorf("Expected %v, got %v", specialBytes, bytesSlice[0].Value)
	}

	// Check that String() method encodes back correctly (should be uppercase)
	expectedString := "[" + strings.ToUpper(encoded) + "]"
	if bs.String() != expectedString {
		t.Errorf("String() = %v, want %v", bs.String(), expectedString)
	}
}

func TestBytesHexSlice_Constructor(t *testing.T) {
	// Test with non-nil slice
	initialBytes := []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("test"))}
	result := types.BytesHexSlice(&initialBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Test that it implements pflag.Value
	var _ = result

	// Test that the initial value is preserved
	if result.String() != "[74657374]" {
		t.Errorf("Expected [54455354], got %s", result.String())
	}
}

func TestBytesHexSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice
	var nilBytes []*wrapperspb.BytesValue
	result := types.BytesHexSlice(&nilBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Test that it implements pflag.Value
	var _ = result

	// Test that it works with nil input
	if result.String() != "[]" {
		t.Errorf("Expected '[]', got %s", result.String())
	}
}

func TestBytesHexSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		append  string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:    "append to empty slice",
			initial: "",
			append:  "48656C6C6F", // "Hello" in hex
			want:    []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello"))},
			wantErr: false,
		},
		{
			name:    "append to existing slice",
			initial: "48656C6C6F", // "Hello" in hex
			append:  "576F726C64", // "World" in hex
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
		{
			name:    "append invalid hex",
			initial: "48656C6C6F", // "Hello" in hex
			append:  "invalid-hex!",
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
			},
			wantErr: true,
		},
		{
			name:    "append with whitespace",
			initial: "48656C6C6F",     // "Hello" in hex
			append:  "  576F726C64  ", // "World" in hex with spaces
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)
			bsTyped := bs.(*types.BytesHexSliceValue)

			// Set initial value if provided
			if tt.initial != "" {
				if err := bs.Set(tt.initial); err != nil {
					t.Fatalf("Failed to set initial value: %v", err)
				}
			}

			// Test Append
			err := bsTyped.Append(tt.append)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("Append() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i].Value) != string(wantVal.Value) {
						t.Errorf("Append()[%d] = %v, want %v", i, bytesSlice[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesHexSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		replace []string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:    "replace empty with values",
			initial: []string{},
			replace: []string{"48656C6C6F", "576F726C64"}, // "Hello", "World"
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
		{
			name:    "replace existing values",
			initial: []string{"48656C6C6F"},             // "Hello"
			replace: []string{"576F726C64", "54657374"}, // "World", "Test"
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("World")),
				wrapperspb.Bytes([]byte("Test")),
			},
			wantErr: false,
		},
		{
			name:    "replace with invalid hex",
			initial: []string{"48656C6C6F"}, // "Hello"
			replace: []string{"invalid-hex!"},
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")), // Should remain unchanged on error
			},
			wantErr: true,
		},
		{
			name:    "replace with whitespace",
			initial: []string{"48656C6C6F"},                     // "Hello"
			replace: []string{"  576F726C64  ", "  54657374  "}, // "World", "Test" with spaces
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("World")),
				wrapperspb.Bytes([]byte("Test")),
			},
			wantErr: false,
		},
		{
			name:    "replace with empty slice",
			initial: []string{"48656C6C6F"}, // "Hello"
			replace: []string{},
			want:    []*wrapperspb.BytesValue{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)

			// Set initial values if provided
			for _, val := range tt.initial {
				if err := bs.Set(val); err != nil {
					t.Fatalf("Failed to set initial value: %v", err)
				}
			}

			// Test Replace
			bsTyped := bs.(*types.BytesHexSliceValue)
			err := bsTyped.Replace(tt.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("Replace() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i].Value) != string(wantVal.Value) {
						t.Errorf("Replace()[%d] = %v, want %v", i, bytesSlice[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesHexSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*types.BytesHexSliceValue)
		want  []string
	}{
		{
			name: "get slice from empty",
			setup: func(bs *types.BytesHexSliceValue) {
				// Leave empty
			},
			want: []string{},
		},
		{
			name: "get slice with values",
			setup: func(bs *types.BytesHexSliceValue) {
				bs.Set("48656C6C6F") // "Hello"
				bs.Set("576F726C64") // "World"
			},
			want: []string{"48656C6C6F", "576F726C64"},
		},
		{
			name: "get slice with special bytes",
			setup: func(bs *types.BytesHexSliceValue) {
				bs.Set("000102") // Binary data
				bs.Set("FFFE")   // High bytes
			},
			want: []string{"000102", "FFFE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)

			// Setup the test scenario
			bsTyped := bs.(*types.BytesHexSliceValue)
			tt.setup(bsTyped)

			// Test GetSlice
			got := bsTyped.GetSlice()
			if len(got) != len(tt.want) {
				t.Errorf("GetSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, wantVal := range tt.want {
				if got[i] != wantVal {
					t.Errorf("GetSlice()[%d] = %v, want %v", i, got[i], wantVal)
				}
			}
		})
	}
}

func TestBytesHexSliceValue_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "odd length hex should fail",
			input:   "48656C6", // 7 characters (odd length)
			wantErr: true,
		},
		{
			name:    "single byte hex",
			input:   "FF", // Single byte
			wantErr: false,
		},
		{
			name:    "max byte values",
			input:   "FFFFFFFF", // All 0xFF bytes
			wantErr: false,
		},
		{
			name:    "mixed case hex",
			input:   "48656c6C6f20576f726C64", // Mixed case
			wantErr: false,
		},
		{
			name:    "hex with leading/trailing spaces",
			input:   "  48656C6C6F20576F726C64  ", // Spaces around
			wantErr: false,
		},
		{
			name:    "hex with invalid characters",
			input:   "48656C6C6G", // Contains 'G'
			wantErr: true,
		},
		{
			name:    "hex with symbols",
			input:   "48656C6C6!", // Contains '!'
			wantErr: true,
		},
		{
			name:    "hex with 0x prefix should fail",
			input:   "0x48656C6C6F", // 0x prefix
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)

			err := bs.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytesHexSliceValue_LargeData(t *testing.T) {
	// Test with large data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	hexStr := hex.EncodeToString(largeData)
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	err := bs.Set(hexStr)
	if err != nil {
		t.Fatalf("Set() with large data failed: %v", err)
	}

	if len(bytesSlice) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(bytesSlice))
	}

	if len(bytesSlice[0].Value) != len(largeData) {
		t.Errorf("Data length mismatch: got %d, want %d", len(bytesSlice[0].Value), len(largeData))
	}

	// Verify data integrity
	for i := 0; i < len(largeData); i++ {
		if bytesSlice[0].Value[i] != largeData[i] {
			t.Errorf("Data mismatch at index %d: got %v, want %v", i, bytesSlice[0].Value[i], largeData[i])
			break
		}
	}
}

func TestBytesHexNativeSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		append  string
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "append to empty slice",
			initial: "",
			append:  "48656c6c6f", // "Hello" in hex
			want:    [][]byte{[]byte("Hello")},
			wantErr: false,
		},
		{
			name:    "append to existing slice",
			initial: "48656c6c6f", // "Hello" in hex
			append:  "576f726c64", // "World" in hex
			want: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
			},
			wantErr: false,
		},
		{
			name:    "append invalid hex",
			initial: "48656c6c6f", // "Hello" in hex
			append:  "invalid-hex!",
			want: [][]byte{
				[]byte("Hello"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesSlice [][]byte
			nbhs := types.BytesHexSlice(&bytesSlice)

			// Set initial value if provided
			if tt.initial != "" {
				if err := nbhs.Set(tt.initial); err != nil {
					t.Fatalf("Failed to set initial value: %v", err)
				}
			}

			// Test Append
			nbhsTyped := nbhs.(*types.BytesHexNativeSliceValue)
			err := nbhsTyped.Append(tt.append)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("Append() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i]) != string(wantVal) {
						t.Errorf("Append()[%d] = %v, want %v", i, bytesSlice[i], wantVal)
					}
				}
			}
		})
	}
}

func TestBytesHexNativeSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		replace []string
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "replace empty with values",
			initial: []string{},
			replace: []string{"48656c6c6f", "576f726c64"}, // "Hello", "World"
			want: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
			},
			wantErr: false,
		},
		{
			name:    "replace existing values",
			initial: []string{"48656c6c6f"},             // "Hello"
			replace: []string{"576f726c64", "54657374"}, // "World", "Test"
			want: [][]byte{
				[]byte("World"),
				[]byte("Test"),
			},
			wantErr: false,
		},
		{
			name:    "replace with invalid hex",
			initial: []string{"48656c6c6f"}, // "Hello"
			replace: []string{"invalid-hex!"},
			want: [][]byte{
				[]byte("Hello"), // Should remain unchanged on error
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesSlice [][]byte
			nbhs := types.BytesHexSlice(&bytesSlice)

			// Set initial values if provided
			for _, val := range tt.initial {
				if err := nbhs.Set(val); err != nil {
					t.Fatalf("Failed to set initial value: %v", err)
				}
			}

			// Test Replace
			nbhsTyped := nbhs.(*types.BytesHexNativeSliceValue)
			err := nbhsTyped.Replace(tt.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("Replace() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i]) != string(wantVal) {
						t.Errorf("Replace()[%d] = %v, want %v", i, bytesSlice[i], wantVal)
					}
				}
			}
		})
	}
}

func TestBytesHexNativeSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*types.BytesHexNativeSliceValue)
		want  []string
	}{
		{
			name: "get slice from empty",
			setup: func(bs *types.BytesHexNativeSliceValue) {
				// Leave empty
			},
			want: []string{},
		},
		{
			name: "get slice with values",
			setup: func(bs *types.BytesHexNativeSliceValue) {
				bs.Set("48656c6c6f") // "Hello"
				bs.Set("576f726c64") // "World"
			},
			want: []string{"48656C6C6F", "576F726C64"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesSlice [][]byte
			nbhs := types.BytesHexSlice(&bytesSlice)

			// Setup the test scenario
			tt.setup(nbhs.(*types.BytesHexNativeSliceValue))

			// Test GetSlice
			nbhsTyped := nbhs.(*types.BytesHexNativeSliceValue)
			got := nbhsTyped.GetSlice()
			if len(got) != len(tt.want) {
				t.Errorf("GetSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, wantVal := range tt.want {
				if got[i] != wantVal {
					t.Errorf("GetSlice()[%d] = %v, want %v", i, got[i], wantVal)
				}
			}
		})
	}
}

func TestBytesHexSliceValue_ErrorRecovery(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	// Set a valid value first
	if err := bs.Set("48656C6C6F"); err != nil { // "Hello"
		t.Fatalf("Failed to set valid value: %v", err)
	}

	// Try to set an invalid value
	if err := bs.Set("invalid-hex!"); err == nil {
		t.Error("Expected error for invalid hex, but got none")
	}

	// Verify the original value is still there
	if len(bytesSlice) != 1 {
		t.Errorf("Expected 1 element after failed Set, got %d", len(bytesSlice))
	}
	if string(bytesSlice[0].Value) != "Hello" {
		t.Errorf("Expected 'Hello' to remain after failed Set, got %s", bytesSlice[0].Value)
	}

	// Test Replace with partial failure
	bsTyped := bs.(*types.BytesHexSliceValue)
	if err := bsTyped.Replace([]string{"576F726C64", "invalid-hex!"}); err == nil {
		t.Error("Expected error for invalid hex in Replace, but got none")
	}

	// Verify state is consistent after failed Replace
	if len(bytesSlice) != 1 {
		t.Errorf("Expected 1 element after failed Replace, got %d", len(bytesSlice))
	}
}

func TestBytesHexSliceValue_Constructor(t *testing.T) {
	tests := []struct {
		name           string
		initial        []*wrapperspb.BytesValue
		expectedType   string
		expectedString string
	}{
		{
			name:           "constructor with empty slice",
			initial:        []*wrapperspb.BytesValue{},
			expectedType:   "bytesHexSliceValue",
			expectedString: "[]",
		},
		{
			name:           "constructor with values",
			initial:        []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("test"))},
			expectedType:   "bytesHexSliceValue",
			expectedString: "[74657374]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.BytesHexSlice(&tt.initial)

			if result == nil {
				t.Fatal("BytesHexSlice() returned nil")
			}

			if got := result.Type(); got != tt.expectedType {
				t.Errorf("Type() = %v, want %v", got, tt.expectedType)
			}

			if got := result.String(); got != tt.expectedString {
				t.Errorf("String() = %v, want %v", got, tt.expectedString)
			}
		})
	}
}

func TestBytesHexNativeSliceValue_Constructor(t *testing.T) {
	tests := []struct {
		name           string
		initial        [][]byte
		expectedType   string
		expectedString string
	}{
		{
			name:           "constructor with empty slice",
			initial:        [][]byte{},
			expectedType:   "BytesHexNativeSliceValue",
			expectedString: "[]",
		},
		{
			name:           "constructor with values",
			initial:        [][]byte{[]byte("test")},
			expectedType:   "BytesHexNativeSliceValue",
			expectedString: "[74657374]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.BytesHexSlice(&tt.initial)

			if result == nil {
				t.Fatal("BytesHexSlice() returned nil")
			}

			if got := result.Type(); got != tt.expectedType {
				t.Errorf("Type() = %v, want %v", got, tt.expectedType)
			}

			if got := result.String(); got != tt.expectedString {
				t.Errorf("String() = %v, want %v", got, tt.expectedString)
			}
		})
	}
}

func TestBytesHexSliceValue_PflagInterface(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesHexSliceValue implements pflag.Value interface
	bs := types.BytesHexSlice(&bytesSlice)
	fs.Var(bs, "bytes", "test bytes hex slice")

	// Test setting single flag
	args := []string{"--bytes", "48656C6C6F"} // "Hello" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Test setting another single flag (should append)
	args = []string{"--bytes", "576F726C64"} // "World" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Test that comma-separated values fail
	args = []string{"--bytes", "54686973204973,412054657374"} // "This Is, A Test" in hex - should fail
	if err := fs.Parse(args); err == nil {
		t.Errorf("Expected error for comma-separated values, but got none")
	}

	// Verify only the first two values are present
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 values, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0].Value) != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", (bytesSlice)[0].Value)
	}

	if string(bytesSlice[1].Value) != "World" {
		t.Errorf("Expected 'World', got '%s'", (bytesSlice)[1].Value)
	}
}

// Tests for [][]byte (native byte slice)
func TestBytesHexNativeSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]byte
		wantErr bool
	}{
		{
			name:  "single hex",
			input: "48656c6c6f20576f726c64",
			want:  [][]byte{[]byte("Hello World")},
		},
		{
			name:    "comma separated should fail",
			input:   "48656c6c6f,576f726c64,54657374",
			wantErr: true,
		},
		{
			name:    "comma with spaces should fail",
			input:   " 48656c6c6f , 576f726c64 , 54657374 ",
			wantErr: true,
		},
		{
			name:    "invalid hex in comma list should fail",
			input:   "48656c6c6f,invalid,54657374",
			wantErr: true,
		},
		{
			name:    "quoted comma separated should fail",
			input:   `"48656c6c6f","576f726c64"`,
			wantErr: true,
		},
		{
			name:    "mixed case comma separated should fail",
			input:   "48656C6C6F,776f726c64",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesSlice [][]byte
			nbhs := types.BytesHexSlice(&bytesSlice)
			err := nbhs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesHexNativeSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("BytesHexNativeSliceValue.Set() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i]) != string(wantVal) {
						t.Errorf("BytesHexNativeSliceValue.Set()[%d] = %v, want %v", i, bytesSlice[i], wantVal)
					}
				}
			}
		})
	}
}

func TestBytesHexNativeSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input [][]byte
		want  string
	}{
		{
			name:  "empty slice",
			input: [][]byte{},
			want:  "[]",
		},
		{
			name:  "single bytes",
			input: [][]byte{[]byte("Hello World")},
			want:  "[48656c6c6f20576f726c64]",
		},
		{
			name:  "multiple values",
			input: [][]byte{[]byte("Hello"), []byte("World"), []byte("Test")},
			want:  "[48656c6c6f,576f726c64,54657374]", // CSV format output as per implementation
		},
		{
			name:  "binary data",
			input: [][]byte{{0x00, 0x01, 0x02, 0x03}, {0xFF, 0xFE, 0xFD}},
			want:  "[00010203,fffefd]", // CSV format output
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nbhs := types.BytesHexSlice(&tt.input)
			if got := nbhs.String(); got != tt.want {
				t.Errorf("BytesHexNativeSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexNativeSliceValue_Type(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)
	if got := nbhs.Type(); got != "BytesHexNativeSliceValue" {
		t.Errorf("BytesHexNativeSliceValue.Type() = %v, want BytesHexNativeSliceValue", got)
	}
}

func TestBytesHexNativeSliceValue_ImplementsPflagValue(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)
	var _ = nbhs
}

func TestBytesHexNativeSliceValue_MultipleSets(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)

	// First set - single value
	err := nbhs.Set("48656c6c6f") // "Hello" in hex
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set - another single value
	err = nbhs.Set("576f726c64") // "World" in hex
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Third set - comma separated should fail
	err = nbhs.Set("54657374") // "Test" in hex - single value
	if err != nil {
		t.Fatalf("Third Set() error = %v", err)
	}

	// Verify all values are present
	expected := []string{"Hello", "World", "Test"}
	if len(bytesSlice) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(bytesSlice))
	}

	for i, expectedVal := range expected {
		if string(bytesSlice[i]) != expectedVal {
			t.Errorf("Expected %v at index %d, got %v", expectedVal, i, bytesSlice[i])
		}
	}
}

func TestBytesHexNativeSliceValue_PflagInterface(t *testing.T) {
	var bytesSlice [][]byte
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesHexNativeSliceValue implements pflag.Value interface
	nbhs := types.BytesHexSlice(&bytesSlice)
	fs.Var(nbhs, "hex-bytes", "test native hex bytes slice")

	// Test setting single flag
	args := []string{"--hex-bytes", "48656c6c6f"} // "Hello" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Test setting another single flag (should append)
	args = []string{"--hex-bytes", "576f726c64"} // "World" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Test that comma-separated values fail
	args = []string{"--hex-bytes", "54657374,6d6f7265"} // "Test,more" in hex - should fail
	if err := fs.Parse(args); err == nil {
		t.Errorf("Expected error for comma-separated values, but got none")
	}

	// Verify only the first two values are present
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 values, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0]) != "Hello" {
		t.Errorf("Expected 'Hello' at index 0, got %v", bytesSlice[0])
	}

	if string(bytesSlice[1]) != "World" {
		t.Errorf("Expected 'World' at index 1, got %v", bytesSlice[1])
	}
}
