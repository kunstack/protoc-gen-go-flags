package module

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	"github.com/kunstack/protoc-gen-go-flags/utils"
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func isNowStr(s *string) bool {
	if s == nil {
		return false
	}
	nowStr := strings.ToLower(*s)
	return nowStr == "now" || nowStr == "now()"
}

func (m *Module) parseTimestamp(t string, layouts []string) (*timestamppb.Timestamp, error) {
	for _, format := range layouts {
		value, err := time.Parse(utils.ParseTimeFormat(format), t)
		if err == nil {
			return &timestamppb.Timestamp{
				Seconds: value.Unix(),
				Nanos:   int32(value.Nanosecond()),
			}, nil
		}
	}
	return nil, fmt.Errorf("timestamp '%s' could not be parsed with any of the provided formats: %v", t, layouts)
}

func (m *Module) checkTimestamp(ft pgs.FieldType, r *flags.TimestampFlag) {
	m.checkCommon(ft, r, pgs.MessageT, pgs.TimestampWKT, false)
	if len(r.Formats) == 0 {
		m.Failf("at least one format must be specified for timestamp flag")
	}
	// Check for duplicate formats
	formats := make(map[string]struct{})
	// Validate that formats are not empty strings
	for i, format := range r.Formats {
		if format == "" {
			m.Failf("timestamp format at index %d is empty", i)
		}
		if _, ok := formats[format]; ok {
			m.Failf("timestamp format '%s' at index %d is duplicated", format, i)
		}
		formats[format] = struct{}{}
	}
	if r.Default != nil && *r.Default != "" && !isNowStr(r.Default) {
		_, err := m.parseTimestamp(*r.Default, r.Formats)
		if err != nil {
			m.Failf("timestamp default value '%s' is invalid: %v", *r.Default, err)
		}
	}
}

func (m *Module) checkTimestampSlice(ft FieldType, r *flags.RepeatedTimestampFlag) {
	m.checkCommon(ft, r, pgs.MessageT, pgs.TimestampWKT, true)
	if len(r.Formats) == 0 {
		m.Failf("at least one format must be specified for timestamp flag")
	}
	// Check for duplicate formats
	formats := make(map[string]struct{})
	// Validate that formats are not empty strings
	for i, format := range r.Formats {
		if format == "" {
			m.Failf("timestamp format at index %d is empty", i)
		}
		if _, ok := formats[format]; ok {
			m.Failf("timestamp format '%s' at index %d is duplicated", format, i)
		}
		formats[format] = struct{}{}
	}

	for i, item := range r.Default {
		if item == "" {
			m.Failf("timestamp default value at index %d is empty", i)
		}
		if isNowStr(&item) {
			continue
		}
		_, err := m.parseTimestamp(item, r.Formats)
		if err != nil {
			m.Failf("timestamp default value '%s' at index %d  is invalid: %v", item, i, err)
		}
	}
}

func (m *Module) genTimestampDefaults(f pgs.Field, name pgs.Name, flag *flags.TimestampFlag) string {
	var (
		declBuilder = &strings.Builder{}
		timeBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	if flag.Default == nil || *flag.Default == "" {
		return ""
	}
	if isNowStr(flag.Default) {
		_, _ = fmt.Fprintf(timeBuilder, "timestamppb.Now()")
	} else {
		t, err := m.parseTimestamp(*flag.Default, flag.GetFormats())
		if err != nil {
			m.Failf("timestamp default value '%s' is invalid: %v", *flag.Default, err)
		}
		_, _ = fmt.Fprintf(timeBuilder, "&timestamppb.Timestamp{Seconds: %d, Nanos: %d}", t.Seconds, t.Nanos)
	}

	_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = %s
			}
		`,
		name, name, timeBuilder.String(),
	)

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

// genTimestampSliceDefaults generates default value assignment code for repeated timestamp fields
func (m *Module) genTimestampSliceDefaults(f pgs.Field, name pgs.Name, flag *flags.RepeatedTimestampFlag) string {
	if flag.Default == nil || len(flag.GetDefault()) == 0 {
		return ""
	}

	var code strings.Builder

	defaultValues := make([]string, len(flag.GetDefault()))

	for i, defaultValue := range flag.Default {
		if isNowStr(&defaultValue) {
			defaultValues[i] = "timestamppb.Now()"
		} else {
			t, err := m.parseTimestamp(defaultValue, flag.GetFormats())
			if err != nil {
				m.Failf("timestamp default value '%s' is invalid: %v", defaultValue, err)
			}
			defaultValues[i] = fmt.Sprintf("&timestamppb.Timestamp{Seconds: %d, Nanos: %d}", t.Seconds, t.Nanos)
		}
	}

	_, _ = fmt.Fprintf(&code, `
		if len(x.%s) == 0 {
		x.%s = %s{%s}
	}
`, name, name, m.getFieldTypeName(f), strings.Join(defaultValues, ","))
	return code.String()
}

func (m *Module) genTimestamp(f pgs.Field, name pgs.Name, flag *flags.TimestampFlag) string {
	var (
		declBuilder    = &strings.Builder{}
		formatsBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}
	_, _ = fmt.Fprint(formatsBuilder,
		`[]string{`,
		strings.Join(
			lo.Map(
				flag.GetFormats(),
				func(item string, _ int) string {
					return strconv.Quote(item)
				},
			),
			",",
		),
		`}`,
	)

	_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
		name, name, m.getFieldTypeName(f),
	)

	_, _ = fmt.Fprintf(declBuilder, `
		fs.VarP(types.Timestamp(x.%s, %s), builder.Build(%q), %q, %q)
	`,
		name, formatsBuilder.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(),
	)

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genTimestampSlice(f pgs.Field, name pgs.Name, flag *flags.RepeatedTimestampFlag) string {
	var (
		declBuilder    = &strings.Builder{}
		formatsBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}
	_, _ = fmt.Fprint(formatsBuilder,
		`[]string{`,
		strings.Join(
			lo.Map(
				flag.GetFormats(),
				func(item string, _ int) string {
					return strconv.Quote(item)
				},
			),
			",",
		),
		`}`,
	)

	_, _ = fmt.Fprintf(declBuilder, `
		fs.VarP(types.TimestampSlice(&x.%s, %s), builder.Build(%q), %q, %q)
	`,
		name, formatsBuilder.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(),
	)

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
