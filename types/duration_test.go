package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/kunstack/protoc-gen-go-flags/types"
)

func TestDurationValue_String(t *testing.T) {
	tests := []struct {
		name     string
		duration *durationpb.Duration
		expected string
	}{
		{
			name:     "nil duration",
			duration: nil,
			expected: "0s",
		},
		{
			name:     "zero duration",
			duration: &durationpb.Duration{Seconds: 0, Nanos: 0},
			expected: "0s",
		},
		{
			name:     "seconds only",
			duration: &durationpb.Duration{Seconds: 5, Nanos: 0},
			expected: "5s",
		},
		{
			name:     "nanoseconds only",
			duration: &durationpb.Duration{Seconds: 0, Nanos: 500000000},
			expected: "500ms",
		},
		{
			name:     "seconds and nanoseconds",
			duration: &durationpb.Duration{Seconds: 2, Nanos: 500000000},
			expected: "2.5s",
		},
		{
			name:     "negative duration",
			duration: &durationpb.Duration{Seconds: -1, Nanos: 0},
			expected: "-1s",
		},
		{
			name:     "large duration",
			duration: &durationpb.Duration{Seconds: 3600, Nanos: 0},
			expected: "1h0m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dv *types.DurationValue
			if tt.duration != nil {
				dv = types.Duration(tt.duration)
			}
			assert.Equal(t, tt.expected, dv.String())
		})
	}
}

func TestDurationValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *durationpb.Duration
		expectError bool
	}{
		{
			name:     "valid seconds",
			input:    "5s",
			expected: &durationpb.Duration{Seconds: 5, Nanos: 0},
		},
		{
			name:     "valid milliseconds",
			input:    "500ms",
			expected: &durationpb.Duration{Seconds: 0, Nanos: 500000000},
		},
		{
			name:     "valid minutes",
			input:    "2m",
			expected: &durationpb.Duration{Seconds: 120, Nanos: 0},
		},
		{
			name:     "valid hours",
			input:    "1h",
			expected: &durationpb.Duration{Seconds: 3600, Nanos: 0},
		},
		{
			name:     "valid complex duration",
			input:    "1h30m45s",
			expected: &durationpb.Duration{Seconds: 5445, Nanos: 0},
		},
		{
			name:     "valid fractional seconds",
			input:    "2.5s",
			expected: &durationpb.Duration{Seconds: 2, Nanos: 500000000},
		},
		{
			name:     "valid microseconds",
			input:    "100us",
			expected: &durationpb.Duration{Seconds: 0, Nanos: 100000},
		},
		{
			name:     "valid nanoseconds",
			input:    "50ns",
			expected: &durationpb.Duration{Seconds: 0, Nanos: 50},
		},
		{
			name:        "invalid format",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:     "zero duration",
			input:    "0s",
			expected: &durationpb.Duration{Seconds: 0, Nanos: 0},
		},
		{
			name:     "negative duration",
			input:    "-30s",
			expected: &durationpb.Duration{Seconds: -30, Nanos: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var duration durationpb.Duration
			dv := types.Duration(&duration)

			err := dv.Set(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Seconds, duration.Seconds)
				assert.Equal(t, tt.expected.Nanos, duration.Nanos)
			}
		})
	}
}

func TestDurationValue_Set_NilReceiver(t *testing.T) {
	var dv *types.DurationValue
	err := dv.Set("5s")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot set nil Duration")
}

func TestDurationValue_Type(t *testing.T) {
	duration := &durationpb.Duration{}
	dv := types.Duration(duration)
	assert.Equal(t, "durationValue", dv.Type())
}

func TestDurationValue_RoundTrip(t *testing.T) {
	// Test that we can set a duration and get the same string representation
	testDurations := []string{
		"1s",
		"500ms",
		"2m30s",
		"1h15m",
		"100us",
		"50ns",
		"0s",
		"-30s",
	}

	for _, durationStr := range testDurations {
		t.Run(durationStr, func(t *testing.T) {
			var duration durationpb.Duration
			dv := types.Duration(&duration)

			// Set the duration
			err := dv.Set(durationStr)
			require.NoError(t, err)

			// Parse the original duration for comparison
			expectedDuration, err := time.ParseDuration(durationStr)
			require.NoError(t, err)

			// Check that the protobuf duration matches
			protobufDuration := duration.AsDuration()
			assert.Equal(t, expectedDuration, protobufDuration)

			// Check that the string representation is reasonable
			// (Note: String() uses Go's duration formatting which may differ slightly)
			dvStr := dv.String()
			assert.NotEmpty(t, dvStr)
		})
	}
}

func TestDurationValue_EdgeCases(t *testing.T) {
	t.Run("very large nanoseconds", func(t *testing.T) {
		var duration durationpb.Duration
		dv := types.Duration(&duration)

		// Set a duration with nanoseconds > 1e9
		err := dv.Set("1.999999999s")
		require.NoError(t, err)
		assert.Equal(t, int64(1), duration.Seconds)
		assert.Equal(t, int32(999999999), duration.Nanos)
	})

	t.Run("fractional conversion", func(t *testing.T) {
		var duration durationpb.Duration
		dv := types.Duration(&duration)

		// Test that fractional seconds are properly converted
		err := dv.Set("0.123s")
		require.NoError(t, err)
		assert.Equal(t, int64(0), duration.Seconds)
		assert.Equal(t, int32(123000000), duration.Nanos)
	})
}
