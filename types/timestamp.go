package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/kunstack/protoc-gen-go-flags/utils"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ pflag.Value = (*TimestampValue)(nil)

type TimestampValue struct {
	layouts []string
	wrap    *timestamppb.Timestamp
}

func (t *TimestampValue) String() string {
	if t.wrap != nil {
		w := t.wrap.AsTime()
		if !w.IsZero() {
			return w.Format(time.RFC3339Nano)
		}
	}
	return ""
}

func (t *TimestampValue) Set(s string) error {
	for _, format := range t.layouts {
		actualFormat := utils.ParseTimeFormat(format)
		value, err := time.Parse(actualFormat, s)
		if err == nil {
			t.wrap.Seconds = value.Unix()
			t.wrap.Nanos = int32(value.Nanosecond())
			return nil
		}
	}
	return fmt.Errorf("invalid time format `%s` must be one of: %s", s, strings.Join(t.layouts, ","))
}

func (t *TimestampValue) Type() string {
	return "timestampValue"
}

func Timestamp(wrap *timestamppb.Timestamp, formats []string) *TimestampValue {
	return &TimestampValue{wrap: wrap, layouts: formats}
}
