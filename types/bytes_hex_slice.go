package types

import (
	"encoding/hex"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-go-flags/utils"
)

var (
	_ pflag.Value = (*BytesHexSliceValue)(nil)
	_ pflag.Value = (*BytesHexNativeSliceValue)(nil)
)

type BytesHexSliceValue struct {
	value   *[]*wrapperspb.BytesValue
	changed bool
}

func (s *BytesHexSliceValue) Set(val string) error {
	decodedBytes, err := hex.DecodeString(strings.TrimSpace(val))
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

func (s *BytesHexSliceValue) Append(val string) error {
	decodedBytes, err := hex.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	*s.value = append(*s.value, wrapperspb.Bytes(decodedBytes))
	return nil
}

func (s *BytesHexSliceValue) Replace(val []string) error {
	out := make([]*wrapperspb.BytesValue, len(val))
	for i, d := range val {
		decodedBytes, err := hex.DecodeString(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = wrapperspb.Bytes(decodedBytes)
	}
	*s.value = out
	return nil
}

func (s *BytesHexSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strings.ToUpper(hex.EncodeToString(d.Value))
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *BytesHexSliceValue) Type() string {
	return "bytesHexSliceValue"
}

// String defines a "native" format for this hexadecimal slice flag value.
func (s *BytesHexSliceValue) String() string {
	hexStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		hexStrSlice[i] = strings.ToUpper(hex.EncodeToString(b.Value))
	}
	out, _ := utils.WriteAsCSV(hexStrSlice)

	return "[" + out + "]"
}

type BytesHexNativeSliceValue struct {
	value   *[][]byte
	changed bool
}

func (s *BytesHexNativeSliceValue) Set(val string) error {
	decodedBytes, err := hex.DecodeString(strings.TrimSpace(val))
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

func (s *BytesHexNativeSliceValue) Append(val string) error {
	decodedBytes, err := hex.DecodeString(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	*s.value = append(*s.value, decodedBytes)
	return nil
}

func (s *BytesHexNativeSliceValue) Replace(val []string) error {
	out := make([][]byte, len(val))
	for i, d := range val {
		decodedBytes, err := hex.DecodeString(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = decodedBytes
	}
	*s.value = out
	return nil
}

func (s *BytesHexNativeSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strings.ToUpper(hex.EncodeToString(d))
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *BytesHexNativeSliceValue) Type() string {
	return "BytesHexNativeSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *BytesHexNativeSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = hex.EncodeToString(b)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)

	return "[" + out + "]"
}

func BytesHexSlice[T []*wrapperspb.BytesValue | [][]byte](v *T) pflag.Value {
	switch wrap := any(v).(type) {
	case *[]*wrapperspb.BytesValue:
		return &BytesHexSliceValue{value: wrap}
	case *[][]byte:
		return &BytesHexNativeSliceValue{value: wrap}
	default:
		// This should never happen due to type constraints
		panic("BytesSlice: unsupported type")
	}
}
