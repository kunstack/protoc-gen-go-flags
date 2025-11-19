// Copyright 2021 Aapeli <aapeli.nian@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import (
	"strings"

	"github.com/spf13/pflag"
)

// Delimiter constants are used for separating hierarchical flag names.
// These provide common naming conventions for flag organization.
const (
	// DelimiterDot is the default delimiter used for separating hierarchical flag names.
	// It is set to "." (dot) to enable dot-notation flag naming (e.g., "server.port").
	DelimiterDot = "."

	// DelimiterDash is a delimiter that uses hyphen/minus character for flag naming.
	// Useful for kebab-case flag names (e.g., "server-port").
	DelimiterDash = "-"

	// DelimiterUnderscore is a delimiter that uses underscore character for flag naming.
	// Useful for snake_case flag names (e.g., "server_port").
	DelimiterUnderscore = "_"

	// DelimiterColon is a delimiter that uses colon character for flag naming.
	// Useful for namespace-style flag names (e.g., "server:port").
	DelimiterColon = ":"
)

// Options holds configuration for flag name generation and formatting.
// It contains settings for prefix handling, delimiter usage, and custom name transformations.
type Options struct {
	Prefix    []string            // Prefix segments to prepend to flag names for hierarchical organization
	Delimiter string              // Separator used between name components (default: ".")
	Renamer   func(string) string // Custom function to transform flag names after prefix application
}

// Option is a functional option pattern type that modifies Options instances.
// Functions of this type can be passed to configuration functions to customize
// flag generation behavior.
type Option func(*Options)

// Renamer is a function type that transforms flag names. It takes a flag name
// as input and returns the transformed name. Common use cases include case
// conversion (ToLower, ToUpper), adding prefixes/suffixes, or applying custom
// naming conventions.
type Renamer func(string) string

// This package provides protobuf extensions for generating AddFlags methods
// using the protoc-gen-go-flags plugin.

// Defaulter is an interface that types can implement to provide default values
// for their fields. The SetDefaults method should be called to initialize
// all fields with their configured default values.
type Defaulter interface {
	SetDefaults()
}

// Flagger is an interface that protobuf-generated structs implement to expose
// their fields as command-line flags. The AddFlags method binds the struct's
// fields to a pflag.FlagSet, allowing automatic generation of CLI flags
// from protobuf message definitions.
//
// Parameters:
//   - fs: The pflag.FlagSet to which flags will be added
//   - prefix: Optional prefix strings that will be prepended to flag names
//     for hierarchical flag organization (e.g., "server", "database")
type Flagger interface {
	AddFlags(fs *pflag.FlagSet, opts ...Option)
}

// WithDelimiter returns an Option that sets the delimiter used for separating
// hierarchical flag name components. The default delimiter is "." (dot).
//
// Example:
//
//	WithDelimiter("-") would create flags like "server-port" instead of "server.port"
func WithDelimiter(delimiter string) Option {
	return func(o *Options) {
		o.Delimiter = delimiter
	}
}

// WithRenamer returns an Option that sets a custom name transformation function.
// This function is applied to flag names after prefixes are applied but before
// the flag is registered. Common use cases include snake_case conversion,
// kebab-case conversion, or adding namespace prefixes.
//
// Example:
//
//	WithRenamer(strings.ToLower) would convert all flag names to lowercase
func WithRenamer(renamer Renamer) Option {
	return func(o *Options) {
		o.Renamer = renamer
	}
}

// WithPrefix returns an Option that adds prefix segments to flag names.
// Prefixes are useful for organizing flags hierarchically, such as by service
// or module name. Empty strings are filtered out, and trimming of delimiter
// characters is handled during the Build phase.
//
// Example:
//
//	WithPrefix("server", "database") with field "port" creates "server.database.port"
func WithPrefix(prefix ...string) Option {
	var nonEmpty []string
	for _, part := range prefix {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return func(o *Options) {
		o.Prefix = append(o.Prefix, nonEmpty...)
	}
}

// NameBuilder constructs full flag names by combining prefixes, base names,
// and applying custom transformations. It encapsulates the naming logic and
// formatting rules for generating consistent flag identifiers.
type NameBuilder struct {
	options *Options
}

// Build generates the complete flag name by joining the configured prefix
// with the provided name using the configured delimiter, then applying
// the custom renamer function if one is set.
//
// Parameters:
//   - name: The base flag name to build upon
//
// Returns:
//
//	The fully qualified flag name including prefixes and transformations
func (n NameBuilder) Build(name string) string {
	// Trim delimiter characters from prefix segments
	var trimmedPrefix []string
	for _, part := range n.options.Prefix {
		trimmed := strings.Trim(part, n.options.Delimiter)
		if trimmed != "" {
			trimmedPrefix = append(trimmedPrefix, trimmed)
		}
	}

	// Join prefix and name with delimiter, then apply renamer
	flagName := strings.Join(append(trimmedPrefix, name), n.options.Delimiter)
	if n.options.Renamer != nil {
		return n.options.Renamer(flagName)
	}
	return flagName
}

// NewNameBuilder creates a new NameBuilder with the provided configuration options.
// If no options are provided, it uses sensible defaults (dot delimiter, identity renamer).
//
// Parameters:
//   - opts: Optional configuration functions to customize naming behavior
//
// Returns:
//
//	A configured NameBuilder instance ready to generate flag names
func NewNameBuilder(opts ...Option) NameBuilder {
	options := &Options{
		Delimiter: DelimiterDot,
		Renamer:   func(s string) string { return s },
	}
	for _, opt := range opts {
		opt(options)
	}
	return NameBuilder{options: options}
}
