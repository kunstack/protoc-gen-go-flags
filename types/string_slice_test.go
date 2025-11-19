package types_test

import (
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestStringSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.StringValue
		want  string
	}{
		{
			name:  "empty slice",
			input: []*wrapperspb.StringValue{},
			want:  "[]",
		},
		{
			name:  "single string",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello")},
			want:  "[hello]",
		},
		{
			name:  "multiple strings",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world"), wrapperspb.String("test")},
			want:  "[hello,world,test]",
		},
		{
			name:  "strings with special characters",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello@world"), wrapperspb.String("test#123")},
			want:  "[hello@world,test#123]",
		},
		{
			name:  "strings with unicode",
			input: []*wrapperspb.StringValue{wrapperspb.String("こんにちは"), wrapperspb.String("world")},
			want:  "[こんにちは,world]",
		},
		{
			name:  "empty strings",
			input: []*wrapperspb.StringValue{wrapperspb.String(""), wrapperspb.String("")},
			want:  "[,]",
		},
		{
			name:  "mixed content",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello123"), wrapperspb.String("world!@#"), wrapperspb.String("test")},
			want:  "[hello123,world!@#,test]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := types.StringSlice(&tt.input)
			if got := ss.String(); got != tt.want {
				t.Errorf("StringSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceValue_Type(t *testing.T) {
	flags := []*wrapperspb.StringValue{}
	ss := types.StringSlice(&flags)
	if got := ss.Type(); got != "stringSliceValue" {
		t.Errorf("StringSliceValue.Type() = %v, want %v", got, "stringSliceValue")
	}
}

func TestStringSliceValue_PflagInterface(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that StringSliceValue implements pflag.Value interface
	ss := types.StringSlice(&flags)
	fs.Var(ss, "strings", "test string slice")

	// Test setting flags
	args := []string{"--strings", "hello", "--strings", "world", "--strings", "test"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Check the values through the StringSliceValue interface
	if len(flags) != 3 {
		t.Errorf("Expected 3 values, got %d", len(flags))
	}

	expected := []string{"hello", "world", "test"}
	for i, exp := range expected {
		if (flags)[i].Value != exp {
			t.Errorf("Value[%d] = %v, want %v", i, (flags)[i].Value, exp)
		}
	}
}

func TestStringSliceFunction(t *testing.T) {
	flags := []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world")}
	ss := types.StringSlice(&flags)

	if ss == nil {
		t.Error("StringSlice() returned nil")
	}

	if len(flags) != 2 {
		t.Errorf("StringSlice() length = %v, want %v", len(flags), 2)
	}

	if (flags)[0].Value != "hello" || (flags)[1].Value != "world" {
		t.Errorf("Values not as expected: %v, %v", (flags)[0].Value, (flags)[1].Value)
	}
}

func TestStringSliceFunction_Nil(t *testing.T) {
	var flags []*wrapperspb.StringValue
	ss := types.StringSlice(&flags)

	if ss == nil {
		t.Error("StringSlice() returned nil for nil slice")
	}

	// Should be able to call Set on nil slice
	err := ss.Set("test")
	if err != nil {
		t.Errorf("StringSlice().Set() failed on nil slice: %v", err)
	}

	if len(flags) != 1 || (flags)[0].Value != "test" {
		t.Errorf("StringSlice() failed to initialize nil slice properly")
	}
}

func TestStringSliceValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"single char", "a", []string{"a"}},
		{"whitespace", "   ", []string{"   "}},
		{"tabs", "\t\t", []string{"\t\t"}},
		{"newlines", "\"hello\nworld\"", []string{"\"hello\nworld\""}},
		{"unicode", "こんにちは", []string{"こんにちは"}},
		{"emoji", "👋", []string{"👋"}},
		{"numbers", "12345", []string{"12345"}},
		{"mixed", "hello 123 🌍", []string{"hello 123 🌍"}},
		{"quotes", "\"hello\"", []string{"\"hello\""}},
		{"brackets", "[hello]", []string{"[hello]"}},
		{"url", "https://example.com", []string{"https://example.com"}},
		{"email", "test@example.com", []string{"test@example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)
			err := ss.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set string slice value %q: %v", tt.input, err)
			}

			if len(flags) != len(tt.want) {
				t.Errorf("StringSlice edge case: input=%q, got length=%d, want length=%d",
					tt.input, len(flags), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if (flags)[i].Value != want {
					t.Errorf("StringSlice edge case: input=%q, got[%d]=%v, want=%v",
						tt.input, i, (flags)[i].Value, want)
				}
			}
		})
	}
}

func TestStringSliceValue_ChangedFlagBehavior(t *testing.T) {
	// Test the changed flag behavior indirectly through Set() method
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// First Set call should initialize the slice
	err := ss.Set("hello")
	if err != nil {
		t.Fatalf("First Set() failed: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("After first Set() expected 1 value, got %d", len(flags))
	}

	// Second Set call should append to the slice
	err = ss.Set("world")
	if err != nil {
		t.Fatalf("Second Set() failed: %v", err)
	}

	if len(flags) != 2 {
		t.Errorf("After second Set() expected 2 values, got %d", len(flags))
	}

	// Verify both values are present
	if (flags)[0].Value != "hello" || (flags)[1].Value != "world" {
		t.Errorf("Values not as expected: got %v and %v, want hello and world",
			(flags)[0].Value, (flags)[1].Value)
	}
}

func TestStringSliceValue_LongStrings(t *testing.T) {
	// Test with very long strings
	longString := string(make([]byte, 1000))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	err := ss.Set(longString)
	if err != nil {
		t.Fatalf("Failed to set very long string: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("Expected 1 value after setting long string, got %d", len(flags))
	}

	if (flags)[0].Value != longString {
		t.Errorf("Long string value mismatch")
	}
}

func TestStringSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		append  string
		want    []string
	}{
		{
			name:    "append to empty slice",
			initial: []string{},
			append:  "hello",
			want:    []string{"hello"},
		},
		{
			name:    "append to existing slice",
			initial: []string{"hello", "world"},
			append:  "test",
			want:    []string{"hello", "world", "test"},
		},
		{
			name:    "append string with spaces",
			initial: []string{"hello"},
			append:  "hello world",
			want:    []string{"hello", "hello world"},
		},
		{
			name:    "append string with special characters",
			initial: []string{"hello"},
			append:  "hello@world.com",
			want:    []string{"hello", "hello@world.com"},
		},
		{
			name:    "append unicode string",
			initial: []string{"hello"},
			append:  "こんにちは",
			want:    []string{"hello", "こんにちは"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			flags := make([]*wrapperspb.StringValue, len(tt.initial))
			for i, v := range tt.initial {
				flags[i] = wrapperspb.String(v)
			}
			ss := types.StringSlice(&flags)

			// Test Append
			err := ss.Append(tt.append)
			if err != nil {
				t.Fatalf("Append() failed: %v", err)
			}

			// Verify result
			if len(flags) != len(tt.want) {
				t.Errorf("Append() length = %v, want %v", len(flags), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if flags[i].Value != want {
					t.Errorf("Append()[%d] = %v, want %v", i, flags[i].Value, want)
				}
			}
		})
	}
}

func TestStringSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		replace []string
		want    []string
	}{
		{
			name:    "replace empty with values",
			initial: []string{},
			replace: []string{"hello", "world"},
			want:    []string{"hello", "world"},
		},
		{
			name:    "replace existing values",
			initial: []string{"old1", "old2"},
			replace: []string{"new1", "new2", "new3"},
			want:    []string{"new1", "new2", "new3"},
		},
		{
			name:    "replace with empty slice",
			initial: []string{"hello", "world"},
			replace: []string{},
			want:    []string{},
		},
		{
			name:    "replace with single value",
			initial: []string{"hello", "world", "test"},
			replace: []string{"single"},
			want:    []string{"single"},
		},
		{
			name:    "replace with special characters",
			initial: []string{"hello"},
			replace: []string{"hello@world", "test#123", "unicode: こんにちは"},
			want:    []string{"hello@world", "test#123", "unicode: こんにちは"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			flags := make([]*wrapperspb.StringValue, len(tt.initial))
			for i, v := range tt.initial {
				flags[i] = wrapperspb.String(v)
			}
			ss := types.StringSlice(&flags)

			// Test Replace
			err := ss.Replace(tt.replace)
			if err != nil {
				t.Fatalf("Replace() failed: %v", err)
			}

			// Verify result
			if len(flags) != len(tt.want) {
				t.Errorf("Replace() length = %v, want %v", len(flags), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if flags[i].Value != want {
					t.Errorf("Replace()[%d] = %v, want %v", i, flags[i].Value, want)
				}
			}
		})
	}
}

func TestStringSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func() []*wrapperspb.StringValue
		want  []string
	}{
		{
			name:  "get from empty slice",
			setup: func() []*wrapperspb.StringValue { return []*wrapperspb.StringValue{} },
			want:  []string{},
		},
		{
			name: "get from single value slice",
			setup: func() []*wrapperspb.StringValue {
				return []*wrapperspb.StringValue{wrapperspb.String("hello")}
			},
			want: []string{"hello"},
		},
		{
			name: "get from multiple values slice",
			setup: func() []*wrapperspb.StringValue {
				return []*wrapperspb.StringValue{
					wrapperspb.String("hello"),
					wrapperspb.String("world"),
					wrapperspb.String("test"),
				}
			},
			want: []string{"hello", "world", "test"},
		},
		{
			name: "get with special characters",
			setup: func() []*wrapperspb.StringValue {
				return []*wrapperspb.StringValue{
					wrapperspb.String("hello@world.com"),
					wrapperspb.String("test#123"),
					wrapperspb.String("path/to/file"),
				}
			},
			want: []string{"hello@world.com", "test#123", "path/to/file"},
		},
		{
			name: "get with unicode and emoji",
			setup: func() []*wrapperspb.StringValue {
				return []*wrapperspb.StringValue{
					wrapperspb.String("こんにちは"),
					wrapperspb.String("🌍"),
					wrapperspb.String("hello"),
				}
			},
			want: []string{"こんにちは", "🌍", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := tt.setup()
			ss := types.StringSlice(&flags)

			// Test GetSlice
			got := ss.GetSlice()

			if len(got) != len(tt.want) {
				t.Errorf("GetSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("GetSlice()[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}

func TestStringSliceValue_SequentialOperations(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// Test sequential operations that simulate real command-line usage
	// Note: Set() behavior: first call initializes, subsequent calls append
	operations := []struct {
		name     string
		fn       func() error
		expected []string
	}{
		{"Set initial value", func() error { return ss.Set("initial") }, []string{"initial"}},
		{"Append value 1", func() error { return ss.Append("appended1") }, []string{"initial", "appended1"}},
		{"Append value 2", func() error { return ss.Append("appended2") }, []string{"initial", "appended1", "appended2"}},
		{"Set new value (appends)", func() error { return ss.Set("new_value") }, []string{"initial", "appended1", "appended2", "new_value"}},
		{"Append value 3", func() error { return ss.Append("appended3") }, []string{"initial", "appended1", "appended2", "new_value", "appended3"}},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			if err := op.fn(); err != nil {
				t.Fatalf("%s failed: %v", op.name, err)
			}

			// Verify state after each operation
			if len(flags) != len(op.expected) {
				t.Errorf("After %s: expected %d elements, got %d", op.name, len(op.expected), len(flags))
			}
			for i, want := range op.expected {
				if flags[i].Value != want {
					t.Errorf("After %s: Element[%d] = %v, want %v", op.name, i, flags[i].Value, want)
				}
			}
		})
	}
}

func TestStringSliceValue_CSVSpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.StringValue
		want  string
	}{
		{
			name: "strings with commas",
			input: []*wrapperspb.StringValue{
				wrapperspb.String("hello, world"),
				wrapperspb.String("test,123"),
			},
			want: "[\"hello, world\",\"test,123\"]",
		},
		{
			name: "strings with quotes",
			input: []*wrapperspb.StringValue{
				wrapperspb.String(`say "hello"`),
				wrapperspb.String(`it's "test"`),
			},
			want: `["say ""hello""","it's ""test"""]`,
		},
		{
			name: "strings with newlines",
			input: []*wrapperspb.StringValue{
				wrapperspb.String("hello\nworld"),
				wrapperspb.String("test\n123"),
			},
			want: "[\"hello\nworld\",\"test\n123\"]",
		},
		{
			name: "mixed special characters",
			input: []*wrapperspb.StringValue{
				wrapperspb.String(`hello, "world"`),
				wrapperspb.String("test\n123, abc"),
			},
			want: `["hello, ""world""","test
123, abc"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := types.StringSlice(&tt.input)
			if got := ss.String(); got != tt.want {
				t.Errorf("StringSliceValue.String() with CSV special chars = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceValue_MassiveData(t *testing.T) {
	// Test with massive amount of data
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// Add 1000 strings
	for i := 0; i < 1000; i++ {
		value := "string_" + string(rune('0'+i%10)) + "_" + string(rune('a'+i%26))
		if err := ss.Set(value); err != nil {
			t.Fatalf("Failed to set string %d: %v", i, err)
		}
	}

	if len(flags) != 1000 {
		t.Errorf("Expected 1000 elements, got %d", len(flags))
	}

	// Test String() with massive data
	result := ss.String()
	if len(result) == 0 {
		t.Error("String() returned empty result for massive data")
	}

	// Test GetSlice() with massive data
	slice := ss.GetSlice()
	if len(slice) != 1000 {
		t.Errorf("GetSlice() returned %d elements, want 1000", len(slice))
	}
}

func TestStringSliceValue_ErrorRecovery(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// Set initial values
	if err := ss.Set("hello"); err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}
	if err := ss.Set("world"); err != nil {
		t.Fatalf("Failed to set second value: %v", err)
	}

	// Test that operations maintain state consistency
	_ = len(flags) // original length for reference

	// Test Replace with empty slice
	if err := ss.Replace([]string{}); err != nil {
		t.Fatalf("Replace with empty slice failed: %v", err)
	}
	if len(flags) != 0 {
		t.Errorf("Expected empty slice after Replace, got %d elements", len(flags))
	}

	// Test that we can recover and add new values
	if err := ss.Set("new1"); err != nil {
		t.Fatalf("Failed to set after empty replace: %v", err)
	}
	if err := ss.Set("new2"); err != nil {
		t.Fatalf("Failed to set second value after empty replace: %v", err)
	}
	if len(flags) != 2 {
		t.Errorf("Expected 2 elements after recovery, got %d", len(flags))
	}
}

func TestStringSliceValue_NilElements(t *testing.T) {
	// Test behavior with nil elements in the slice - this tests edge case handling
	// Note: nil elements should not normally occur, but we test graceful handling
	flags := []*wrapperspb.StringValue{
		wrapperspb.String("hello"),
		{}, // Empty StringValue (not nil, but zero value)
		wrapperspb.String("world"),
	}
	ss := types.StringSlice(&flags)

	// Test String() with empty elements
	result := ss.String()
	if result == "" {
		t.Error("String() should handle empty elements gracefully")
	}

	// Test GetSlice() with empty elements
	slice := ss.GetSlice()
	if len(slice) != 3 {
		t.Errorf("GetSlice() should return all elements, got %d", len(slice))
	}
}

func TestStringSliceValue_SpecialValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "JSON string",
			input: `{"key": "value", "number": 123}`,
			want:  `{"key": "value", "number": 123}`,
		},
		{
			name:  "XML string",
			input: `<root><child>value</child></root>`,
			want:  `<root><child>value</child></root>`,
		},
		{
			name:  "SQL query",
			input: `SELECT * FROM users WHERE id = '123' AND name = "John"`,
			want:  `SELECT * FROM users WHERE id = '123' AND name = "John"`,
		},
		{
			name:  "file path",
			input: `/usr/local/bin/program --arg "value with spaces"`,
			want:  `/usr/local/bin/program --arg "value with spaces"`,
		},
		{
			name:  "regex pattern",
			input: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			want:  `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		},
		{
			name:  "base64 string",
			input: `SGVsbG8gV29ybGQhISE=`,
			want:  `SGVsbG8gV29ybGQhISE=`,
		},
		{
			name:  "UUID",
			input: `550e8400-e29b-41d4-a716-446655440000`,
			want:  `550e8400-e29b-41d4-a716-446655440000`,
		},
		{
			name:  "IPv6 address",
			input: `2001:0db8:85a3:0000:0000:8a2e:0370:7334`,
			want:  `2001:0db8:85a3:0000:0000:8a2e:0370:7334`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)

			err := ss.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set special value %q: %v", tt.input, err)
			}

			if len(flags) != 1 {
				t.Errorf("Expected 1 value, got %d", len(flags))
				return
			}

			if flags[0].Value != tt.want {
				t.Errorf("Special value test: input=%q, got=%q, want=%q",
					tt.input, flags[0].Value, tt.want)
			}

			// Test round-trip through String() and GetSlice()
			_ = ss.String()
			slice := ss.GetSlice()

			if len(slice) != 1 {
				t.Errorf("GetSlice() returned %d elements, want 1", len(slice))
			}

			if slice[0] != tt.want {
				t.Errorf("Round-trip failed: got=%q, want=%q", slice[0], tt.want)
			}
		})
	}
}

