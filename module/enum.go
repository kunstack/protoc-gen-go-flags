package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) checkEnum(ft pgs.FieldType, flag *flags.EnumFlag, pt pgs.ProtoType, wrapper pgs.WellKnownType) {
	m.checkCommon(ft, flag, pt, wrapper, false)

	typ, ok := ft.(interface {
		Enum() pgs.Enum
	})

	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}

	defined := typ.Enum().Values()

	if flag.Default != nil {
		for _, val := range defined {
			if val.Value() == *flag.Default {
				return
			}
		}
		m.Failf("unexpected enum value %d for %s", *flag.Default, typ.Enum().Name())
	}
}

func (m *Module) checkEnumSlice(ft FieldType, flag *flags.RepeatedEnumFlag, pt pgs.ProtoType, wrapper pgs.WellKnownType) {
	m.checkCommon(ft, flag, pt, wrapper, true)
	typ, ok := ft.(interface {
		Enum() pgs.Enum
	})

	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}

	if len(flag.Default) > 0 {
		defined := typ.Enum().Values()
		index := make(map[int32]struct{}, len(defined))
		for _, val := range defined {
			index[val.Value()] = struct{}{}
		}
		for _, val := range flag.Default {
			if _, ok := index[val]; !ok {
				m.Failf("unexpected enum value %d for %s", val, typ.Enum().Name())
			}
		}
	}
}

func (m *Module) genEnum(f pgs.Field, name pgs.Name, flag *flags.EnumFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	if f.HasOptionalKeyword() {
		_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
			name, name, m.getFieldTypeName(f),
		)
		_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.Enum(x.%s), builder.Build(%q), %q, %q)
		`,
			name, flag.Name, flag.GetShort(), flag.GetUsage(),
		)
	} else {
		_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.Enum(&x.%s), builder.Build(%q), %q, %q)
		`,
			name, flag.Name, flag.GetShort(), flag.GetUsage(),
		)
	}

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genEnumSlice(f pgs.Field, name pgs.Name, flag *flags.RepeatedEnumFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.EnumSlice(&x.%s), builder.Build(%q), %q, %q)
		`,
		name, flag.Name, flag.GetShort(), flag.GetUsage(),
	)

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
