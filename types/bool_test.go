package types_test

import (
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBoolValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.BoolValue
		expected string
	}{
		{
			name:     "true value",
			value:    wrapperspb.Bool(true),
			expected: "true",
		},
		{
			name:     "false value",
			value:    wrapperspb.Bool(false),
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := types.Bool(tt.value)
			if got := b.String(); got != tt.expected {
				t.Errorf("BoolValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoolValue_Set(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue bool
		wantErr   bool
	}{
		{
			name:      "true lowercase",
			input:     "true",
			wantValue: true,
		},
		{
			name:      "true uppercase",
			input:     "TRUE",
			wantValue: true,
		},
		{
			name:      "true mixed case",
			input:     "True",
			wantValue: true,
		},
		{
			name:      "false lowercase",
			input:     "false",
			wantValue: false,
		},
		{
			name:      "false uppercase",
			input:     "FALSE",
			wantValue: false,
		},
		{
			name:      "false mixed case",
			input:     "False",
			wantValue: false,
		},
		{
			name:      "1",
			input:     "1",
			wantValue: true,
		},
		{
			name:      "0",
			input:     "0",
			wantValue: false,
		},
		{
			name:      "t",
			input:     "t",
			wantValue: true,
		},
		{
			name:      "f",
			input:     "f",
			wantValue: false,
		},
		{
			name:      "T",
			input:     "T",
			wantValue: true,
		},
		{
			name:      "F",
			input:     "F",
			wantValue: false,
		},
		{
			name:    "invalid value",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "2",
			input:   "2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b types.BoolValue
			err := b.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if b.Value != tt.wantValue {
					t.Errorf("BoolValue.Set() = %v, want %v", b.Value, tt.wantValue)
				}
			}
		})
	}
}

func TestBoolValue_Type(t *testing.T) {
	var b types.BoolValue
	if got := b.Type(); got != "boolValue" {
		t.Errorf("BoolValue.Type() = %v, want %v", got, "boolValue")
	}
}

func TestBoolValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*types.BoolValue)(nil)
}

func TestBoolValue_PflagIntegrationTrue(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var b types.BoolValue

	f := fs.VarPF(&b, "bool", "", "Test bool")
	f.NoOptDefVal = "true"

	err := fs.Parse([]string{"--bool", "true"})
	if err != nil {
		t.Fatalf("Failed to parse bool: %v", err)
	}

	if !b.Value {
		t.Errorf("Parsed bool = %v, want true", b.Value)
	}
}

func TestBoolValue_PflagIntegrationFalse(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var b types.BoolValue

	fs.Var(&b, "bool", "Test bool")

	err := fs.Parse([]string{"--bool", "false"})
	if err != nil {
		t.Fatalf("Failed to parse bool: %v", err)
	}

	if b.Value {
		t.Errorf("Parsed bool = %v, want false", b.Value)
	}

	fs2 := pflag.NewFlagSet("test", pflag.ContinueOnError)

	f := fs2.VarPF(&b, "bool", "", "Test bool")
	f.NoOptDefVal = "true"

	err = fs2.Parse([]string{"--bool"})
	if err != nil {
		t.Fatalf("Failed to parse bool: %v", err)
	}

	if !b.Value {
		t.Errorf("Parsed bool = %v, want true", b.Value)
	}
}

func TestBoolValue_PflagIntegrationNoOptDef(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var b types.BoolValue

	f := fs.VarPF(&b, "bool", "", "Test bool")
	f.NoOptDefVal = "true"

	err := fs.Parse([]string{"--bool"})
	if err != nil {
		t.Fatalf("Failed to parse bool: %v", err)
	}

	if !b.Value {
		t.Errorf("Parsed bool = %v, want true", b.Value)
	}
}

func TestBoolValue_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var b types.BoolValue

	fs.Var(&b, "bool", "Test bool")

	err := fs.Parse([]string{"--bool", "invalid"})
	if err == nil {
		t.Fatal("Expected error for invalid bool value")
	}
}

func TestBool_Constructor(t *testing.T) {
	// Test with true value
	trueWrap := wrapperspb.Bool(true)
	b := types.Bool(trueWrap)

	if b.Value != true {
		t.Errorf("Bool(true) = %v, want true", b.Value)
	}

	// Test with false value
	falseWrap := wrapperspb.Bool(false)
	b = types.Bool(falseWrap)

	if b.Value != false {
		t.Errorf("Bool(false) = %v, want false", b.Value)
	}
}

func TestBool_ConstructorNil(t *testing.T) {
	b := types.Bool(nil)

	// Test that it works with nil input
	// When nil, it should default to false value
	if b != nil && b.Value != false {
		t.Errorf("Bool(nil) = %v, want false", b.Value)
	}
}

func TestBoolValue_StringAfterSet(t *testing.T) {
	var b types.BoolValue

	// Test setting to true
	err := b.Set("true")
	if err != nil {
		t.Fatalf("Failed to set bool to true: %v", err)
	}

	if got := b.String(); got != "true" {
		t.Errorf("BoolValue.String() after Set(true) = %v, want %v", got, "true")
	}

	// Test setting to false
	err = b.Set("false")
	if err != nil {
		t.Fatalf("Failed to set bool to false: %v", err)
	}

	if got := b.String(); got != "false" {
		t.Errorf("BoolValue.String() after Set(false) = %v, want %v", got, "false")
	}
}

func TestBoolValue_MultipleSets(t *testing.T) {
	var b types.BoolValue

	// Test multiple sets
	err := b.Set("true")
	if err != nil {
		t.Fatalf("Failed to set bool to true: %v", err)
	}

	if !b.Value {
		t.Errorf("BoolValue after Set(true) = %v, want true", b.Value)
	}

	err = b.Set("false")
	if err != nil {
		t.Fatalf("Failed to set bool to false: %v", err)
	}

	if b.Value {
		t.Errorf("BoolValue after Set(false) = %v, want false", b.Value)
	}

	err = b.Set("1")
	if err != nil {
		t.Fatalf("Failed to set bool to 1: %v", err)
	}

	if !b.Value {
		t.Errorf("BoolValue after Set(1) = %v, want true", b.Value)
	}
}
