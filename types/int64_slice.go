package types

import (
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-go-flags/utils"
)

var _ pflag.Value = (*Int64SliceValue)(nil)

type Int64SliceValue struct {
	value   *[]*wrapperspb.Int64Value
	changed bool
}

func (s *Int64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]*wrapperspb.Int64Value, len(ss))
	for i, d := range ss {
		var err error
		var temp int64
		temp, err = strconv.ParseInt(strings.TrimSpace(d), 10, 64)
		if err != nil {
			return err
		}
		out[i] = wrapperspb.Int64(temp)
	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *Int64SliceValue) Type() string {
	return "int64SliceValue"
}

// String defines a "native" format for this int64 slice flag value.
func (s *Int64SliceValue) String() string {
	int64StrSlice := make([]string, len(*s.value))
	for i, v := range *s.value {
		int64StrSlice[i] = strconv.FormatInt(v.Value, 10)
	}
	out, _ := utils.WriteAsCSV(int64StrSlice)

	return "[" + out + "]"
}

func Int64Slice(v *[]*wrapperspb.Int64Value) *Int64SliceValue {
	return &Int64SliceValue{value: v}
}
