package module

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-go-flags/flags"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

// validateBytesEncoding validates the encoding type for bytes fields.
// Returns an error if the encoding type is not supported.
func (m *Module) validateBytesEncoding(encoding flags.BytesEncodingType) error {
	switch encoding {
	case flags.BytesEncodingType_BYTES_ENCODING_TYPE_UNSPECIFIED,
		flags.BytesEncodingType_BYTES_ENCODING_TYPE_BASE64,
		flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX:
		return nil
	default:
		return fmt.Errorf("invalid bytes encoding type: %v", encoding)
	}
}

// validateBytesDefault validates a default bytes value based on the encoding type.
// If index is >= 0, it's for repeated fields and will be used in error messages.
func (m *Module) validateBytesDefault(data []byte, encoding flags.BytesEncodingType, index int) error {
	if len(data) == 0 {
		return nil
	}

	switch encoding {
	case flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX:
		// Verify the default value is valid hexadecimal
		if _, err := hex.DecodeString(string(data)); err != nil {
			if index >= 0 {
				return fmt.Errorf("bytes default value at index %d is not valid hex: %v", index, err)
			}
			return fmt.Errorf("bytes default value is not valid hex: %v", err)
		}

	case flags.BytesEncodingType_BYTES_ENCODING_TYPE_BASE64, flags.BytesEncodingType_BYTES_ENCODING_TYPE_UNSPECIFIED:
		// Verify the default value is valid base64
		if _, err := base64.StdEncoding.DecodeString(string(data)); err != nil {
			if index >= 0 {
				return fmt.Errorf("bytes default value at index %d is not valid base64: %v", index, err)
			}
			return fmt.Errorf("bytes default value is not valid base64: %v", err)
		}
	}

	return nil
}

// checkBytes validates the configuration of a bytes flag field.
// It performs comprehensive validation including encoding type, default values,
// and required fields for bytes-type protobuf fields.
func (m *Module) checkBytes(ft pgs.FieldType, r *flags.BytesFlag) {
	// Perform common validation first (name, type compatibility, etc.)
	m.checkCommon(ft, r, pgs.BytesT, pgs.BytesValueWKT, false)

	// Ensure usage description is provided for the flag
	if r.GetUsage() == "" {
		m.Failf("usage is required for bytes flag")
	}

	// Validate the encoding type specification using the common method
	if err := m.validateBytesEncoding(r.GetEncoding()); err != nil {
		m.Failf(err.Error())
	}

	// Validate default value if provided - must be valid for the specified encoding
	if r.Default != nil && len(r.GetDefault()) > 0 {
		if err := m.validateBytesDefault(r.GetDefault(), r.GetEncoding(), -1); err != nil {
			m.Failf(err.Error())
		}
	}

	// Ensure deprecated flags have proper deprecation messages
	if r.GetDeprecated() && r.GetDeprecatedUsage() == "" {
		m.Failf("deprecated bytes flag must provide deprecated_usage message")
	}
}

// checkBytesSlice validates the configuration of a repeated bytes flag field.
// Similar to checkBytes but handles slice/repeated fields with multiple default values.
func (m *Module) checkBytesSlice(ft FieldType, r *flags.RepeatedBytesFlag) {
	// Perform common validation first for repeated fields (name, type compatibility, etc.)
	m.checkCommon(ft, r, pgs.BytesT, pgs.BytesValueWKT, true)

	// Ensure usage description is provided for the repeated flag
	if r.GetUsage() == "" {
		m.Failf("usage is required for repeated bytes flag")
	}

	// Validate the encoding type specification using the common method
	if err := m.validateBytesEncoding(r.GetEncoding()); err != nil {
		m.Failf(err.Error())
	}

	// Validate each default value if provided - must be valid for the specified encoding
	for i, defaultBytes := range r.Default {
		if err := m.validateBytesDefault(defaultBytes, r.GetEncoding(), i); err != nil {
			m.Failf(err.Error())
		}
	}

	// Ensure deprecated flags have proper deprecation messages
	if r.GetDeprecated() && r.GetDeprecatedUsage() == "" {
		m.Failf("deprecated repeated bytes flag must provide deprecated_usage message")
	}
}

// genBytes generates the flag binding code for a bytes field.
// It handles both regular bytes fields and google.protobuf.BytesValue wrapper types,
// supporting both base64 and hex encoding formats.
//
// Parameters:
//   - f: The protobuf field for type information
//   - name: The field name for code generation
//   - flag: The bytes flag configuration
//   - wk: Well-known type information (e.g., google.protobuf.BytesValue)
func (m *Module) genBytes(f pgs.Field, name pgs.Name, flag *flags.BytesFlag, wk pgs.WellKnownType) string {
	// Configure the flag and check if it's disabled
	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	var (
		wrapper       = "Bytes"
		nativeWrapper = "BytesBase64VarP"
		declBuilder   = &strings.Builder{}
	)

	// Set wrapper based on encoding
	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		wrapper = "BytesHex"
		nativeWrapper = "BytesHexVarP"
	}

	// Handle google.protobuf.BytesValue wrapper types
	if wk != "" && wk != pgs.UnknownWKT {
		_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
			name, name, m.getFieldTypeName(f),
		)
		_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.%s(x.%s), builder.Build(%q), %q, %q)
		`,
			wrapper, name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
		)
	} else {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.%s(&x.%s, builder.Build(%q), %q, x.%s, %q)
			`,
			nativeWrapper, name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
	}
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

