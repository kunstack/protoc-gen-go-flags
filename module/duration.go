package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) genDuration(f pgs.Field, name pgs.Name, flag *flags.DurationFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
		name, name, m.getFieldTypeName(f),
	)

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.Duration(x.%s), builder.Build(%q), %q, %q)
		`,
		name, flag.Name, flag.GetShort(), flag.GetUsage(),
	)

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genDurationSlice(f pgs.Field, name pgs.Name, flag *flags.RepeatedDurationFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.DurationSlice(&x.%s), builder.Build(%q), %q, %q)
		`,
		name, flag.Name, flag.GetShort(), flag.GetUsage(),
	)

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

// genDurationDefaults generates default value assignment code for duration fields.
// It handles duration parsing and creates appropriate default assignment code.
//
// Parameters:
//   - f: The protobuf field (should be a duration field)
//   - name: The Go field name
//   - flag: The duration flag configuration
//   - wk: Well-known type information
//
// Returns:
//   - Generated Go code for duration default assignment
//   - Empty string if no default should be generated
func (m *Module) genDurationDefaults(f pgs.Field, name pgs.Name, flag *flags.DurationFlag) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	// Check if default value is configured
	if flag.Default == nil || *flag.Default == "" {
		return ""
	}

	durationStr := *flag.Default

	// Parse the duration string
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		m.Failf("duration default value '%s' is invalid: %v", durationStr, err)
		return ""
	}

	nanos := duration.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	// For regular duration fields
	_, _ = fmt.Fprintf(declBuilder, `
		if x.%s == nil {
			x.%s = &durationpb.Duration{Seconds: %d, Nanos: %d}
		}`, name, name, int64(secs), int32(nanos))

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

// genDurationSliceDefaults generates default value assignment code for repeated duration fields
func (m *Module) genDurationSliceDefaults(f pgs.Field, name pgs.Name, flag *flags.RepeatedDurationFlag, wk pgs.WellKnownType) string {
	if flag.Default == nil || len(flag.GetDefault()) == 0 {
		return ""
	}

	var code strings.Builder

	// Check if the slice is empty before setting defaults
	code.WriteString(fmt.Sprintf(`
	if len(x.%s) == 0 {`, name))

	for i, defaultValue := range flag.Default {
		varName := fmt.Sprintf("value%d", i)
		code.WriteString(fmt.Sprintf(`
		%s, _ := time.ParseDuration(%q)
		x.%s = append(x.%s, %s)`, varName, defaultValue, name, name, varName))
	}

	code.WriteString(`
	}`)
	return code.String()
}
