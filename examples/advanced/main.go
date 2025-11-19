package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/kunstack/protoc-gen-go-flags/examples/advanced/proto"
)

func main() {
	// Create configuration instance
	config := &proto.AdvancedConfig{}

	// Set default values
	config.SetDefaults()

	// Create flag set
	fs := pflag.NewFlagSet("advanced", pflag.ExitOnError)

	// Add flags from protobuf configuration
	config.AddFlags(fs)

	// Parse command-line arguments
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Display configuration
	fmt.Println("Advanced Configuration Example")
	fmt.Println("==============================")
	fmt.Println()

	fmt.Println("Scalar Types:")
	fmt.Printf("  Name: %s\n", config.Name)
	fmt.Printf("  Workers: %d\n", config.WorkerCount)
	fmt.Printf("  Max Memory: %d bytes\n", config.MaxMemory)
	fmt.Printf("  Rate Limit: %d req/s\n", config.RateLimit)
	fmt.Printf("  Ratio: %.2f\n", config.Ratio)
	fmt.Printf("  Enabled: %v\n", config.Enabled)
	fmt.Printf("  API Key: %s\n", string(config.ApiKey))
	fmt.Println()

	fmt.Println("Special Types:")
	if config.Timeout != nil {
		fmt.Printf("  Timeout: %v\n", config.Timeout.AsDuration())
	}
	if config.CreatedAt != nil {
		fmt.Printf("  Created At: %v\n", config.CreatedAt.AsTime())
	}
	fmt.Println()

	fmt.Println("Repeated Fields:")
	fmt.Printf("  Servers: %s\n", strings.Join(config.Servers, ", "))
	fmt.Printf("  Ports: %v\n", config.Ports)
	fmt.Println()

	fmt.Println("Map Fields:")
	fmt.Printf("  Labels: %v\n", config.Labels)
	fmt.Printf("  Limits: %v\n", config.Limits)
	fmt.Println()

	fmt.Println("Enum:")
	fmt.Printf("  Log Level: %s (%d)\n", config.LogLevel.String(), config.LogLevel)
	fmt.Println()

	fmt.Println("Advanced Features:")
	fmt.Printf("  Deprecated Flag: %v\n", config.DeprecatedFlag)
	fmt.Printf("  Internal Token: %s\n", config.InternalToken)
	fmt.Println()

	fmt.Println("Configuration loaded successfully!")
	fmt.Println("\nTry running with different flags:")
	fmt.Println("  ./bin/advanced --name MyApp --workers 8 --timeout 60s")
	fmt.Println("  ./bin/advanced --servers api1.com,api2.com --labels env=prod,region=us")
	fmt.Println("  ./bin/advanced --log-level LOG_LEVEL_INFO --limits cpu=2000,memory=4096")
}
