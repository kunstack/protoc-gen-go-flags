package types_test

import (
	"testing"
	"time"

	"github.com/kunstack/protoc-gen-go-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampSliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []*timestamppb.Timestamp
		expected string
	}{
		{
			name:     "empty slice",
			input:    []*timestamppb.Timestamp{},
			expected: "[]",
		},
		{
			name:     "single timestamp",
			input:    []*timestamppb.Timestamp{timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC))},
			expected: "[2023-12-25T10:30:00Z]",
		},
		{
			name: "multiple timestamps",
			input: []*timestamppb.Timestamp{
				timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
				timestamppb.New(time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC)),
			},
			expected: "[2023-12-25T10:30:00Z,2023-12-26T15:45:30Z]",
		},
		{
			name: "timestamps with nanoseconds",
			input: []*timestamppb.Timestamp{
				timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 123456789, time.UTC)),
				timestamppb.New(time.Date(2023, 12, 26, 15, 45, 30, 987654321, time.UTC)),
			},
			expected: "[2023-12-25T10:30:00.123456789Z,2023-12-26T15:45:30.987654321Z]",
		},
		{
			name: "nil timestamp in slice",
			input: []*timestamppb.Timestamp{
				timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
				nil,
				timestamppb.New(time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC)),
			},
			expected: "[2023-12-25T10:30:00Z,,2023-12-26T15:45:30Z]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tss := types.TimestampSlice(&tt.input, []string{time.RFC3339})
			if got := tss.String(); got != tt.expected {
				t.Errorf("TimestampSliceValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTimestampSliceValue_Set(t *testing.T) {
	tests := []struct {
		name         string
		formats      []string
		input        string
		wantTimes    []time.Time
		wantErr      bool
		errorMessage string
	}{
		{
			name:    "valid RFC3339",
			formats: []string{time.RFC3339},
			input:   "2023-12-25T10:30:00Z",
			wantTimes: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
		},
		{
			name:    "valid RFC822",
			formats: []string{time.RFC822},
			input:   "25 Dec 23 10:30 UTC",
			wantTimes: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
		},
		{
			name:    "valid RFC1123",
			formats: []string{time.RFC1123},
			input:   "Mon, 25 Dec 2023 10:30:00 UTC",
			wantTimes: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
		},
		{
			name:         "invalid format",
			formats:      []string{time.RFC3339},
			input:        "invalid-timestamp",
			wantErr:      true,
			errorMessage: "invalid time format `invalid-timestamp` must be one of: 2006-01-02T15:04:05Z07:00",
		},
		{
			name:         "empty string",
			formats:      []string{time.RFC3339},
			input:        "",
			wantErr:      true,
			errorMessage: "invalid time format `` must be one of: 2006-01-02T15:04:05Z07:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamps := make([]*timestamppb.Timestamp, 0)
			tss := types.TimestampSlice(&timestamps, tt.formats)
			err := tss.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimestampSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMessage {
				t.Errorf("TimestampSliceValue.Set() error message = %v, want %v", err.Error(), tt.errorMessage)
			}
			if !tt.wantErr {
				if len(timestamps) != len(tt.wantTimes) {
					t.Errorf("TimestampSliceValue.Set() length = %v, want %v", len(timestamps), len(tt.wantTimes))
					return
				}
				for i, want := range tt.wantTimes {
					got := timestamps[i].AsTime()
					if !got.Equal(want) {
						t.Errorf("TimestampSliceValue.Set() time[%d] = %v, want %v", i, got, want)
					}
				}
			}
		})
	}
}

func TestTimestampSliceValue_Type(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})
	if got := tss.Type(); got != "timestampSliceValue" {
		t.Errorf("TimestampSliceValue.Type() = %v, want %v", got, "timestampSliceValue")
	}
}

func TestTimestampSliceValue_PflagInterface(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that TimestampSliceValue implements pflag.Value interface
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})
	fs.Var(tss, "timestamps", "test timestamp slice")

	// Test setting flags
	args := []string{"--timestamps", "2023-12-25T10:30:00Z", "--timestamps", "2023-12-26T15:45:30Z"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Check the values through the TimestampSliceValue interface
	if len(timestamps) != 2 {
		t.Errorf("Expected 2 values, got %d", len(timestamps))
	}

	expected := []time.Time{
		time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
	}
	for i, exp := range expected {
		got := timestamps[i].AsTime()
		if !got.Equal(exp) {
			t.Errorf("Value[%d] = %v, want %v", i, got, exp)
		}
	}
}

func TestTimestampSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		initial []time.Time
		append  string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "append to empty slice",
			initial: []time.Time{},
			append:  "2023-12-25T10:30:00Z",
			want:    []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
		},
		{
			name:    "append to existing slice",
			initial: []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
			append:  "2023-12-26T15:45:30Z",
			want: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
			},
		},
		{
			name:    "append invalid timestamp",
			initial: []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
			append:  "invalid-timestamp",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			timestamps := make([]*timestamppb.Timestamp, len(tt.initial))
			for i, v := range tt.initial {
				timestamps[i] = timestamppb.New(v)
			}
			tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

			// Test Append
			err := tss.Append(tt.append)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Append() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify result
				if len(timestamps) != len(tt.want) {
					t.Errorf("Append() length = %v, want %v", len(timestamps), len(tt.want))
					return
				}

				for i, want := range tt.want {
					got := timestamps[i].AsTime()
					if !got.Equal(want) {
						t.Errorf("Append()[%d] = %v, want %v", i, got, want)
					}
				}
			}
		})
	}
}

func TestTimestampSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		initial []time.Time
		replace []string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "replace empty with values",
			initial: []time.Time{},
			replace: []string{"2023-12-25T10:30:00Z", "2023-12-26T15:45:30Z"},
			want: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
			},
		},
		{
			name:    "replace existing values",
			initial: []time.Time{time.Date(2023, 12, 24, 0, 0, 0, 0, time.UTC)},
			replace: []string{"2023-12-25T10:30:00Z", "2023-12-26T15:45:30Z", "2023-12-27T20:15:45Z"},
			want: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
				time.Date(2023, 12, 27, 20, 15, 45, 0, time.UTC),
			},
		},
		{
			name:    "replace with empty slice",
			initial: []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
			replace: []string{},
			want:    []time.Time{},
		},
		{
			name:    "replace with invalid timestamp",
			initial: []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
			replace: []string{"2023-12-25T10:30:00Z", "invalid-timestamp"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			timestamps := make([]*timestamppb.Timestamp, len(tt.initial))
			for i, v := range tt.initial {
				timestamps[i] = timestamppb.New(v)
			}
			tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

			// Test Replace
			err := tss.Replace(tt.replace)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Replace() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify result
				if len(timestamps) != len(tt.want) {
					t.Errorf("Replace() length = %v, want %v", len(timestamps), len(tt.want))
					return
				}

				for i, want := range tt.want {
					got := timestamps[i].AsTime()
					if !got.Equal(want) {
						t.Errorf("Replace()[%d] = %v, want %v", i, got, want)
					}
				}
			}
		})
	}
}

func TestTimestampSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func() []*timestamppb.Timestamp
		want  []string
	}{
		{
			name:  "get from empty slice",
			setup: func() []*timestamppb.Timestamp { return []*timestamppb.Timestamp{} },
			want:  []string{},
		},
		{
			name: "get from single value slice",
			setup: func() []*timestamppb.Timestamp {
				return []*timestamppb.Timestamp{timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC))}
			},
			want: []string{"2023-12-25T10:30:00Z"},
		},
		{
			name: "get from multiple values slice",
			setup: func() []*timestamppb.Timestamp {
				return []*timestamppb.Timestamp{
					timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
					timestamppb.New(time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC)),
					timestamppb.New(time.Date(2023, 12, 27, 20, 15, 45, 0, time.UTC)),
				}
			},
			want: []string{
				"2023-12-25T10:30:00Z",
				"2023-12-26T15:45:30Z",
				"2023-12-27T20:15:45Z",
			},
		},
		{
			name: "get with nil elements",
			setup: func() []*timestamppb.Timestamp {
				return []*timestamppb.Timestamp{
					timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
					nil,
					timestamppb.New(time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC)),
				}
			},
			want: []string{
				"2023-12-25T10:30:00Z",
				"", // empty string for nil
				"2023-12-26T15:45:30Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamps := tt.setup()
			tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

			// Test GetSlice
			got := tss.GetSlice()

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

func TestTimestampSliceValue_MultipleFormats(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339, time.RFC822, time.RFC1123})

	// Test different formats
	testCases := []struct {
		input    string
		expected time.Time
	}{
		{"2023-12-25T10:30:00Z", time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
		{"25 Dec 23 10:30 UTC", time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
		{"Mon, 25 Dec 2023 10:30:00 UTC", time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		err := tss.Set(tc.input)
		if err != nil {
			t.Errorf("Failed to set timestamp %q: %v", tc.input, err)
			continue
		}
	}

	if len(timestamps) != len(testCases) {
		t.Errorf("Expected %d timestamps, got %d", len(testCases), len(timestamps))
	}

	for i, tc := range testCases {
		got := timestamps[i].AsTime()
		if !got.Equal(tc.expected) {
			t.Errorf("Timestamp[%d] = %v, want %v", i, got, tc.expected)
		}
	}
}

func TestTimestampSliceValue_TimezoneHandling(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC1123Z})

	// Test timezone offset handling
	err := tss.Set("Mon, 25 Dec 2023 10:30:00 +0200")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := time.Date(2023, 12, 25, 8, 30, 0, 0, time.UTC) // 10:30 +02:00 = 8:30 UTC
	got := timestamps[0].AsTime()

	if !got.Equal(expected) {
		t.Errorf("Timezone conversion failed: got %v, want %v", got, expected)
	}
}

