package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func (m *Module) checkMap(typ FieldType, flag *flags.MapFlag) {
	if flag == nil {
		return
	}

	// Ensure the field is actually a map type
	fieldType := m.mustFieldType(typ)
	m.Assert(fieldType.IsMap(), "map flag should be used for map fields")

	// Check that usage is provided
	if flag.Usage == "" {
		m.Failf("usage is required for map flag")
	}

	// Check if deprecated flag has proper deprecation usage message
	if flag.Deprecated && flag.DeprecatedUsage == "" {
		m.Failf("deprecated map flag must provide deprecated_usage message")
	}

	// Validate that the format matches the actual field types
	keyElem := fieldType.Key()
	valueElem := fieldType.Element()

	switch flag.GetFormat() {
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
		// Validate key is string
		if keyElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_STRING format requires string keys, but got %v", keyElem.ProtoType())
		}
		// Validate value is string
		if valueElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_STRING format requires string values, but got %v", valueElem.ProtoType())
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
		// Validate key is string
		if keyElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_INT format requires string keys, but got %v", keyElem.ProtoType())
		}
		// Validate value is an integer type (support all integer types)
		switch valueElem.ProtoType() {
		case pgs.Int32T, pgs.SInt32, pgs.SFixed32,
			pgs.Int64T, pgs.SInt64, pgs.SFixed64,
			pgs.UInt32T, pgs.Fixed32T,
			pgs.UInt64T, pgs.Fixed64T:
			// These are all valid integer types
		default:
			m.Failf("STRING_TO_INT format requires integer values, but got %v", valueElem.ProtoType())
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_JSON:
		// JSON format is flexible, no strict type validation needed
		// Just ensure it's actually a map
		break

	case flags.MapFormatType_MAP_FORMAT_TYPE_UNSPECIFIED:
		// Default to JSON format, no additional validation needed
		break

	default:
		m.Failf("unknown map format type: %v", flag.GetFormat())
	}
}

func (m *Module) genMap(f pgs.Field, name pgs.Name, flag *flags.MapFlag) string {
	var (
		declBuilder = &strings.Builder{}
	)

	// Validate the field type matches the specified format
	if !f.Type().IsMap() {
		m.Failf("field %s is not a map type", name)
		return ""
	}

	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	// Determine the format to use
	mapFormat := flag.GetFormat()

	// If unspecified, default to JSON format for backward compatibility
	if mapFormat == flags.MapFormatType_MAP_FORMAT_TYPE_UNSPECIFIED {
		mapFormat = flags.MapFormatType_MAP_FORMAT_TYPE_JSON
	}

	keyType := f.Type().Key()
	valueType := f.Type().Element()

	switch mapFormat {
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
		if keyType.ProtoType() != pgs.StringT {
			m.Failf("field %s key type is not string for STRING_TO_STRING format", name)
			return ""
		}
		if valueType.ProtoType() != pgs.StringT {
			m.Failf("field %s value type is not string for STRING_TO_STRING format", name)
			return ""
		}
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
		if keyType.ProtoType() != pgs.StringT {
			m.Failf("field %s key type is not string for STRING_TO_INT format", name)
			return ""
		}
		// Support all integer types for STRING_TO_INT format
		switch valueType.ProtoType() {
		case pgs.Int32T, pgs.SInt32, pgs.SFixed32,
			pgs.Int64T, pgs.SInt64, pgs.SFixed64,
			pgs.UInt32T, pgs.Fixed32T,
			pgs.UInt64T, pgs.Fixed64T:
			// These are all valid integer types
		default:
			m.Failf("field %s value type is not a valid integer type for STRING_TO_INT format", name)
			return ""
		}
	}

	// Generate flag binding based on format
	switch mapFormat {
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
		_, _ = fmt.Fprintf(declBuilder, `
				fs.StringToStringVarP(&x.%s, builder.Build(%q), %q, x.%s, %q)
			`,
			name, flag.GetName(), flag.GetShort(), name, flag.GetUsage(),
		)

	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
		// For string-to-int maps, determine the specific int type based on the field
		valueType := f.Type().Element()

		switch valueType.ProtoType() {
		case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
			_, _ = fmt.Fprintf(declBuilder, `
					fs.StringToInt64VarP(&x.%s, builder.Build(%q), %q, x.%s, %q)
				`,
				name, flag.GetName(), flag.GetShort(), name, flag.GetUsage(),
			)

		case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
			_, _ = fmt.Fprintf(declBuilder, `
					fs.VarP(types.StringToInt32(&x.%s), builder.Build(%q), %q, %q)
				`,
				name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
			)

		case pgs.UInt32T, pgs.Fixed32T:
			_, _ = fmt.Fprintf(declBuilder, `
					fs.VarP(types.StringToUint32(&x.%s), builder.Build(%q), %q, %q)
				`,
				name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
			)

		case pgs.UInt64T, pgs.Fixed64T:
			_, _ = fmt.Fprintf(declBuilder, `
					fs.VarP(types.StringToUint64(&x.%s), builder.Build(%q), %q, %q)
				`,
				name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
			)

		default:
			m.Failf("field %s value type is not a valid integer type for STRING_TO_INT format", name)
			return ""
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_JSON:
		// For JSON format, use the existing JSON handling
		_, _ = fmt.Fprintf(declBuilder, `
				fs.VarP(types.JSON(&x.%s), builder.Build(%q), %q, %q)
			`,
			name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
		)

	}
	declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
