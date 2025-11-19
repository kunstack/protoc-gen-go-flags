package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) checkMessage(typ pgs.FieldType, flag *flags.MessageFlag) {
	if !flag.GetNested() {
		return
	}
	m.mustType(typ, pgs.MessageT, pgs.UnknownWKT)
	if typ, ok := typ.(Repeatable); ok {
		m.Assert(!typ.IsRepeated(), "message flag does not support repeated fields")
	}
}

func (m *Module) genMessageDefaults(f pgs.Field, name pgs.Name, flag *flags.MessageFlag) string {
	var (
		declBuilder = &strings.Builder{}
	)
	if !flag.GetNested() {
		return fmt.Sprint("\n// ", name, ": flags disabled by [(flags.value).message = {nested: false}]")
	}
	if flag.GetNested() {
		_, _ = fmt.Fprintf(declBuilder, `
				if x.%s == nil {
					x.%s = new(%s)
				}
        	`,
			name, name, m.getFieldTypeName(f),
		)
	}
	_, _ = fmt.Fprintf(declBuilder, `
			if v, ok := interface{}(x.%s).(flags.Defaulter); ok {
				v.SetDefaults()
			}
        `,
		name,
	)
	return declBuilder.String()
}

func (m *Module) genMessage(f pgs.Field, name pgs.Name, flag *flags.MessageFlag) string {
	var (
		declBuilder = &strings.Builder{}
	)
	if !flag.GetNested() {
		return fmt.Sprint("\n// ", name, ": flags disabled by [(flags.value).message = {nested: false}]")
	}
	prefix := flag.GetName()
	if prefix == "" {
		// use field name instead
		prefix = strings.ToLower(f.Name().String())
	}
	if flag.GetNested() {
		_, _ = fmt.Fprintf(declBuilder, `
				if x.%s == nil {
					x.%s = new(%s)
				}
        	`,
			name, name, m.getFieldTypeName(f),
		)
	}
	_, _ = fmt.Fprintf(declBuilder, `
			if v, ok := interface{}(x.%s).(flags.Flagger); ok {
				v.AddFlags(fs, append(opts, flags.WithPrefix(%q))...)
			}
        `,
		name, prefix,
	)
	return declBuilder.String()
}