func TestTimestampSliceValue_Constructor(t *testing.T) {
	timestamps := []*timestamppb.Timestamp{
		timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
	}

	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

	// Test that the returned TimestampSliceValue has correct fields
	if len(timestamps) != 1 {
		t.Errorf("Expected 1 timestamp, got %d", len(timestamps))
	}

	// Test String() method on constructed timestamp slice
	expected := "[2023-12-25T10:30:00Z]"
	if got := tss.String(); got != expected {
		t.Errorf("TimestampSlice().String() = %v, want %v", got, expected)
	}
}

func TestTimestampSliceValue_ConstructorNil(t *testing.T) {
	tss := types.TimestampSlice(nil, []string{time.RFC3339})

	if tss == nil {
		t.Error("TimestampSlice(nil) returned nil")
	}

	if tss.Type() != "timestampSliceValue" {
		t.Errorf("TimestampSlice(nil).Type() = %v, want %v", tss.Type(), "timestampSliceValue")
	}
}

func TestTimestampSliceValue_SequentialOperations(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

	// Test sequential operations that simulate real command-line usage
	operations := []struct {
		name     string
		fn       func() error
		expected []time.Time
	}{
		{
			name:     "Set initial value",
			fn:       func() error { return tss.Set("2023-12-25T10:30:00Z") },
			expected: []time.Time{time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)},
		},
		{
			name: "Append value 1",
			fn:   func() error { return tss.Append("2023-12-26T15:45:30Z") },
			expected: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
			},
		},
		{
			name: "Append value 2",
			fn:   func() error { return tss.Append("2023-12-27T20:15:45Z") },
			expected: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
				time.Date(2023, 12, 27, 20, 15, 45, 0, time.UTC),
			},
		},
		{
			name: "Set new value (appends)",
			fn:   func() error { return tss.Set("2023-12-28T08:20:10Z") },
			expected: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
				time.Date(2023, 12, 27, 20, 15, 45, 0, time.UTC),
				time.Date(2023, 12, 28, 8, 20, 10, 0, time.UTC),
			},
		},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			if err := op.fn(); err != nil {
				t.Fatalf("%s failed: %v", op.name, err)
			}

			// Verify state after each operation
			if len(timestamps) != len(op.expected) {
				t.Errorf("After %s: expected %d elements, got %d", op.name, len(op.expected), len(timestamps))
			}
			for i, want := range op.expected {
				got := timestamps[i].AsTime()
				if !got.Equal(want) {
					t.Errorf("After %s: Element[%d] = %v, want %v", op.name, i, got, want)
				}
			}
		})
	}
}

func TestTimestampSliceValue_CommandLineScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario []string
		want     []time.Time
	}{
		{
			name:     "typical flag usage with multiple timestamps",
			scenario: []string{"2023-12-25T10:30:00Z", "2023-12-26T15:45:30Z", "2023-12-27T20:15:45Z"},
			want: []time.Time{
				time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 15, 45, 30, 0, time.UTC),
				time.Date(2023, 12, 27, 20, 15, 45, 0, time.UTC),
			},
		},
		{
			name:     "log timestamps",
			scenario: []string{"2023-12-25T08:00:00Z", "2023-12-25T12:30:15Z", "2023-12-25T18:45:30Z"},
			want: []time.Time{
				time.Date(2023, 12, 25, 8, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 25, 12, 30, 15, 0, time.UTC),
				time.Date(2023, 12, 25, 18, 45, 30, 0, time.UTC),
			},
		},
		{
			name:     "backup schedule timestamps",
			scenario: []string{"2023-12-25T02:00:00Z", "2023-12-26T02:00:00Z", "2023-12-27T02:00:00Z"},
			want: []time.Time{
				time.Date(2023, 12, 25, 2, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 26, 2, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 27, 2, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamps := make([]*timestamppb.Timestamp, 0)
			tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

			// Simulate command-line argument processing
			for _, arg := range tt.scenario {
				if err := tss.Set(arg); err != nil {
					t.Fatalf("Failed to set argument %q: %v", arg, err)
				}
			}

			// Verify results
			if len(timestamps) != len(tt.want) {
				t.Errorf("Expected %d elements, got %d", len(tt.want), len(timestamps))
				return
			}

			for i, want := range tt.want {
				got := timestamps[i].AsTime()
				if !got.Equal(want) {
					t.Errorf("Element[%d] = %v, want %v", i, got, want)
				}
			}

			// Test String() representation
			str := tss.String()
			if str == "" {
				t.Error("String() returned empty result")
			}

			// Test GetSlice() functionality
			slice := tss.GetSlice()
			if len(slice) != len(tt.want) {
				t.Errorf("GetSlice() returned %d elements, want %d", len(slice), len(tt.want))
			}
		})
	}
}

