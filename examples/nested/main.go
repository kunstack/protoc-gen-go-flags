package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	proto "github.com/kunstack/protoc-gen-go-flags/examples/nested/proto"
)

func main() {
	// Create configuration instance
	config := &proto.AppConfig{}

	// Set default values
	config.SetDefaults()

	// Create flag set
	fs := pflag.NewFlagSet("app", pflag.ExitOnError)

	// Add flags from protobuf configuration
	config.AddFlags(fs)

	// Parse command-line arguments
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Display configuration
	fmt.Println("Application Configuration:")
	fmt.Println("=========================")
	fmt.Println()

	fmt.Println("Server:")
	fmt.Printf("  Host: %s\n", config.Server.Host)
	fmt.Printf("  Port: %d\n", config.Server.Port)
	fmt.Printf("  HTTPS: %v\n", config.Server.Https)
	fmt.Println()

	fmt.Println("Database:")
	fmt.Printf("  URL: %s\n", config.Database.Url)
	fmt.Printf("  Max Connections: %d\n", config.Database.MaxConnections)
	fmt.Println()

	fmt.Println("Logging:")
	fmt.Printf("  Level: %s\n", config.Logging.Level)
	fmt.Printf("  Format: %s\n", config.Logging.Format)
	fmt.Println()

	fmt.Println("Configuration loaded successfully!")
}
