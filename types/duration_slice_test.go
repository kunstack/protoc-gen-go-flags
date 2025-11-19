package types_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/kunstack/protoc-gen-go-flags/types"
)

func TestDurationSliceValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []*durationpb.Duration
		expectError bool
	}{
		{
			name:     "single duration",
			input:    "5s",
			expected: []*durationpb.Duration{durationpb.New(5 * 1000000000)},
		},
		{
			name:  "multiple durations",
			input: "1s,2m,3h",
			expected: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
				durationpb.New(10800 * 1000000000),
			},
		},
		{
			name:  "durations with spaces",
			input: " 1s , 2m , 3h ",
			expected: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
				durationpb.New(10800 * 1000000000),
			},
		},
		{
			name:  "mixed time units",
			input: "500ms,100us,50ns",
			expected: []*durationpb.Duration{
				durationpb.New(500000000),
				durationpb.New(100000),
				durationpb.New(50),
			},
		},
		{
			name:  "fractional durations",
			input: "2.5s,1.5m",
			expected: []*durationpb.Duration{
				durationpb.New(2500000000),
				durationpb.New(90000000000),
			},
		},
		{
			name:  "negative durations",
			input: "-30s,-1m",
			expected: []*durationpb.Duration{
				durationpb.New(-30000000000),
				durationpb.New(-60000000000),
			},
		},
		{
			name:     "zero duration",
			input:    "0s",
			expected: []*durationpb.Duration{durationpb.New(0)},
		},
		{
			name:        "invalid duration format",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "mixed valid and invalid",
			input:       "1s,invalid,3h",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var durations []*durationpb.Duration
			dsv := types.DurationSlice(&durations)

			err := dsv.Set(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, durations, len(tt.expected))
				for i, expected := range tt.expected {
					assert.Equal(t, expected.Seconds, durations[i].Seconds)
					assert.Equal(t, expected.Nanos, durations[i].Nanos)
				}
			}
		})
	}
}

func TestDurationSliceValue_Append(t *testing.T) {
	tests := []struct {
		name        string
		initial     []*durationpb.Duration
		appendValue string
		expected    []*durationpb.Duration
		expectError bool
	}{
		{
			name:        "append to empty slice",
			initial:     []*durationpb.Duration{},
			appendValue: "5s",
			expected: []*durationpb.Duration{
				durationpb.New(5 * 1000000000),
			},
		},
		{
			name: "append to existing slice",
			initial: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
			},
			appendValue: "2m",
			expected: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
			},
		},
		{
			name:        "append with spaces",
			initial:     []*durationpb.Duration{},
			appendValue: "  5s  ",
			expected: []*durationpb.Duration{
				durationpb.New(5 * 1000000000),
			},
		},
		{
			name:        "append invalid duration",
			initial:     []*durationpb.Duration{},
			appendValue: "invalid",
			expectError: true,
		},
		{
			name:        "append empty string",
			initial:     []*durationpb.Duration{},
			appendValue: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			durations := make([]*durationpb.Duration, len(tt.initial))
			copy(durations, tt.initial)
			dsv := types.DurationSlice(&durations)

			err := dsv.Append(tt.appendValue)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, durations, len(tt.expected))
				for i, expected := range tt.expected {
					assert.Equal(t, expected.Seconds, durations[i].Seconds)
					assert.Equal(t, expected.Nanos, durations[i].Nanos)
				}
			}
		})
	}
}

func TestDurationSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name        string
		initial     []*durationpb.Duration
		replaceWith []string
		expected    []*durationpb.Duration
		expectError bool
	}{
		{
			name: "replace existing values",
			initial: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(2 * 1000000000),
			},
			replaceWith: []string{"5s", "10s"},
			expected: []*durationpb.Duration{
				durationpb.New(5 * 1000000000),
				durationpb.New(10 * 1000000000),
			},
		},
		{
			name:        "replace with empty slice",
			initial:     []*durationpb.Duration{durationpb.New(1 * 1000000000)},
			replaceWith: []string{},
			expected:    []*durationpb.Duration{},
		},
		{
			name:        "replace empty with values",
			initial:     []*durationpb.Duration{},
			replaceWith: []string{"1s", "2s"},
			expected: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(2 * 1000000000),
			},
		},
		{
			name:        "replace with invalid duration",
			initial:     []*durationpb.Duration{},
			replaceWith: []string{"1s", "invalid"},
			expectError: true,
		},
		{
			name:        "replace with values containing spaces",
			initial:     []*durationpb.Duration{},
			replaceWith: []string{" 1s ", " 2m "},
			expected: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			durations := make([]*durationpb.Duration, len(tt.initial))
			copy(durations, tt.initial)
			dsv := types.DurationSlice(&durations)

			err := dsv.Replace(tt.replaceWith)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, durations, len(tt.expected))
				for i, expected := range tt.expected {
					assert.Equal(t, expected.Seconds, durations[i].Seconds)
					assert.Equal(t, expected.Nanos, durations[i].Nanos)
				}
			}
		})
	}
}

func TestDurationSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name     string
		initial  []*durationpb.Duration
		expected []string
	}{
		{
			name:     "empty slice",
			initial:  []*durationpb.Duration{},
			expected: []string{},
		},
		{
			name: "single duration",
			initial: []*durationpb.Duration{
				durationpb.New(5 * 1000000000),
			},
			expected: []string{"5s"},
		},
		{
			name: "multiple durations",
			initial: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
				durationpb.New(10800 * 1000000000),
			},
			expected: []string{"1s", "2m0s", "3h0m0s"},
		},
		{
			name: "fractional durations",
			initial: []*durationpb.Duration{
				durationpb.New(2500000000),
				durationpb.New(500000000),
			},
			expected: []string{"2.5s", "500ms"},
		},
		{
			name: "negative durations",
			initial: []*durationpb.Duration{
				durationpb.New(-30000000000),
				durationpb.New(-60000000000),
			},
			expected: []string{"-30s", "-1m0s"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			durations := make([]*durationpb.Duration, len(tt.initial))
			copy(durations, tt.initial)
			dsv := types.DurationSlice(&durations)

			result := dsv.GetSlice()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDurationSliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		initial  []*durationpb.Duration
		expected string
	}{
		{
			name:     "empty slice",
			initial:  []*durationpb.Duration{},
			expected: "[]",
		},
		{
			name: "single duration",
			initial: []*durationpb.Duration{
				durationpb.New(5 * 1000000000),
			},
			expected: "[5s]",
		},
		{
			name: "multiple durations",
			initial: []*durationpb.Duration{
				durationpb.New(1 * 1000000000),
				durationpb.New(120 * 1000000000),
				durationpb.New(10800 * 1000000000),
			},
			expected: "[1s,2m0s,3h0m0s]",
		},
		{
			name: "fractional durations",
			initial: []*durationpb.Duration{
				durationpb.New(2500000000),
				durationpb.New(500000000),
			},
			expected: "[2.5s,500ms]",
		},
		{
			name: "negative durations",
			initial: []*durationpb.Duration{
				durationpb.New(-30000000000),
				durationpb.New(-60000000000),
			},
			expected: "[-30s,-1m0s]",
		},
		{
			name: "zero duration",
			initial: []*durationpb.Duration{
				durationpb.New(0),
			},
			expected: "[0s]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			durations := make([]*durationpb.Duration, len(tt.initial))
			copy(durations, tt.initial)
			dsv := types.DurationSlice(&durations)

			result := dsv.String()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDurationSliceValue_Type(t *testing.T) {
	var durations []*durationpb.Duration
	dsv := types.DurationSlice(&durations)
	assert.Equal(t, "DurationSliceValue", dsv.Type())
}

func TestDurationSliceValue_MultipleSets(t *testing.T) {
	var durations []*durationpb.Duration
	dsv := types.DurationSlice(&durations)

	// First set
	err := dsv.Set("1s,2s")
	require.NoError(t, err)
	assert.Len(t, durations, 2)

	// Second set should append due to changed flag
	err = dsv.Set("3s,4s")
	require.NoError(t, err)
	assert.Len(t, durations, 4) // Should have 4 elements (appended)

	// Verify the values
	assert.Equal(t, int64(1), durations[0].Seconds)
	assert.Equal(t, int64(2), durations[1].Seconds)
	assert.Equal(t, int64(3), durations[2].Seconds)
	assert.Equal(t, int64(4), durations[3].Seconds)
}

func TestDurationSliceValue_EdgeCases(t *testing.T) {
	t.Run("very large nanoseconds", func(t *testing.T) {
		var durations []*durationpb.Duration
		dsv := types.DurationSlice(&durations)

		err := dsv.Set("1.999999999s,2.000000001s")
		require.NoError(t, err)
		assert.Len(t, durations, 2)
		assert.Equal(t, int64(1), durations[0].Seconds)
		assert.Equal(t, int32(999999999), durations[0].Nanos)
		assert.Equal(t, int64(2), durations[1].Seconds)
		assert.Equal(t, int32(1), durations[1].Nanos)
	})

	t.Run("complex duration parsing", func(t *testing.T) {
		var durations []*durationpb.Duration
		dsv := types.DurationSlice(&durations)

		err := dsv.Set("1h30m45s,2h15m30s")
		require.NoError(t, err)
		assert.Len(t, durations, 2)
		assert.Equal(t, int64(5445), durations[0].Seconds) // 1h30m45s = 5445s
		assert.Equal(t, int64(8130), durations[1].Seconds) // 2h15m30s = 8130s
	})

	t.Run("empty elements handling", func(t *testing.T) {
		var durations []*durationpb.Duration
		dsv := types.DurationSlice(&durations)

		// Test with empty elements in comma-separated list
		err := dsv.Set("1s,,3s")
		assert.Error(t, err) // Should fail on empty duration string
	})
}

func TestDurationSliceValue_Constructor(t *testing.T) {
	t.Run("constructor with nil slice", func(t *testing.T) {
		var durations []*durationpb.Duration
		dsv := types.DurationSlice(&durations)
		assert.NotNil(t, dsv)
		assert.Equal(t, "[]", dsv.String())
	})

	t.Run("constructor with existing slice", func(t *testing.T) {
		initial := []*durationpb.Duration{
			durationpb.New(5 * 1000000000),
		}
		dsv := types.DurationSlice(&initial)
		assert.NotNil(t, dsv)
		assert.Equal(t, "[5s]", dsv.String())
	})
}

func TestDurationSliceValue_ImplementsInterfaces(t *testing.T) {
	var durations []*durationpb.Duration
	dsv := types.DurationSlice(&durations)

	// Test that it implements pflag.Value
	var _ pflag.Value = dsv

	// Test that it implements pflag.SliceValue
	var _ pflag.SliceValue = dsv
}
