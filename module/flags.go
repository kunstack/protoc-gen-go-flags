package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) genMark(flag commonFlag) string {
	var declBuilder = &strings.Builder{}
	if flag.GetHidden() {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.MarkHidden(%q)
			`,
			flag.GetName())
	}
	if flag.GetDeprecated() {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.MarkDeprecated(%q, %q)
			`,
			flag.GetName(), flag.GetDeprecatedUsage())
	}
	return declBuilder.String()
}

func (m *Module) processRepeatedFlag(f pgs.Field, name pgs.Name, repeated *flags.RepeatedFlags) string {
	if repeated == nil {
		return ""
	}
	wk := pgs.UnknownWKT
	if emb := f.Type().Element().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}

	switch r := repeated.Type.(type) {
	case *flags.RepeatedFlags_Float:
		return m.genCommonSlice(f, name, r.Float, wk, "FloatSlice", "Float32SliceVarP")
	case *flags.RepeatedFlags_Double:
		return m.genCommonSlice(f, name, r.Double, wk, "DoubleSlice", "Float64SliceVarP")
	case *flags.RepeatedFlags_Int32:
		return m.genCommonSlice(f, name, r.Int32, wk, "Int32Slice", "Int32SliceVarP")
	case *flags.RepeatedFlags_Int64:
		return m.genCommonSlice(f, name, r.Int64, wk, "Int64Slice", "Int64SliceVarP")
	case *flags.RepeatedFlags_Uint32:
		return m.genCommonSlice(f, name, r.Uint32, wk, "UInt32Slice", "Uint32SliceVarP")
	case *flags.RepeatedFlags_Uint64:
		return m.genCommonSlice(f, name, r.Uint64, wk, "UInt64Slice", "Uint64SliceVarP")
	case *flags.RepeatedFlags_Sint32:
		return m.genCommonSlice(f, name, r.Sint32, wk, "Int32Slice", "Int32SliceVarP")
	case *flags.RepeatedFlags_Sint64:
		return m.genCommonSlice(f, name, r.Sint64, wk, "Int64Slice", "Int64SliceVarP")
	case *flags.RepeatedFlags_Fixed32:
		return m.genCommonSlice(f, name, r.Fixed32, wk, "UInt32Slice", "Uint32SliceVarP")
	case *flags.RepeatedFlags_Fixed64:
		return m.genCommonSlice(f, name, r.Fixed64, wk, "UInt32Slice", "Uint64SliceVarP")
	case *flags.RepeatedFlags_Sfixed32:
		return m.genCommonSlice(f, name, r.Sfixed32, wk, "Int32Slice", "Int32SliceVarP")
	case *flags.RepeatedFlags_Sfixed64:
		return m.genCommonSlice(f, name, r.Sfixed64, wk, "Int64Slice", "Int64SliceVarP")
	case *flags.RepeatedFlags_Bool:
		return m.genCommonSlice(f, name, r.Bool, wk, "BoolSlice", "BoolSliceVarP")
	case *flags.RepeatedFlags_String_:
		return m.genCommonSlice(f, name, r.String_, wk, "StringSlice", "StringSliceVarP")
	case *flags.RepeatedFlags_Bytes:
		return m.genBytesSlice(name, r.Bytes)
	case *flags.RepeatedFlags_Enum:
		return m.genEnumSlice(f, name, r.Enum, wk)
	case *flags.RepeatedFlags_Duration:
		return m.genDurationSlice(f, name, r.Duration, wk)
	case *flags.RepeatedFlags_Timestamp:
		return m.genTimestampSlice(f, name, r.Timestamp)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", f.Type)
	}
	return ""
}

func (m *Module) genFieldFlags(f pgs.Field) string {
	m.Push(f.Name().String())
	defer m.Pop()
	var field flags.FieldFlags
	ok, err := f.Extension(flags.E_Value, &field)
	if err != nil || !ok {
		return ""
	}

	wk := pgs.UnknownWKT
	if emb := f.Type().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}

	// 检查是否是oneof字段，发出警告并跳过处理
	if f.InRealOneOf() {
		m.Failf("oneof field '%s' is not supported for flag generation. Please use regular fields instead.", f.Name())
		return ""
	}

	name := m.ctx.Name(f)
	switch r := field.Type.(type) {
	case *flags.FieldFlags_Float:
		return m.genCommon(f, name, r.Float, wk, "Float", "Float32VarP")
	case *flags.FieldFlags_Double:
		return m.genCommon(f, name, r.Double, wk, "Double", "Float64VarP")
	case *flags.FieldFlags_Int32:
		return m.genCommon(f, name, r.Int32, wk, "Int32", "Int32VarP")
	case *flags.FieldFlags_Int64:
		return m.genCommon(f, name, r.Int64, wk, "Int64", "Int64VarP")
	case *flags.FieldFlags_Uint32:
		return m.genCommon(f, name, r.Uint32, wk, "UInt32", "Uint32VarP")
	case *flags.FieldFlags_Uint64:
		return m.genCommon(f, name, r.Uint64, wk, "UInt64", "Uint64VarP")
	case *flags.FieldFlags_Sint32:
		return m.genCommon(f, name, r.Sint32, wk, "Int32", "Int32VarP")
	case *flags.FieldFlags_Sint64:
		return m.genCommon(f, name, r.Sint64, wk, "Int64", "Int64VarP")
	case *flags.FieldFlags_Fixed32:
		return m.genCommon(f, name, r.Fixed32, wk, "UInt32", "Uint32VarP")
	case *flags.FieldFlags_Fixed64:
		return m.genCommon(f, name, r.Fixed64, wk, "UInt64", "Uint64VarP")
	case *flags.FieldFlags_Sfixed32:
		return m.genCommon(f, name, r.Sfixed32, wk, "Int32", "Int32VarP")
	case *flags.FieldFlags_Sfixed64:
		return m.genCommon(f, name, r.Sfixed64, wk, "Int64", "Int64VarP")
	case *flags.FieldFlags_Bool:
		return m.genCommon(f, name, r.Bool, wk, "Bool", "BoolVarP")
	case *flags.FieldFlags_String_:
		return m.genCommon(f, name, r.String_, wk, "String", "StringVarP")
	case *flags.FieldFlags_Bytes:
		return m.genBytes(f, name, r.Bytes, wk)
	case *flags.FieldFlags_Enum:
		return m.genEnum(f, name, r.Enum, wk)
	case *flags.FieldFlags_Duration:
		return m.genDuration(f, name, r.Duration, wk)
	case *flags.FieldFlags_Timestamp:
		return m.genTimestamp(f, name, r.Timestamp)
	case *flags.FieldFlags_Message:
		return m.genMessage(f, name, r.Message)
	case *flags.FieldFlags_Map:
		return m.genMap(f, name, r.Map)
	case *flags.FieldFlags_Repeated:
		return m.processRepeatedFlag(f, name, r.Repeated)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)
	}
	return fmt.Sprint("\n// ", f.Name())
}
