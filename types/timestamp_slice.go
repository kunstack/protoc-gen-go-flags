package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/kunstack/protoc-gen-go-flags/utils"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ pflag.Value      = (*TimestampSliceValue)(nil)
	_ pflag.SliceValue = (*TimestampSliceValue)(nil)
)

type TimestampSliceValue struct {
	value   *[]*timestamppb.Timestamp
	layouts []string
	changed bool
}

func (t *TimestampSliceValue) Set(val string) error {
	for _, format := range t.layouts {
		actualFormat := utils.ParseTimeFormat(format)
		parsedTime, err := time.Parse(actualFormat, val)
		if err == nil {
			ts := timestamppb.New(parsedTime)
			if !t.changed {
				*t.value = []*timestamppb.Timestamp{ts}
				t.changed = true
			} else {
				*t.value = append(*t.value, ts)
			}
			return nil
		}
	}
	return fmt.Errorf("invalid time format `%s` must be one of: %s", val, strings.Join(t.layouts, ","))
}

func (t *TimestampSliceValue) Append(val string) error {
	for _, format := range t.layouts {
		actualFormat := utils.ParseTimeFormat(format)
		parsedTime, err := time.Parse(actualFormat, val)
		if err == nil {
			*t.value = append(*t.value, timestamppb.New(parsedTime))
			return nil
		}
	}
	return fmt.Errorf("invalid time format `%s` must be one of: %s", val, strings.Join(t.layouts, ","))
}

func (t *TimestampSliceValue) Replace(val []string) error {
	out := make([]*timestamppb.Timestamp, 0, len(val))
	for _, timestampStr := range val {
		var parsed bool
		for _, format := range t.layouts {
			actualFormat := utils.ParseTimeFormat(format)
			parsedTime, err := time.Parse(actualFormat, timestampStr)
			if err == nil {
				out = append(out, timestamppb.New(parsedTime))
				parsed = true
				break
			}
		}
		if !parsed {
			return fmt.Errorf("invalid time format `%s` must be one of: %s", timestampStr, strings.Join(t.layouts, ","))
		}
	}
	*t.value = out
	return nil
}

func (t *TimestampSliceValue) GetSlice() []string {
	out := make([]string, len(*t.value))
	for i, ts := range *t.value {
		if ts != nil {
			t := ts.AsTime()
			if !t.IsZero() {
				out[i] = ts.AsTime().Format(time.RFC3339Nano)
			}
		}
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (t *TimestampSliceValue) Type() string {
	return "timestampSliceValue"
}

// String defines a "native" format for this timestamp slice flag value.
func (t *TimestampSliceValue) String() string {
	timestampStrSlice := make([]string, len(*t.value))
	for i, ts := range *t.value {
		if ts != nil {
			t := ts.AsTime()
			if !t.IsZero() {
				timestampStrSlice[i] = ts.AsTime().Format(time.RFC3339Nano)
			}
		}
	}
	out, _ := utils.WriteAsCSV(timestampStrSlice)
	return "[" + out + "]"
}

func TimestampSlice(v *[]*timestamppb.Timestamp, formats []string) *TimestampSliceValue {
	return &TimestampSliceValue{value: v, layouts: formats}
}
