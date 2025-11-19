# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

protoc-gen-go-flags is a Go-based protocol buffer compiler plugin that generates command-line flag bindings for protobuf messages. It creates `AddFlags` methods that integrate with the `spf13/pflag` library to automatically generate CLI flags from protobuf message definitions.

## Git Commit Guidelines

**IMPORTANT**: All git commit messages MUST be written in English.

When creating commits:
- Use English for all commit messages, titles, and descriptions
- Follow conventional commit format: `type(scope): description`
- Common types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
- Use the imperative mood ("add feature" not "added feature")
- Include detailed description in the commit body when necessary

Example:
```
docs: add comprehensive integration guide to Chinese README

- Add step-by-step tutorial for project integration
- Include troubleshooting section
- Improve FAQ with practical examples
```

## Development Commands

### Build and Code Quality
```bash
# Format Go code and protobuf files
make fmt

# Run static analysis
make vet

# Run comprehensive linting
make lint

# Clean up module dependencies
make tidy

# Run tests
make test

# Clean build artifacts
make clean
```

### Code Generation
```bash
# Install required tools (buf, protoc-gen-go, golangci-lint)
make deps

# Generate Go code from protobuf definitions
make generate
```

## Architecture

### Core Components

**main.go**: Entry point that initializes the protoc-gen-star plugin framework and registers the flags module.

**module/module.go**: Contains the core `Module` struct that implements the protoc-gen-star interface. This module:
- Processes protobuf messages and fields
- Generates `AddFlags` methods using Go templates
- Maps protobuf field types to appropriate pflag methods
- Handles field options and message-level configuration

**flags/**: Defines protobuf extensions for configuring flag generation:
- `annotations.proto`: Protocol buffer extensions for field and message options
- `annotations.pb.go`: Generated protobuf Go code
- `interface.go`: Interface definition and types

### Key Design Patterns

1. **Template-based Generation**: Uses Go's `text/template` to generate flag binding code
2. **Extension-based Configuration**: Uses protobuf field and message options to control flag generation
3. **Type Mapping**: Automatically maps protobuf types to appropriate pflag methods (BoolVarP, StringVarP, etc.)
4. **Prefix Support**: Allows hierarchical flag organization through prefix parameters

### Extension Usage

The plugin supports configuration via protobuf options:

```protobuf
syntax = "proto3";
package demo;
import "flags/annotations.proto";

// Message-level options
message DemoMessage {
  bool help = 1 [(flags.value).bool = {
    name: "help"
    short: "h"
    usage: "Show help message flag"
    deprecated: true
    hidden: true
    deprecated_usage: "This flag is deprecated, use --greeting instead"
  }];
}
```

## Dependencies

- **github.com/lyft/protoc-gen-star/v2**: Protocol buffer code generation framework
- **github.com/spf13/pflag**: POSIX/GNU-style command-line flag parsing
- **google.golang.org/protobuf**: Google's protocol buffer implementation

## Testing Guidelines

When writing test cases, prioritize using the `github.com/stretchr/testify` package for assertions. Use methods like:
- `assert.NoError(t, err)` for error checking
- `assert.Equal(t, expected, actual)` for value comparisons
- `assert.Len(t, slice, length)` for length checks
- `assert.Error(t, err)` for expecting errors

This provides cleaner, more readable test code with better failure messages.

## Buf Configuration

The project uses buf for protobuf linting and generation:
- `buf.yaml`: Module configuration for buf.build/linka-cloud/protoc-gen-defaults
- `buf.gen.yaml`: Code generation configuration with source-relative paths
