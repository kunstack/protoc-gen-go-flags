# Advanced Types Example

This example demonstrates all supported protobuf types and advanced features in protoc-gen-go-flags.

## Overview

This comprehensive example showcases:
- All scalar types (int32, int64, uint32, uint64, float, double, bool, string, bytes)
- Special types (Duration, Timestamp)
- Repeated fields (arrays)
- Map fields with different formats
- Enum types
- Wrapper types
- Deprecated flags
- Hidden flags

## Project Structure

```
advanced/
├── go.mod
├── main.go
├── buf.yaml
├── buf.gen.yaml
├── Makefile
├── README.md
└── proto/
    └── config.proto    # Comprehensive type examples
```

## Configuration Features

### Scalar Types
- `name` (string) - Application name
- `worker_count` (int32) - Number of worker threads
- `max_memory` (int64) - Maximum memory in bytes
- `rate_limit` (uint32) - Rate limit per second
- `ratio` (float) - Processing ratio
- `enabled` (bool) - Enable/disable feature
- `api_key` (bytes) - API key in base64

### Special Types
- `timeout` (Duration) - Connection timeout
- `created_at` (Timestamp) - Creation timestamp

### Repeated Fields
- `servers` ([]string) - List of server addresses
- `ports` ([]int32) - List of port numbers

### Map Fields
- `labels` (map[string]string) - Key-value labels
- `limits` (map[string]int32) - Resource limits

### Enum Type
- `log_level` - Log level (DEBUG, INFO, WARN, ERROR)

## Usage

### Build and Run

```bash
# Generate code and build
make build

# Run with defaults
./bin/advanced

# Run with custom configuration
./bin/advanced \
  --name "MyApp" \
  --workers 8 \
  --timeout 60s \
  --servers api1.example.com,api2.example.com \
  --labels "env=prod,region=us-west" \
  --log-level INFO

# Show help
./bin/advanced --help
```

### Example Commands

**Basic scalar types:**
```bash
./bin/advanced \
  --name "MyService" \
  --workers 4 \
  --max-memory 1073741824 \
  --rate-limit 1000 \
  --ratio 0.95 \
  --enabled
```

**Duration and Timestamp:**
```bash
./bin/advanced \
  --timeout 30s \
  --created-at "2024-01-01T00:00:00Z"
```

**Repeated fields:**
```bash
# Comma-separated values
./bin/advanced --servers "api1.com,api2.com,api3.com"

# Multiple values
./bin/advanced --ports 8080 --ports 8081 --ports 8082
```

**Map fields:**
```bash
# JSON format
./bin/advanced --labels '{"env":"prod","region":"us-west"}'

# STRING_TO_STRING format
./bin/advanced --labels "env=prod,region=us-west"

# STRING_TO_INT format
./bin/advanced --limits "cpu=1000,memory=2048"
```

**Enum:**
```bash
./bin/advanced --log-level INFO
# Accepts: DEBUG (0), INFO (1), WARN (2), ERROR (3)
```

**Bytes (Base64):**
```bash
./bin/advanced --api-key "YXBpa2V5MTIzNDU2"
```

## Type Features

### 1. Scalar Types

All Go basic types are supported with appropriate pflag methods:
- `String`, `Int32`, `Int64`, `Uint32`, `Uint64`
- `Float32`, `Float64`, `Bool`
- `BytesBase64`, `BytesHex`

### 2. Duration

Accepts time duration strings:
- `30s` - 30 seconds
- `5m` - 5 minutes
- `2h` - 2 hours
- `1h30m` - 1 hour 30 minutes

### 3. Timestamp

Supports multiple formats:
- RFC3339: `2024-01-01T15:04:05Z`
- ISO8601: `2024-01-01`
- RFC1123: `Mon, 01 Jan 2024 15:04:05 MST`

### 4. Repeated Fields

Can specify values in two ways:
- Comma-separated: `--servers "s1,s2,s3"`
- Multiple flags: `--servers s1 --servers s2 --servers s3`

### 5. Map Fields

Three format options:

**JSON** (default):
```bash
--labels '{"key": "value"}'
```

**STRING_TO_STRING**:
```bash
--labels "key1=val1,key2=val2"
```

**STRING_TO_INT**:
```bash
--limits "cpu=1000,memory=2048"
```

### 6. Enum Types

Accepts either:
- Enum name: `--log-level INFO`
- Enum number: `--log-level 1`

### 7. Deprecated Flags

Flags can be marked as deprecated with a custom message:
```protobuf
bool old_flag = 1 [(flags.value).bool = {
  deprecated: true
  deprecated_usage: "Use --new-flag instead"
}];
```

### 8. Hidden Flags

Flags can be hidden from help output:
```protobuf
string internal = 1 [(flags.value).string = {
  hidden: true
}];
```

## Next Steps

After exploring all types, check out:
- **hierarchical/** - Learn about custom prefixes and delimiters for better flag organization