// genBytesSlice generates the flag binding code for a repeated bytes field.
// Creates a slice-based flag that can accept multiple values, supporting
// both base64 and hex encoding formats.
//
// Parameters:
//   - f: The protobuf field (unused but kept for interface consistency)
//   - name: The field name for code generation
//   - flag: The repeated bytes flag configuration
//   - wk: Well-known type information (unused for slice types)
func (m *Module) genBytesSlice(name pgs.Name, flag *flags.RepeatedBytesFlag) string {
	// Configure the flag and check if it's disabled
	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	var (
		wrapper     = "BytesSlice"
		declBuilder = &strings.Builder{}
	)

	// Set wrapper based on encoding
	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		wrapper = "BytesHexSlice"
	}

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.%s(&x.%s), builder.Build(%q), %q, %q)
		`,
		wrapper, name, flag.GetName(), flag.GetShort(), flag.GetUsage())

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

// genBytesDefaults generates the default value assignment code for bytes fields.
// It handles both regular bytes fields and google.protobuf.BytesValue wrapper types,
// with support for hex and base64 encoding formats.
//
// Parameters:
//   - f: The protobuf field (unused but kept for interface consistency)
//   - name: The field name for code generation
//   - flag: The bytes flag configuration
//   - wk: Well-known type information (e.g., google.protobuf.BytesValue)
func (m *Module) genBytesDefaults(f pgs.Field, name pgs.Name, flag *flags.BytesFlag, wk pgs.WellKnownType) string {
	// Return empty string if no default value is configured
	if flag.Default == nil || len(flag.GetDefault()) == 0 {
		return ""
	}

	fieldName := name.String()
	defaultBytes := flag.GetDefault()
	encoding := flag.GetEncoding()
	isWrapper := wk != "" && wk != pgs.UnknownWKT

	switch encoding {
	case flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX:
		if isWrapper {
			return fmt.Sprintf(`
			if x.%s == nil {
				x.%s = &wrapperspb.BytesValue{Value:  utils.MustDecodeHex(%q)}
			}`, fieldName, fieldName, defaultBytes)
		}
		return fmt.Sprintf(`
			if len(x.%s) == 0 {
				x.%s =  utils.MustDecodeHex(%q)
			}`, fieldName, fieldName, defaultBytes)

	case flags.BytesEncodingType_BYTES_ENCODING_TYPE_BASE64, flags.BytesEncodingType_BYTES_ENCODING_TYPE_UNSPECIFIED:
		if isWrapper {
			return fmt.Sprintf(`
			if x.%s == nil {
				x.%s = &wrapperspb.BytesValue{Value:  utils.MustDecodeBase64(%q)}
			}`, fieldName, fieldName, defaultBytes)
		}
		return fmt.Sprintf(`
			if len(x.%s) == 0 {
				x.%s =  utils.MustDecodeBase64(%q)
			}`, fieldName, fieldName, defaultBytes)
	}
	return ""
}

// genBytesSliceDefaults generates the default value assignment code for repeated bytes fields.
// It handles both regular bytes slice fields and google.protobuf.BytesValue wrapper types,
// with support for hex and base64 encoding formats.
//
// Parameters:
//   - f: The protobuf field (unused but kept for interface consistency)
//   - name: The field name for code generation
//   - flag: The repeated bytes flag configuration
//   - wk: Well-known type information (unused for slice types)
func (m *Module) genBytesSliceDefaults(f pgs.Field, name pgs.Name, flag *flags.RepeatedBytesFlag, wk pgs.WellKnownType) string {
	// Return empty string if no default value is configured
	if flag.Default == nil || len(flag.GetDefault()) == 0 {
		return ""
	}

	var code strings.Builder

	// Check if the slice is empty before setting defaults (similar to genCommonDefaults)
	code.WriteString(fmt.Sprintf(`
		if len(x.%s) == 0 {`, name))

	defaultValues := make([]string, len(flag.GetDefault()))

	// Generate default assignments for each value in the slice
	for i, defaultBytes := range flag.Default {
		switch flag.GetEncoding() {
		case flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX:
			if wk != "" && wk != pgs.UnknownWKT {
				defaultValues[i] = fmt.Sprintf("{Value:  utils.MustDecodeHex(%q) }", defaultBytes)
			} else {
				defaultValues[i] = fmt.Sprintf(" utils.MustDecodeHex(%q)", defaultBytes)
			}
		case flags.BytesEncodingType_BYTES_ENCODING_TYPE_BASE64, flags.BytesEncodingType_BYTES_ENCODING_TYPE_UNSPECIFIED:
			if wk != "" && wk != pgs.UnknownWKT {
				defaultValues[i] = fmt.Sprintf("{Value:  utils.MustDecodeBase64(%q) }", defaultBytes)
			} else {
				defaultValues[i] = fmt.Sprintf(" utils.MustDecodeBase64(%q)", defaultBytes)
			}
		}
	}

	// Append the decoded bytes to the slice
	code.WriteString(fmt.Sprintf(`
			x.%s = %s{%s}`, name, m.getFieldTypeName(f), strings.Join(defaultValues, ",")))

	code.WriteString(`
		}`)

	return code.String()
}
