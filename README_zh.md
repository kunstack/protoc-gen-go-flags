# protoc-gen-go-flags

[![Go Report Card](https://goreportcard.com/badge/github.com/kunstack/protoc-gen-go-flags)](https://goreportcard.com/report/github.com/kunstack/protoc-gen-go-flags)
[![Go Reference](https://pkg.go.dev/badge/github.com/kunstack/protoc-gen-go-flags.svg)](https://pkg.go.dev/github.com/kunstack/protoc-gen-go-flags)

中文文档 | [English](README.md)

protoc-gen-go-flags 是一个基于 Go 语言的 Protocol Buffer 编译器插件，用于为 protobuf 消息自动生成命令行标志绑定。它能够根据 protobuf 消息定义自动生成 `AddFlags` 方法，与 `spf13/pflag` 库无缝集成，为您的 protobuf 消息提供强大的命令行参数支持。

## 为什么使用 protoc-gen-go-flags

如果您的项目满足以下任一条件，protoc-gen-go-flags 将大大简化您的开发工作：

- ✅ 使用 Protocol Buffers 定义配置结构
- ✅ 需要为 CLI 应用提供命令行参数支持
- ✅ 希望避免手动编写重复的标志绑定代码
- ✅ 想要保持配置定义和 CLI 接口的一致性
- ✅ 需要支持复杂的嵌套配置结构

**传统方式 vs protoc-gen-go-flags：**

传统方式需要为每个配置字段手动编写标志绑定：
```go
// 手动方式：繁琐且容易出错
fs.StringVar(&config.Host, "host", "localhost", "Server host")
fs.Int32Var(&config.Port, "port", 8080, "Server port")
fs.BoolVar(&config.Verbose, "verbose", false, "Enable verbose")
// ... 为每个字段重复编写
```

使用 protoc-gen-go-flags：
```go
// 自动生成：简洁且类型安全
config.AddFlags(fs)
```

## 目录

- [为什么使用 protoc-gen-go-flags](#为什么使用-protoc-gen-go-flags)
- [特性](#特性)
- [快速开始](#快速开始)
  - [前置要求](#前置要求)
  - [安装](#安装)
  - [基本用法](#基本用法)
- [完整集成教程](#完整集成教程)
  - [步骤 1：准备项目](#步骤-1准备项目)
  - [步骤 2：添加标志注解依赖](#步骤-2添加标志注解依赖)
  - [步骤 3：定义 Protobuf 消息](#步骤-3定义-protobuf-消息)
  - [步骤 4：生成代码](#步骤-4生成代码)
  - [步骤 5：在应用中使用](#步骤-5在应用中使用)
- [使用示例](#使用示例)
- [支持的类型](#支持的类型)
- [配置选项](#配置选项)
- [分层标志组织](#分层标志组织)
- [常见问题](#常见问题)
- [贡献](#贡献)
- [许可证](#许可证)

## 特性

- 🚀 **自动化代码生成**：从 protobuf 消息自动生成命令行标志绑定
- 🎯 **类型全覆盖**：支持所有 protobuf 类型（标量类型、枚举、repeated、map、消息等）
- 🔧 **高度可配置**：支持自定义标志名称、简写、用法文本、默认值等
- 📦 **嵌套消息支持**：为嵌套消息生成层级化标志
- 🏗️ **分层组织**：通过前缀支持分层标志命名（支持点号、破折号、下划线、冒号分隔符）
- 🔒 **最佳实践**：生成符合 Go 规范的代码，支持私有/公有方法
- 💾 **默认值支持**：为所有类型提供默认值设置
- 🚦 **废弃标志**：支持废弃标志和隐藏标志
- 🔄 **包别名**：智能处理包名冲突，避免编译错误

## 快速开始

### 前置要求

在开始之前，请确保您的开发环境满足以下要求：

- **Go 1.18+**：protoc-gen-go-flags 需要 Go 1.18 或更高版本
- **Protocol Buffers 编译器（protoc）**：用于编译 .proto 文件
  ```bash
  # macOS
  brew install protobuf

  # Ubuntu/Debian
  apt-get install protobuf-compiler

  # 或从官方下载: https://github.com/protocolbuffers/protobuf/releases
  ```
- **protoc-gen-go**：Go 的 protobuf 代码生成器
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```

### 安装

安装 protoc-gen-go-flags 插件：

```bash
go install github.com/kunstack/protoc-gen-go-flags@latest
```

验证安装：
```bash
protoc-gen-go-flags --version
```

### 基本用法

**1. 定义带有标志选项的 protobuf 消息：**

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";

option go_package = "github.com/example/project;example";

message Config {
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: 8080
    }];

    bool verbose = 3 [(flags.value).bool = {
        name: "verbose"
        short: "v"
        usage: "Enable verbose logging"
    }];
}
```

**2. 生成代码：**

```bash
protoc -I. -I flags --go_out=paths=source_relative:. --flags_out=paths=source_relative:. config.proto
```

**3. 在应用中使用：**

```go
package main

import (
    "fmt"
    "os"

    pb "github.com/example/project"
    "github.com/spf13/pflag"
)

func main() {
    var config pb.Config

    // 创建标志集并添加标志
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
    config.AddFlags(fs)

    // 解析标志
    fs.Parse(os.Args[1:])

    // 使用配置（直接访问字段）
    fmt.Printf("Server: %s:%d (verbose: %v)\n",
        config.Host, config.Port, config.Verbose)
}
```

### AddFlags vs SetDefaults 的区别

- **AddFlags 方法**：将配置字段注册为命令行标志，让用户可以通过 CLI 参数传入值
- **SetDefaults 方法**：设置字段的默认值，在没有用户提供参数时使用

**使用场景**：
- 如果用户希望默认值在标志解析之前就生效，应该先调用 `SetDefaults()`
- 如果只是希望从命令行读取配置，可以只使用 `AddFlags()`
- 最佳实践是两者结合使用，既提供默认值，又允许用户覆盖

**调用示例**：

```go
var config pb.Config

// 方法1：只使用 AddFlags（用户必须提供所有值）
config.AddFlags(fs)

// 方法2：结合使用（推荐）
config.SetDefaults()  // 先设置默认值
config.AddFlags(fs)   // 再添加标志覆盖

// 方法3：在自定义标志集中使用
customFS := pflag.NewFlagSet("custom", pflag.ExitOnError)
config.AddFlags(customFS)
```

## 集成教程

本节提供完整的分步教程，帮助您在自己的项目中集成 protoc-gen-go-flags。

### 步骤 1：准备项目

创建一个新的 Go 项目（或使用现有项目）：

```bash
mkdir myapp
cd myapp
go mod init github.com/yourname/myapp
```

安装必要的依赖：

```bash
# 安装 pflag 库
go get github.com/spf13/pflag

# 安装 protobuf 运行时
go get google.golang.org/protobuf

# 安装 protoc-gen-go-flags 运行时库
go get github.com/kunstack/protoc-gen-go-flags/flags
```

创建项目结构：

```bash
myapp/
├── go.mod
├── main.go          # 应用入口
└── proto/
    └── config.proto # protobuf 定义
```

### 步骤 2：添加标志注解依赖

您需要将 protoc-gen-go-flags 的注解文件添加到您的项目中。有两种方式：

#### 方式 1：使用 Buf Schema Registry（推荐）

在您的 `buf.yaml` 中添加依赖：

```yaml
version: v2
deps:
  - buf.build/kunstack/protoc-gen-go-flags
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

然后在 `buf.gen.yaml` 中配置代码生成：

```yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go
    out: .
    opt: paths=source_relative
  - local: protoc-gen-go-flags
    out: .
    opt: paths=source_relative
```

运行 buf 命令更新依赖并生成代码：

```bash
buf mod update
buf generate
```

#### 方式 2：直接复制文件

从 [protoc-gen-go-flags 仓库](https://github.com/kunstack/protoc-gen-go-flags/tree/main/flags) 下载 `annotations.proto` 文件到您的项目：

```bash
mkdir -p proto/flags
curl -o proto/flags/annotations.proto \
  https://raw.githubusercontent.com/kunstack/protoc-gen-go-flags/main/flags/annotations.proto
```

项目结构更新为：

```bash
myapp/
├── go.mod
├── main.go
└── proto/
    ├── config.proto
    └── flags/
        └── annotations.proto
```

### 步骤 3：定义 Protobuf 消息

在 `proto/config.proto` 中定义您的配置：

```protobuf
syntax = "proto3";

package myapp.config;

// 导入标志注解
import "flags/annotations.proto";

option go_package = "github.com/yourname/myapp/proto;config";

message ServerConfig {
    // 启用空消息生成
    option (flags.allow_empty) = true;

    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "服务器主机地址"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "服务器端口"
        default: 8080
    }];

    bool debug = 3 [(flags.value).bool = {
        name: "debug"
        short: "d"
        usage: "启用调试模式"
    }];
}
```

### 步骤 4：生成代码

根据您在步骤2中选择的方式，使用相应的命令生成代码：

#### 使用 buf（如果选择了方式1）

如果您在步骤2中选择了 Buf Schema Registry 方式，代码已经在运行 `buf generate` 时生成。

生成的文件：
- `proto/config.pb.go` - 标准的 protobuf Go 代码
- `proto/config.pb.flags.go` - 标志绑定代码

#### 使用 protoc（如果选择了方式2）

如果您在步骤2中选择了直接复制文件方式，使用 protoc 命令生成：

```bash
protoc \
  -I./proto \
  -I./proto/flags \
  --go_out=. \
  --go_opt=paths=source_relative \
  --flags_out=. \
  --flags_opt=paths=source_relative \
  proto/config.proto
```

这将生成两个文件：
- `proto/config.pb.go` - 标准的 protobuf Go 代码
- `proto/config.pb.flags.go` - 标志绑定代码

### 步骤 5：在应用中使用

在 `main.go` 中使用生成的代码：

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/pflag"
    "github.com/yourname/myapp/proto"
)

func main() {
    // 创建配置实例
    config := &proto.ServerConfig{}

    // 设置默认值（可选但推荐）
    config.SetDefaults()

    // 创建标志集
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)

    // 添加标志
    config.AddFlags(fs)

    // 解析命令行参数
    if err := fs.Parse(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
        os.Exit(1)
    }

    // 使用配置
    fmt.Printf("Starting server...\n")
    fmt.Printf("  Host: %s\n", config.Host)
    fmt.Printf("  Port: %d\n", config.Port)
    fmt.Printf("  Debug: %v\n", config.Debug)

    // 在这里启动您的应用...
}
```

### 步骤 6：编译和运行

编译应用：

```bash
go build -o myapp
```

运行并测试命令行参数：

```bash
# 使用默认值
./myapp

# 输出：
# Starting server...
#   Host: localhost
#   Port: 8080
#   Debug: false

# 自定义参数
./myapp --host 0.0.0.0 --port 3000 --debug

# 输出：
# Starting server...
#   Host: 0.0.0.0
#   Port: 3000
#   Debug: true

# 使用短选项
./myapp -H 127.0.0.1 -p 9000 -d

# 查看帮助
./myapp --help
```

### 完整项目示例

根据您选择的方式，项目结构略有不同：

#### 使用 Buf Schema Registry 的项目结构

```bash
myapp/
├── go.mod
├── go.sum
├── main.go
├── buf.yaml              # Buf 配置
├── buf.gen.yaml          # 代码生成配置
├── buf.lock              # 依赖锁文件（生成）
├── Makefile              # 可选：自动化构建
└── proto/
    ├── config.proto
    ├── config.pb.go          # 生成
    └── config.pb.flags.go    # 生成
```

**Makefile 示例**（使用 buf）：

```makefile
.PHONY: generate build run clean

# 生成 protobuf 代码
generate:
	buf mod update
	buf generate

# 构建应用
build: generate
	go build -o bin/myapp .

# 运行应用
run: build
	./bin/myapp

# 清理生成的文件
clean:
	rm -f proto/*.pb.go proto/*.pb.flags.go
	rm -rf bin/
```

#### 使用直接复制文件的项目结构

```bash
myapp/
├── go.mod
├── go.sum
├── main.go
├── Makefile              # 可选：自动化构建
└── proto/
    ├── config.proto
    ├── config.pb.go          # 生成
    ├── config.pb.flags.go    # 生成
    └── flags/
        └── annotations.proto
```

**Makefile 示例**（使用 protoc）：

```makefile
.PHONY: generate build run clean

# 生成 protobuf 代码
generate:
	protoc \
	  -I./proto \
	  -I./proto/flags \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --flags_out=. \
	  --flags_opt=paths=source_relative \
	  proto/*.proto

# 构建应用
build: generate
	go build -o bin/myapp .

# 运行应用
run: build
	./bin/myapp

# 清理生成的文件
clean:
	rm -f proto/*.pb.go proto/*.pb.flags.go
	rm -rf bin/
```

使用 Makefile：

```bash
# 生成代码
make generate

# 构建
make build

# 运行
make run

# 清理
make clean
```

### 高级集成：嵌套配置

对于复杂的应用，您可能需要嵌套配置：

```protobuf
syntax = "proto3";

package myapp.config;

import "flags/annotations.proto";

option go_package = "github.com/yourname/myapp/proto;config";

message DatabaseConfig {
    string host = 1 [(flags.value).string = {
        name: "db-host"
        usage: "数据库主机"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "db-port"
        usage: "数据库端口"
        default: 5432
    }];
}

message AppConfig {
    option (flags.allow_empty) = true;

    string app_name = 1 [(flags.value).string = {
        name: "app-name"
        usage: "应用名称"
        default: "MyApp"
    }];

    // 嵌套配置
    DatabaseConfig database = 2 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

使用嵌套配置：

```go
config := &proto.AppConfig{}
config.SetDefaults()

fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs)
fs.Parse(os.Args[1:])

fmt.Printf("App: %s\n", config.AppName)
fmt.Printf("DB: %s:%d\n", config.Database.Host, config.Database.Port)
```

命令行使用：

```bash
./myapp --app-name "MyService" --db-db-host db.example.com --db-db-port 3306
```

### 故障排除

#### 问题 1：找不到 annotations.proto

**错误信息**：
```
proto/config.proto:3:1: Import "flags/annotations.proto" was not found.
```

**解决方案**：

- **使用 buf 方式**：确保运行了 `buf mod update` 并且 `buf.yaml` 中正确配置了依赖：
  ```yaml
  deps:
    - buf.build/kunstack/protoc-gen-go-flags
  ```

- **使用 protoc 方式**：确保在 protoc 命令中包含正确的导入路径：
  ```bash
  protoc -I./proto -I./proto/flags ...
  ```

#### 问题 2：生成的代码编译错误

**错误信息**：
```
undefined: flags.Option
undefined: types.Duration
```

**解决方案**：确保已安装运行时库：
```bash
go get github.com/kunstack/protoc-gen-go-flags/flags
go get github.com/kunstack/protoc-gen-go-flags/types
go get github.com/kunstack/protoc-gen-go-flags/utils
```

生成的代码会自动导入这些包，无需手动导入。

#### 问题 3：标志未生效

**现象**：命令行参数没有被读取，配置使用的是零值。

**原因**：可能忘记调用 `SetDefaults()` 或 `AddFlags()`。

**解决方案**：按正确顺序调用：
```go
config := &proto.ServerConfig{}
config.SetDefaults()  // 1. 设置默认值
config.AddFlags(fs)   // 2. 注册标志
fs.Parse(os.Args[1:]) // 3. 解析参数
```

#### 问题 4：buf generate 失败

**错误信息**：
```
Failure: plugin flags: not found
```

**解决方案**：确保 protoc-gen-go-flags 已安装并在 PATH 中：
```bash
# 安装插件
go install github.com/kunstack/protoc-gen-go-flags@latest

# 验证安装
which protoc-gen-go-flags
protoc-gen-go-flags --version
```

#### 问题 5：包名冲突

**现象**：生成的代码中有包名冲突，例如同时使用 `wrapperspb` 和自定义的 `wrapperspb` 包。

**解决方案**：protoc-gen-go-flags 会自动处理包名冲突，为冲突的包生成别名。生成的代码会自动使用别名导入，无需手动处理。

#### 问题 6：Map 格式解析错误

**错误信息**：
```
invalid map format: ...
```

**解决方案**：确保使用正确的格式：
- **JSON 格式**：`--config='{"key": "value"}'`
- **STRING_TO_STRING**：`--labels="key1=value1,key2=value2"`
- **STRING_TO_INT**：`--limits="cpu=1000,memory=2048"`

注意：STRING_TO_INT 只支持整数类型的值。

## 使用示例

### 基本配置示例

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";

option go_package = "github.com/example/project;example";

message ServerConfig {
    option (flags.allow_empty) = true;

    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server host address"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port number"
        default: 8080
    }];

    bool https = 3 [(flags.value).bool = {
        name: "https"
        short: "s"
        usage: "Enable HTTPS"
    }];
}
```

### 分层标志（使用前缀）

```go
// 生成带前缀的标志
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs, flags.WithPrefix("server"))
fs.Parse(os.Args[1:])

// 结果：
// --server.host
// --server.port
// --server.https
```

### 自定义分隔符

```go
// 使用破折号分隔符
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs,
    flags.WithPrefix("server"),
    flags.WithDelimiter("-"))

// 结果：
// --server-host
// --server-port
// --server-https
```

### 嵌套消息

```protobuf
message DatabaseConfig {
    string url = 1 [(flags.value).string = {
        name: "database-url"
        usage: "Database connection URL"
    }];
}

message AppConfig {
    DatabaseConfig database = 1 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

生成的标志：
- `--db.database-url`

### 完整配置示例

```protobuf
syntax = "proto3";

package example;

import "flags/annotations.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/duration.proto";

option go_package = "github.com/example/project;example";

message Config {
    option (flags.allow_empty) = true;

    // 基础类型
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: 8080
    }];

    // 特殊类型
    google.protobuf.Duration timeout = 3 [(flags.value).duration = {
        name: "timeout"
        short: "t"
        usage: "Connection timeout"
        default: "30s"
    }];

    google.protobuf.Timestamp created = 4 [(flags.value).timestamp = {
        name: "created"
        usage: "Creation time"
        formats: ["RFC3339", "ISO8601"]
        default: "2024-01-01T00:00:00Z"
    }];

    // 重复字段
    repeated string servers = 5 [(flags.value).repeated.string = {
        name: "servers"
        short: "s"
        usage: "Server addresses"
        default: ["localhost:8080"]
    }];

    // 映射字段
    map<string, int32> limits = 6 [(flags.value).map = {
        name: "limits"
        usage: "Resource limits"
        format: MAP_FORMAT_TYPE_STRING_TO_INT
        default: "{\"cpu\": 1000, \"memory\": 2048}"
    }];

    // 嵌套消息
    DatabaseConfig database = 7 [(flags.value).message = {
        name: "database"
        nested: true
    }];
}
```

## 支持的类型

protoc-gen-go-flags 支持所有 Protocol Buffer 类型：

### 标量类型

| 类型 | Go 类型 | 默认值支持 | 重复字段支持 | 示例 |
|------|---------|------------|--------------|------|
| `float` | `float32` | ✅ | ✅ | `3.14159` |
| `double` | `float64` | ✅ | ✅ | `2.71828` |
| `int32` | `int32` | ✅ | ✅ | `42` |
| `int64` | `int64` | ✅ | ✅ | `9223372036854775807` |
| `uint32` | `uint32` | ✅ | ✅ | `1000` |
| `uint64` | `uint64` | ✅ | ✅ | `18446744073709551615` |
| `sint32` | `int32` | ✅ | ✅ | `-42` |
| `sint64` | `int64` | ✅ | ✅ | `-9223372036854775808` |
| `fixed32` | `uint32` | ✅ | ✅ | `8080` |
| `fixed64` | `uint64` | ✅ | ✅ | `3000000000` |
| `sfixed32` | `int32` | ✅ | ✅ | `-1000` |
| `sfixed64` | `int64` | ✅ | ✅ | `-3000000000` |
| `bool` | `bool` | ✅ | ✅ | `true`, `false` |
| `string` | `string` | ✅ | ✅ | `"hello world"` |
| `bytes` | `[]byte` | ✅ | ✅ | `"aGVsbG8="` (base64) |

### 特殊类型

| 类型 | Go 类型 | 特性 | 示例 |
|------|---------|------|------|
| `enum` | 枚举类型 | 默认值支持，重复字段 | `1` (枚举值) |
| `google.protobuf.Duration` | `*durationpb.Duration` | 默认值支持，重复字段 | `"30s"`, `"1h"` |
| `google.protobuf.Timestamp` | `*timestamppb.Timestamp` | 多种格式，默认值支持，重复字段 | `"2024-01-01T00:00:00Z"` |
| `google.protobuf.StringValue` | `*wrapperspb.StringValue` | 默认值支持，重复字段 | `"wrapper"` |
| `google.protobuf.Int32Value` | `*wrapperspb.Int32Value` | 默认值支持，重复字段 | `42` |
| `google.protobuf.BoolValue` | `*wrapperspb.BoolValue` | 默认值支持，重复字段 | `true` |

### 复合类型

| 类型 | 格式支持 | 默认值支持 | 示例 |
|------|----------|------------|------|
| `repeated` (所有标量类型) | - | ✅ | 切片类型 |
| `map<string, string>` | JSON, 原生 | ✅ | `{"key": "value"}` |
| `map<string, int32>` | JSON, 原生 | ✅ | `{"key": 123}` |
| `map<string, int64>` | JSON, 原生 | ✅ | `{"key": 456}` |

### 嵌套消息

支持为嵌套消息生成层级化标志，通过 `message` 标志类型配置。

## 配置选项

### 消息级选项

消息级选项控制整个消息的标志生成行为：

```protobuf
message MyMessage {
  // 禁用标志生成
  option (flags.disabled) = true;

  // 生成未导出的标志方法（用于自定义包装）
  option (flags.unexported) = true;

  // 即使没有字段配置也允许生成标志方法
  option (flags.allow_empty) = true;

  // 字段定义...
}
```

| 选项 | 类型 | 描述 |
|------|------|------|
| `flags.disabled` | `bool` | 跳过为此消息生成标志 |
| `flags.unexported` | `bool` | 生成未导出的标志方法 |
| `flags.allow_empty` | `bool` | 即使没有字段配置也生成方法 |

### 字段级选项

字段级选项为单个字段提供详细配置：

```protobuf
string name = 1 [(flags.value).string = {
  name: "custom-name"           // 自定义标志名
  short: "n"                    // 短标志（单字符）
  usage: "Usage text"           // 用法说明
  hidden: false                 // 隐藏标志（不在帮助中显示）
  deprecated: true              // 标记为废弃
  deprecated_usage: "Use --new-flag instead" // 废弃说明
  default: "default-value"      // 默认值
}];
```

#### 通用字段选项

所有字段类型都支持以下选项：

| 选项 | 类型 | 描述 |
|------|------|------|
| `name` | `string` | 自定义标志名（默认为字段名） |
| `short` | `string` | 短标志别名（单字符） |
| `usage` | `string` | 帮助文本（必填） |
| `hidden` | `bool` | 隐藏标志 |
| `deprecated` | `bool` | 废弃标志 |
| `deprecated_usage` | `string` | 废弃说明（废弃标志必填） |

#### 字节类型（bytes）

字节类型支持编码格式选择：

```protobuf
bytes data = 1 [(flags.value).bytes = {
  name: "data"
  usage: "Binary data"
  encoding: BYTES_ENCODING_TYPE_BASE64  // 或 BYTES_ENCODING_TYPE_HEX
  default: "aGVsbG8="
}];
```

支持的编码：
- `BYTES_ENCODING_TYPE_BASE64` - 标准 base64 编码（默认）
- `BYTES_ENCODING_TYPE_HEX` - 十六进制编码

#### 时间戳类型（timestamp）

时间戳类型支持多种时间格式：

```protobuf
google.protobuf.Timestamp created_at = 1 [(flags.value).timestamp = {
  name: "created-at"
  usage: "Creation timestamp"
  formats: ["RFC3339", "ISO8601"]  // 支持的格式
  default: "2024-01-01T00:00:00Z"
}];
```

支持的时间格式：

| 格式名称 | Go 时间格式常量 | 示例 | 说明 |
|---------|----------------|------|------|
| `RFC3339` | `time.RFC3339` | `2024-01-01T15:04:05Z` 或 `2024-01-01T15:04:05+08:00` | RFC 3339 标准格式（推荐） |
| `RFC3339Nano` | `time.RFC3339Nano` | `2024-01-01T15:04:05.999999999Z` | RFC 3339 纳秒精度 |
| `RFC822` | `time.RFC822` | `01 Jan 24 15:04 MST` | RFC 822 格式 |
| `RFC822Z` | `time.RFC822Z` | `01 Jan 24 15:04 -0700` | RFC 822 带数字时区 |
| `RFC850` | `time.RFC850` | `Monday, 01-Jan-24 15:04:05 MST` | RFC 850 格式 |
| `RFC1123` | `time.RFC1123` | `Mon, 01 Jan 2024 15:04:05 MST` | RFC 1123 格式 |
| `RFC1123Z` | `time.RFC1123Z` | `Mon, 01 Jan 2024 15:04:05 -0700` | RFC 1123 带数字时区 |
| `ISO8601` | 自定义 | `2024-01-01` | ISO 8601 日期格式 |
| `ISO8601Time` | 自定义 | `2024-01-01T15:04:05` | ISO 8601 日期时间格式（无时区） |
| `Kitchen` | `time.Kitchen` | `3:04PM` | 厨房时钟格式 |
| `Stamp` | `time.Stamp` | `Jan  1 15:04:05` | 时间戳格式 |
| `StampMilli` | `time.StampMilli` | `Jan  1 15:04:05.000` | 毫秒精度时间戳 |
| `StampMicro` | `time.StampMicro` | `Jan  1 15:04:05.000000` | 微秒精度时间戳 |
| `StampNano` | `time.StampNano` | `Jan  1 15:04:05.000000000` | 纳秒精度时间戳 |
| `DateTime` | 自定义 | `2024-01-01 15:04:05` | 日期时间格式 |
| `DateOnly` | 自定义 | `2024-01-01` | 仅日期格式 |
| `TimeOnly` | 自定义 | `15:04:05` | 仅时间格式 |

> **注意**：
> - 除了预定义格式外，您还可以使用任何有效的 Go 时间格式字符串
> - 格式名称支持 `RFC339`（拼写错误）作为 `RFC3339` 的别名，向后兼容
> - 在 `formats` 数组中可以指定多个格式，解析时会按顺序尝试匹配

#### 持续时间类型（duration）

```protobuf
google.protobuf.Duration timeout = 1 [(flags.value).duration = {
  name: "timeout"
  usage: "Timeout duration"
  default: "30s"
}];
```

支持格式：秒数+单位（如 "30s", "5m", "1h"）

#### 映射类型（map）

映射类型支持多种格式，根据格式类型提供相应的默认值：

**1. JSON 格式（默认）**

```protobuf
map<string, int32> config = 1 [(flags.value).map = {
  name: "config"
  usage: "Configuration key-value pairs"
  format: MAP_FORMAT_TYPE_JSON
  default: "{\"cpu\": 1000, \"memory\": 2048}"
}];
```

命令行使用示例：
```bash
# JSON 格式输入
./myapp --config='{"cpu": 1000, "memory": 2048}'
```

**2. STRING_TO_STRING 格式**

```protobuf
map<string, string> labels = 1 [(flags.value).map = {
  name: "labels"
  usage: "Key-value labels"
  format: MAP_FORMAT_TYPE_STRING_TO_STRING
  default: "env=production,region=us-west"
}];
```

命令行使用示例：
```bash
# 使用逗号分隔的键值对
./myapp --labels="env=production,region=us-west"

# 或分多次指定（会覆盖）
./myapp --labels="env=staging" --labels="region=eu-central"
```

**3. STRING_TO_INT 格式**

```protobuf
map<string, int32> limits = 1 [(flags.value).map = {
  name: "limits"
  usage: "Resource limits"
  format: MAP_FORMAT_TYPE_STRING_TO_INT
  default: "cpu=1000,memory=2048,disk=10000"
}];
```

命令行使用示例：
```bash
# 使用逗号分隔的整数键值对
./myapp --limits="cpu=1000,memory=2048,disk=10000"

# 单个键值对
./myapp --limits="cpu=2000"

# 多次指定会合并
./myapp --limits="cpu=1000,memory=2048" --limits="disk=10000"
```

支持的格式：
- `MAP_FORMAT_TYPE_JSON` - JSON 格式（默认）
  - 默认值示例：`"{\"key\": \"value\"}"`
- `MAP_FORMAT_TYPE_STRING_TO_STRING` - 字符串键值对格式
  - 默认值示例：`"key1=value1,key2=value2"`
  - 使用逗号分隔多个键值对，每个键值对用等号连接
- `MAP_FORMAT_TYPE_STRING_TO_INT` - 字符串键整数值对格式
  - 默认值示例：`"key1=123,key2=456"`
  - 使用逗号分隔多个键值对，值必须是整数
  - **支持的整数类型**：`int32`, `sint32`, `sfixed32`, `int64`, `sint64`, `sfixed64`, `uint32`, `fixed32`, `uint64`, `fixed64`

#### 重复字段（repeated）

```protobuf
syntax = "proto3";
package tests;

import "flags/annotations.proto";

message Example {
  repeated string servers = 1 [(flags.value).repeated.string = {
    name: "servers"
    usage: "Server addresses (can be specified multiple times)"
    default: ["server1"]
  }];
}

```

### 嵌套消息配置

嵌套消息使用 `message` 标志类型：

```protobuf
message NestedConfig {
  string value = 1 [(flags.value).string = { name: "value" }];
}

message MainConfig {
  NestedConfig nested = 1 [(flags.value).message = {
    name: "nested"     // 嵌套消息的前缀名
    nested: true       // 启用嵌套标志生成
  }];
}
```

| 选项 | 类型 | 描述 |
|------|------|------|
| `name` | `string` | 嵌套消息的前缀名（默认为字段名） |
| `nested` | `bool` | 是否生成嵌套标志 |

## 分层标志组织

protoc-gen-go-flags 支持分层组织标志，通过 `WithPrefix` 和 `WithDelimiter` 选项实现。

### 基本前缀

```go
config.AddFlags(fs, flags.WithPrefix("server"))
```

生成：`--server.host`, `--server.port`

### 多级前缀

```go
config.AddFlags(fs, flags.WithPrefix("server", "database"))
```

生成：`--server.database.host`, `--server.database.port`

### 自定义分隔符

```go
config.AddFlags(fs,
  flags.WithPrefix("server"),
  flags.WithDelimiter("-"))  // 破折号
```

生成：`--server-host`, `--server-port`

支持的定界符：
- `flags.DelimiterDot` - 点号（默认）：`server.port`
- `flags.DelimiterDash` - 破折号：`server-port`
- `flags.DelimiterUnderscore` - 下划线：`server_port`
- `flags.DelimiterColon` - 冒号：`server:port`

### 自定义重命名函数

```go
config.AddFlags(fs,
  flags.WithPrefix("Server"),
  flags.WithRenamer(strings.ToLower))
```

生成：`--server-host`（转换为小写）

## 常见问题

### Q: 如何在现有项目中集成 protoc-gen-go-flags？

**A:** 按照以下步骤：
1. 安装插件：`go install github.com/kunstack/protoc-gen-go-flags@latest`
2. 复制 `annotations.proto` 到您的项目
3. 在 `.proto` 文件中添加标志注解
4. 运行 `protoc` 生成代码
5. 在应用中使用生成的 `AddFlags()` 方法

详细步骤请参阅[完整集成教程](#完整集成教程)。

### Q: 如何处理复杂的嵌套配置？

**A:** 使用嵌套消息和 `message` 标志类型：

```protobuf
syntax = "proto3";
package tests;

import "flags/annotations.proto";

message DatabaseConfig {
    string url = 1 [(flags.value).string = {
        name: "url"
        usage: "数据库连接 URL"
    }];
}

message AppConfig {
    DatabaseConfig database = 1 [(flags.value).message = {
        name: "db"
        nested: true
    }];
}
```

这将生成如 `--db-url` 这样的层级标志。

### Q: 如何自定义标志命名（使用前缀或分隔符）？

**A:** 在调用 `AddFlags` 时使用选项：

```go
// 使用前缀
config.AddFlags(fs, flags.WithPrefix("server"))
// 生成：--server.host

// 自定义分隔符
config.AddFlags(fs,
    flags.WithPrefix("server"),
    flags.WithDelimiter("-"))
// 生成：--server-host
```

### Q: 生成的代码报错 "undefined: flags.Option"

**A:** 您需要安装并导入运行时库：

```bash
go get github.com/kunstack/protoc-gen-go-flags/flags
```

```go
import "github.com/kunstack/protoc-gen-go-flags/flags"
```

### Q: 如何跳过特定字段的标志生成？

**A:** 只需不为该字段添加标志注解即可。如果已添加注解，可以设置字段级选项：

```protobuf
string internal_field = 1;  // 不添加标志注解，该字段不会生成标志
```

### Q: 如何设置字段的默认值？

**A:** 在标志注解中使用 `default` 选项：

```protobuf
int32 port = 1 [(flags.value).int32 = {
    name: "port"
    usage: "服务器端口"
    default: 8080  // 设置默认值
}];
```

然后在应用中调用 `config.SetDefaults()` 来应用默认值。

### Q: 支持哪些 protobuf 类型？

**A:** protoc-gen-go-flags 支持所有标准 protobuf 类型：
- 标量类型：string, int32, int64, bool, float, double 等
- 特殊类型：google.protobuf.Duration, Timestamp
- 复合类型：repeated（数组）、map（映射）
- 嵌套消息

详细列表请参阅[支持的类型](#支持的类型)部分。

### Q: 如何在标志中使用环境变量？

**A:** protoc-gen-go-flags 专注于命令行标志绑定。如需环境变量支持，建议结合使用 [viper](https://github.com/spf13/viper) 等配置管理库：

```go
import (
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
)

config := &proto.Config{}
fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
config.AddFlags(fs)

// 绑定到 viper（支持环境变量）
viper.BindPFlags(fs)
viper.AutomaticEnv()

fs.Parse(os.Args[1:])
```

### Q: 生成的文件命名规则是什么？

**A:** 对于 `.proto` 文件，会生成对应的 `.pb.flags.go` 文件：
- `config.proto` → `config.pb.go` + `config.pb.flags.go`
- `server.proto` → `server.pb.go` + `server.pb.flags.go`

### Q: 是否支持 gRPC？

**A:** protoc-gen-go-flags 与 gRPC 完全兼容。您可以在同一个 `.proto` 文件中同时定义 gRPC 服务和标志配置：

```bash
protoc \
    --go_out=. \
    --go-grpc_out=. \
    --flags_out=. \
    your_service.proto
```

## 贡献

欢迎贡献！如果您有建议或发现问题，请：

- 提交 Issue：[GitHub Issues](https://github.com/kunstack/protoc-gen-go-flags/issues)
- 提交 Pull Request：Fork 项目并创建 PR
- 改进文档：帮助完善文档和示例

## 许可证

本项目采用 Apache 2.0 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- [protoc-gen-star](https://github.com/lyft/protoc-gen-star) - 代码生成框架
- [spf13/pflag](https://github.com/spf13/pflag) - 命令行标志库
- [Google Protocol Buffers](https://protobuf.dev/) - 数据序列化协议