func TestTimestampSliceValue_MassiveData(t *testing.T) {
	// Test with massive amount of data
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

	// Add 1000 timestamps
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 1000; i++ {
		ts := baseTime.Add(time.Duration(i) * time.Hour)
		if err := tss.Set(ts.Format(time.RFC3339)); err != nil {
			t.Fatalf("Failed to set timestamp %d: %v", i, err)
		}
	}

	if len(timestamps) != 1000 {
		t.Errorf("Expected 1000 elements, got %d", len(timestamps))
	}

	// Test String() with massive data
	result := tss.String()
	if len(result) == 0 {
		t.Error("String() returned empty result for massive data")
	}

	// Test GetSlice() with massive data
	slice := tss.GetSlice()
	if len(slice) != 1000 {
		t.Errorf("GetSlice() returned %d elements, want 1000", len(slice))
	}
}

func TestTimestampSliceValue_ErrorRecovery(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

	// Set initial values
	if err := tss.Set("2023-12-25T10:30:00Z"); err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}
	if err := tss.Set("2023-12-26T15:45:30Z"); err != nil {
		t.Fatalf("Failed to set second value: %v", err)
	}

	// Test that operations maintain state consistency
	_ = len(timestamps) // original length for reference

	// Test Replace with empty slice
	if err := tss.Replace([]string{}); err != nil {
		t.Fatalf("Replace with empty slice failed: %v", err)
	}
	if len(timestamps) != 0 {
		t.Errorf("Expected empty slice after Replace, got %d elements", len(timestamps))
	}

	// Test that we can recover and add new values
	if err := tss.Set("2023-12-27T20:15:45Z"); err != nil {
		t.Fatalf("Failed to set after empty replace: %v", err)
	}
	if err := tss.Set("2023-12-28T08:20:10Z"); err != nil {
		t.Fatalf("Failed to set second value after empty replace: %v", err)
	}
	if len(timestamps) != 2 {
		t.Errorf("Expected 2 elements after recovery, got %d", len(timestamps))
	}
}

func TestTimestampSliceValue_MultipleOperations(t *testing.T) {
	timestamps := make([]*timestamppb.Timestamp, 0)
	tss := types.TimestampSlice(&timestamps, []string{time.RFC3339})

	// Test sequence: Set -> Append -> Replace -> GetSlice -> String

	// 1. Set initial values
	if err := tss.Set("2023-12-25T10:30:00Z"); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}
	if err := tss.Set("2023-12-26T15:45:30Z"); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	// 2. Append more values
	if err := tss.Append("2023-12-27T20:15:45Z"); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}
	if err := tss.Append("2023-12-28T08:20:10Z"); err != nil {
		t.Fatalf("Append() failed: %v", err)
	}

	// 3. Replace with new values
	if err := tss.Replace([]string{"2023-12-29T12:00:00Z", "2023-12-30T18:30:15Z"}); err != nil {
		t.Fatalf("Replace() failed: %v", err)
	}

	// 4. GetSlice and verify
	slice := tss.GetSlice()
	expectedStrings := []string{
		"2023-12-29T12:00:00Z",
		"2023-12-30T18:30:15Z",
	}
	if len(slice) != len(expectedStrings) {
		t.Errorf("GetSlice() length = %v, want %v", len(slice), len(expectedStrings))
	}
	for i, want := range expectedStrings {
		if slice[i] != want {
			t.Errorf("GetSlice()[%d] = %v, want %v", i, slice[i], want)
		}
	}

	// 5. String and verify
	str := tss.String()
	if str != "[2023-12-29T12:00:00Z,2023-12-30T18:30:15Z]" {
		t.Errorf("String() = %v, want %v", str, "[2023-12-29T12:00:00Z,2023-12-30T18:30:15Z]")
	}

	// 6. Type and verify
	typ := tss.Type()
	if typ != "timestampSliceValue" {
		t.Errorf("Type() = %v, want %v", typ, "timestampSliceValue")
	}
}

func TestTimestampSliceValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*types.TimestampSliceValue)(nil)
}
