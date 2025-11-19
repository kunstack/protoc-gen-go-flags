package module

import (
	"fmt"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) genFieldDefaults(f pgs.Field) string {
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
		return m.genCommonDefaults(f, name, 0, r.Float.Default, wk)
	case *flags.FieldFlags_Double:
		return m.genCommonDefaults(f, name, 0.0, r.Double.Default, wk)
	case *flags.FieldFlags_Int32:
		return m.genCommonDefaults(f, name, 0, r.Int32.Default, wk)
	case *flags.FieldFlags_Int64:
		return m.genCommonDefaults(f, name, 0, r.Int64.Default, wk)
	case *flags.FieldFlags_Uint32:
		return m.genCommonDefaults(f, name, 0, r.Uint32.Default, wk)
	case *flags.FieldFlags_Uint64:
		return m.genCommonDefaults(f, name, 0, r.Uint64.Default, wk)
	case *flags.FieldFlags_Sint32:
		return m.genCommonDefaults(f, name, 0, r.Sint32.Default, wk)
	case *flags.FieldFlags_Sint64:
		return m.genCommonDefaults(f, name, 0, r.Sint64.Default, wk)
	case *flags.FieldFlags_Fixed32:
		return m.genCommonDefaults(f, name, 0, r.Fixed32.Default, wk)
	case *flags.FieldFlags_Fixed64:
		return m.genCommonDefaults(f, name, 0, r.Fixed64.Default, wk)
	case *flags.FieldFlags_Sfixed32:
		return m.genCommonDefaults(f, name, 0, r.Sfixed32.Default, wk)
	case *flags.FieldFlags_Sfixed64:
		return m.genCommonDefaults(f, name, 0, r.Sfixed64.Default, wk)
	case *flags.FieldFlags_Bool:
		return m.genCommonDefaults(f, name, false, r.Bool.Default, wk)
	case *flags.FieldFlags_String_:
		if r.String_.Default != nil {
			return m.genCommonDefaults(f, name, `""`, fmt.Sprint(`"`, r.String_.GetDefault(), `"`), wk)
		}
		return ""
	case *flags.FieldFlags_Bytes:
		return m.genBytesDefaults(f, name, r.Bytes, wk)
	case *flags.FieldFlags_Enum:
		if r.Enum.Default != nil {
			return m.genCommonDefaults(f, name, 0, r.Enum.GetDefault(), wk)
		}
		return ""
	case *flags.FieldFlags_Duration:
		return m.genDurationDefaults(f, name, r.Duration)
	case *flags.FieldFlags_Timestamp:
		return m.genTimestampDefaults(f, name, r.Timestamp)
	case *flags.FieldFlags_Message:
		return m.genMessageDefaults(f, name, r.Message)
	case *flags.FieldFlags_Map:
		return ""
	case *flags.FieldFlags_Repeated:
		return m.processRepeatedDefaults(f, name, r.Repeated)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)
	}
	return fmt.Sprint("\n// ", f.Name())
}

// processRepeatedDefaults handles default value generation for repeated fields
func (m *Module) processRepeatedDefaults(f pgs.Field, name pgs.Name, repeated *flags.RepeatedFlags) string {
	if repeated == nil {
		return ""
	}
	wk := pgs.UnknownWKT
	if emb := f.Type().Element().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}

	switch r := repeated.Type.(type) {
	case *flags.RepeatedFlags_Bytes:
		return m.genBytesSliceDefaults(f, name, r.Bytes, wk)
	case *flags.RepeatedFlags_Float:
		return m.genCommonSliceDefaults(f, name, r.Float.GetDefault(), "%f", wk)
	case *flags.RepeatedFlags_Double:
		return m.genCommonSliceDefaults(f, name, r.Double.GetDefault(), "%f", wk)
	case *flags.RepeatedFlags_Int32:
		return m.genCommonSliceDefaults(f, name, r.Int32.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Int64:
		return m.genCommonSliceDefaults(f, name, r.Int64.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Uint32:
		return m.genCommonSliceDefaults(f, name, r.Uint32.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Uint64:
		return m.genCommonSliceDefaults(f, name, r.Uint64.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Sint32:
		return m.genCommonSliceDefaults(f, name, r.Sint32.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Sint64:
		return m.genCommonSliceDefaults(f, name, r.Sint64.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Fixed32:
		return m.genCommonSliceDefaults(f, name, r.Fixed32.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Fixed64:
		return m.genCommonSliceDefaults(f, name, r.Fixed64.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Sfixed32:
		return m.genCommonSliceDefaults(f, name, r.Sfixed32.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Sfixed64:
		return m.genCommonSliceDefaults(f, name, r.Sfixed64.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Bool:
		return m.genCommonSliceDefaults(f, name, r.Bool.GetDefault(), "%v", wk)
	case *flags.RepeatedFlags_String_:
		return m.genCommonSliceDefaults(f, name, r.String_.GetDefault(), "%q", wk)
	case *flags.RepeatedFlags_Enum:
		return m.genCommonSliceDefaults(f, name, r.Enum.GetDefault(), "%d", wk)
	case *flags.RepeatedFlags_Duration:
		return m.genDurationSliceDefaults(f, name, r.Duration, wk)
	case *flags.RepeatedFlags_Timestamp:
		return m.genTimestampSliceDefaults(f, name, r.Timestamp)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", f.Type)
	}
	return ""
}
