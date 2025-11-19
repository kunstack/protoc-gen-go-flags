package types

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-go-flags/utils"
)

var _ pflag.Value = (*BytesSliceValue)(nil)

type BytesSliceValue struct {
	value   *[]*wrapperspb.BytesValue
	changed bool
}

func (s *BytesSliceValue) Set(val string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = []*wrapperspb.BytesValue{{Value: decodedBytes}}
		s.changed = true
	} else {
		*s.value = append(*s.value, wrapperspb.Bytes(decodedBytes))
	}
	return nil
}

func (s *BytesSliceValue) Append(val string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	*s.value = append(*s.value, wrapperspb.Bytes(decodedBytes))
	return nil
}

func (s *BytesSliceValue) Replace(val []string) error {
	out := make([]*wrapperspb.BytesValue, len(val))
	for i, d := range val {
		decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = wrapperspb.Bytes(decodedBytes)
	}
	*s.value = out
	return nil
}

func (s *BytesSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strings.ToUpper(base64.StdEncoding.EncodeToString(d.Value))
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *BytesSliceValue) Type() string {
	return "bytesSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *BytesSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = base64.StdEncoding.EncodeToString(b.Value)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)
	return "[" + out + "]"
}

// NativeBytesSliceValue is a pflag.Value implementation for handling native [][]byte fields.
// It supports base64 encoding for command-line flag parsing.
type NativeBytesSliceValue struct {
	value   *[][]byte
	changed bool
}

func (s *NativeBytesSliceValue) Set(val string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = [][]byte{decodedBytes}
		s.changed = true
	} else {
		*s.value = append(*s.value, decodedBytes)
	}
	return nil
}

func (s *NativeBytesSliceValue) Append(val string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	*s.value = append(*s.value, decodedBytes)
	return nil
}

func (s *NativeBytesSliceValue) Replace(val []string) error {
	out := make([][]byte, len(val))
	for i, d := range val {
		decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = decodedBytes
	}
	*s.value = out
	return nil
}

func (s *NativeBytesSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strings.ToUpper(base64.StdEncoding.EncodeToString(d))
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *NativeBytesSliceValue) Type() string {
	return "nativeBytesSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *NativeBytesSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = base64.StdEncoding.EncodeToString(b)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)

	return "[" + out + "]"
}

func BytesSlice[T []*wrapperspb.BytesValue | [][]byte](v *T) pflag.Value {
	switch wrap := any(v).(type) {
	case *[]*wrapperspb.BytesValue:
		return &BytesSliceValue{value: wrap}
	case *[][]byte:
		return &NativeBytesSliceValue{value: wrap}
	default:
		// This should never happen due to type constraints
		panic("BytesSlice: unsupported type")
	}
}