func TestStringSliceValue_CommandLineScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario []string
		want     []string
	}{
		{
			name:     "typical flag usage",
			scenario: []string{"web,api", "database", "cache"},
			want:     []string{"web,api", "database", "cache"},
		},
		{
			name:     "environment variables style",
			scenario: []string{"PATH=/usr/bin", "PATH=/usr/local/bin", "PATH=/opt/bin"},
			want:     []string{"PATH=/usr/bin", "PATH=/usr/local/bin", "PATH=/opt/bin"},
		},
		{
			name:     "file paths with spaces",
			scenario: []string{"/home/user/my file.txt", "/home/user/another file.log"},
			want:     []string{"/home/user/my file.txt", "/home/user/another file.log"},
		},
		{
			name:     "URLs as parameters",
			scenario: []string{"https://api.example.com/v1", "https://backup.example.com/v2"},
			want:     []string{"https://api.example.com/v1", "https://backup.example.com/v2"},
		},
		{
			name:     "JSON strings",
			scenario: []string{`{"key": "value"}`, `{"nested": {"data": true}}`},
			want:     []string{`{"key": "value"}`, `{"nested": {"data": true}}`},
		},
		{
			name:     "SQL queries",
			scenario: []string{"SELECT * FROM users", "SELECT id, name FROM products"},
			want:     []string{"SELECT * FROM users", "SELECT id, name FROM products"},
		},
		{
			name:     "regex patterns",
			scenario: []string{"^[a-zA-Z]+$", "\\d{3}-\\d{3}-\\d{4}"},
			want:     []string{"^[a-zA-Z]+$", "\\d{3}-\\d{3}-\\d{4}"},
		},
		{
			name:     "base64 encoded data",
			scenario: []string{"SGVsbG8gV29ybGQ=", "VGVzdCBTdHJpbmc="},
			want:     []string{"SGVsbG8gV29ybGQ=", "VGVzdCBTdHJpbmc="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)

			// Simulate command-line argument processing
			for _, arg := range tt.scenario {
				if err := ss.Set(arg); err != nil {
					t.Fatalf("Failed to set argument %q: %v", arg, err)
				}
			}

			// Verify results
			if len(flags) != len(tt.want) {
				t.Errorf("Expected %d elements, got %d", len(tt.want), len(flags))
				return
			}

			for i, want := range tt.want {
				if flags[i].Value != want {
					t.Errorf("Element[%d] = %q, want %q", i, flags[i].Value, want)
				}
			}

			// Test String() representation
			str := ss.String()
			if str == "" {
				t.Error("String() returned empty result")
			}

			// Test GetSlice() functionality
			slice := ss.GetSlice()
			if len(slice) != len(tt.want) {
				t.Errorf("GetSlice() returned %d elements, want %d", len(slice), len(tt.want))
			}
		})
	}
}

func TestStringSliceValue_MultipleOperations(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// Test sequence: Set -> Append -> Replace -> GetSlice -> String

	// 1. Set initial values
	if err := ss.Set("hello"); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}
	if err := ss.Set("world"); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// 2. Append more values
	if err := ss.Append("test"); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}
	if err := ss.Append("123"); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	// 3. Replace with new values
	if err := ss.Replace([]string{"new1", "new2"}); err != nil {
		t.Fatalf("Replace() failed: %v", err)
	}

	// 4. GetSlice and verify
	slice := ss.GetSlice()
	expected := []string{"new1", "new2"}
	if len(slice) != len(expected) {
		t.Errorf("GetSlice() length = %v, want %v", len(slice), len(expected))
	}
	for i, want := range expected {
		if slice[i] != want {
			t.Errorf("GetSlice()[%d] = %v, want %v", i, slice[i], want)
		}
	}

	// 5. String and verify
	str := ss.String()
	if str != "[new1,new2]" {
		t.Errorf("String() = %v, want %v", str, "[new1,new2]")
	}

	// 6. Type and verify
	typ := ss.Type()
	if typ != "stringSliceValue" {
		t.Errorf("Type() = %v, want %v", typ, "stringSliceValue")
	}
}
