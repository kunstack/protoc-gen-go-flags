package types

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-go-flags/utils"
)

var _ pflag.Value = (*StringSliceValue)(nil)

type StringSliceValue struct {
	value   *[]*wrapperspb.StringValue
	changed bool
}

func (s *StringSliceValue) Set(val string) error {
	if !s.changed {
		*s.value = []*wrapperspb.StringValue{{Value: val}}
		s.changed = true
	} else {
		*s.value = append(*s.value, wrapperspb.String(val))
	}
	return nil
}

func (s *StringSliceValue) Append(val string) error {
	*s.value = append(*s.value, wrapperspb.String(val))
	return nil
}

func (s *StringSliceValue) Replace(val []string) error {
	out := make([]*wrapperspb.StringValue, len(val))
	for i, d := range val {
		out[i] = wrapperspb.String(d)
	}
	*s.value = out
	return nil
}

func (s *StringSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = d.GetValue()
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *StringSliceValue) Type() string {
	return "stringSliceValue"
}

// String defines a "native" format for this string slice flag value.
func (s *StringSliceValue) String() string {
	stringStrSlice := make([]string, len(*s.value))
	for i, str := range *s.value {
		stringStrSlice[i] = str.Value
	}
	out, _ := utils.WriteAsCSV(stringStrSlice)
	return "[" + out + "]"
}

func StringSlice(v *[]*wrapperspb.StringValue) *StringSliceValue {
	return &StringSliceValue{value: v}
}
