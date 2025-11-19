package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	proto "github.com/kunstack/protoc-gen-go-flags/examples/basic/proto"
)

func main() {
	// Create configuration instance
	config := &proto.ServerConfig{}

	// Set default values (optional but recommended)
	config.SetDefaults()

	// Create flag set
	fs := pflag.NewFlagSet("server", pflag.ExitOnError)

	// Add flags from protobuf configuration
	config.AddFlags(fs)

	// Parse command-line arguments
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Display configuration
	fmt.Println("Starting server...")
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Debug: %v\n", config.Debug)

	// Your application logic would go here
	fmt.Println("\nServer configuration loaded successfully!")
	fmt.Println("(This is a demo - no actual server is started)")
}
