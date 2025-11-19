package types_test

import (
	"testing"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBoolSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*wrapperspb.BoolValue
		wantErr bool
	}{
		{
			name:  "single true",
			input: "true",
			want:  []*wrapperspb.BoolValue{wrapperspb.Bool(true)},
		},
		{
			name:  "single false",
			input: "false",
			want:  []*wrapperspb.BoolValue{wrapperspb.Bool(false)},
		},
		{
			name:  "multiple values",
			input: "true,false,true",
			want:  []*wrapperspb.BoolValue{wrapperspb.Bool(true), wrapperspb.Bool(false), wrapperspb.Bool(true)},
		},
		{
			name:  "with spaces",
			input: " true , false ",
			want:  []*wrapperspb.BoolValue{wrapperspb.Bool(true), wrapperspb.Bool(false)},
		},
		{
			name:    "invalid value",
			input:   "true,invalid,false",
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "",
			want:  []*wrapperspb.BoolValue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.BoolValue, 0)
			bs := types.BoolSlice(&flags)
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BoolSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(flags) != len(tt.want) {
					t.Errorf("BoolSliceValue.Set() length = %v, want %v", len(flags), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if (flags)[i].Value != wantVal.Value {
						t.Errorf("BoolSliceValue.Set()[%d] = %v, want %v", i, (flags)[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBoolSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.BoolValue
		want  string
	}{
		{
			name:  "empty slice",
			input: []*wrapperspb.BoolValue{},
			want:  "[]",
		},
		{
			name:  "single true",
			input: []*wrapperspb.BoolValue{wrapperspb.Bool(true)},
			want:  "[true]",
		},
		{
			name:  "multiple values",
			input: []*wrapperspb.BoolValue{wrapperspb.Bool(true), wrapperspb.Bool(false), wrapperspb.Bool(true)},
			want:  "[true,false,true]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := types.BoolSlice(&tt.input)
			if got := bs.String(); got != tt.want {
				t.Errorf("BoolSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolSliceValue_Type(t *testing.T) {
	flags := []*wrapperspb.BoolValue{}
	bs := types.BoolSlice(&flags)
	if got := bs.Type(); got != "boolSliceValue" {
		t.Errorf("BoolSliceValue.Type() = %v, want %v", got, "boolSliceValue")
	}
}

func TestBoolSliceValue_PflagInterface(t *testing.T) {
	flags := make([]*wrapperspb.BoolValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BoolSliceValue implements pflag.Value interface
	bs := types.BoolSlice(&flags)
	fs.Var(bs, "bools", "test boolean slice")

	// Test setting flags
	args := []string{"--bools", "true,false,true"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Check the values through the BoolSliceValue interface
	if len(flags) != 3 {
		t.Errorf("Expected 3 values, got %d", len(flags))
	}

	if !(flags)[0].Value || (flags)[1].Value || !(flags)[2].Value {
		t.Errorf("Values not as expected: %v, %v, %v", (flags)[0].Value, (flags)[1].Value, (flags)[2].Value)
	}
}

func TestBoolSliceFunction(t *testing.T) {
	flags := []*wrapperspb.BoolValue{wrapperspb.Bool(true), wrapperspb.Bool(false)}
	bs := types.BoolSlice(&flags)

	if bs == nil {
		t.Error("BoolSlice() returned nil")
	}

	if len(flags) != 2 {
		t.Errorf("BoolSlice() length = %v, want %v", len(flags), 2)
	}

	if !(flags)[0].Value || (flags)[1].Value {
		t.Errorf("Values not as expected: %v, %v", (flags)[0].Value, (flags)[1].Value)
	}
}
