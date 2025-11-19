package types

import (
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-go-flags/utils"
)

var _ pflag.Value = (*BoolSliceValue)(nil)

type BoolSliceValue struct {
	value   *[]*wrapperspb.BoolValue
	changed bool
}

// Set converts, and assigns, the comma-separated boolean argument string representation as the []*wrapperspb.BoolValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.BoolValue assigned, the newly converted values will be appended.
func (s *BoolSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")
	// read flag arguments with CSV parser
	boolStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse boolean values into slice
	out := make([]*wrapperspb.BoolValue, 0, len(boolStrSlice))
	for _, boolStr := range boolStrSlice {
		b, err := strconv.ParseBool(strings.TrimSpace(boolStr))
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Bool(b))
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
func (s *BoolSliceValue) Type() string {
	return "boolSliceValue"
}

// String defines a "native" format for this boolean slice flag value.
func (s *BoolSliceValue) String() string {
	boolStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		boolStrSlice[i] = strconv.FormatBool(b.Value)
	}
	out, _ := utils.WriteAsCSV(boolStrSlice)

	return "[" + out + "]"
}

func BoolSlice(v *[]*wrapperspb.BoolValue) *BoolSliceValue {
	result := BoolSliceValue{value: v}
	return &result
}
