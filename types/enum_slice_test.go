package types_test

import (
	"testing"

	testtypes "github.com/kunstack/protoc-gen-go-flags/tests"
	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// Test GetSlice function - the main functionality requested
func TestEnumSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func() []testtypes.TestEnum1
		want  []string
	}{
		{
			name: "empty slice returns empty string slice",
			setup: func() []testtypes.TestEnum1 {
				return []testtypes.TestEnum1{}
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enumSlice := tt.setup()
			es := types.EnumSlice(&enumSlice)

			got := es.GetSlice()
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test String method
func TestEnumSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		setup func() []testtypes.TestEnum1
		want  string
	}{
		{
			name: "empty slice",
			setup: func() []testtypes.TestEnum1 {
				return []testtypes.TestEnum1{}
			},
			want: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enumSlice := tt.setup()
			es := types.EnumSlice(&enumSlice)

			got := es.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test Type method
func TestEnumSliceValue_Type(t *testing.T) {
	var enumSlice []testtypes.TestEnum1
	es := types.EnumSlice(&enumSlice)

	got := es.Type()
	assert.Equal(t, "enumSliceValue", got)
}

// Test Set method with invalid enum values
func TestEnumSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "invalid enum value should error",
			input:   "INVALID_ENUM",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var enumSlice []testtypes.TestEnum1
			es := types.EnumSlice(&enumSlice)

			err := es.Set(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test Append method with invalid enum values
func TestEnumSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "append invalid enum should error",
			input:   "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var enumSlice []testtypes.TestEnum1
			es := types.EnumSlice(&enumSlice)

			err := es.Append(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test Replace method
func TestEnumSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr bool
	}{
		{
			name:    "replace with invalid enum values should error",
			input:   []string{"INVALID1", "INVALID2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var enumSlice []testtypes.TestEnum1
			es := types.EnumSlice(&enumSlice)

			err := es.Replace(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test EnumSlice function validation
func TestEnumSlice(t *testing.T) {
	t.Run("nil pointer should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			types.EnumSlice(nil)
		})
	})

	t.Run("non-pointer should panic", func(t *testing.T) {
		enumSlice := []interface{}{}
		assert.Panics(t, func() {
			types.EnumSlice(enumSlice) // Should panic because it's not a pointer
		})
	})

	t.Run("non-slice pointer should panic", func(t *testing.T) {
		var notASlice int
		assert.Panics(t, func() {
			types.EnumSlice(&notASlice)
		})
	})

	t.Run("valid slice pointer should work", func(t *testing.T) {
		var enumSlice []testtypes.TestEnum1
		assert.NotPanics(t, func() {
			es := types.EnumSlice(&enumSlice)
			assert.NotNil(t, es)
		})
	})
}

// Test pflag.Value interface implementation
func TestEnumSliceValue_PflagInterface(t *testing.T) {
	var enumSlice []testtypes.TestEnum1
	es := types.EnumSlice(&enumSlice)

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Should not panic when used as pflag.Value
	assert.NotPanics(t, func() {
		fs.Var(es, "enums", "test enum slice")
	})

	// Test that the flag can be parsed (enum validation will fail with invalid values)
	args := []string{"--enums", "INVALID_ENUM"}
	err := fs.Parse(args)
	assert.Error(t, err) // Should fail because INVALID_ENUM is not a valid enum value
}

// Test CSV parsing in Set method
func TestEnumSliceValue_CSVParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "single invalid value should error",
			input:   "INVALID_ENUM_VALUE",
			wantErr: true,
		},
		{
			name:    "comma separated invalid values should error",
			input:   "INVALID_ENUM1,INVALID_ENUM2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var enumSlice []testtypes.TestEnum1
			es := types.EnumSlice(&enumSlice)

			err := es.Set(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test error recovery and state consistency
func TestEnumSliceValue_ErrorRecovery(t *testing.T) {
	var enumSlice []testtypes.TestEnum1
	es := types.EnumSlice(&enumSlice)

	// Test sequence of operations with invalid values
	err1 := es.Set("INVALID_ENUM_1")
	assert.Error(t, err1)

	err2 := es.Append("INVALID_ENUM_2")
	assert.Error(t, err2)

	err3 := es.Replace([]string{"INVALID_ENUM_3", "INVALID_ENUM_4"})
	assert.Error(t, err3)

	// GetSlice should still work after errors
	result := es.GetSlice()
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

// Benchmark GetSlice performance
func BenchmarkEnumSliceValue_GetSlice(b *testing.B) {
	var enumSlice []testtypes.TestEnum1
	es := types.EnumSlice(&enumSlice)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = es.GetSlice()
	}
}
