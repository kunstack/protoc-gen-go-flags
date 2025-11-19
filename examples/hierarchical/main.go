package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/kunstack/protoc-gen-go-flags/examples/hierarchical/proto"
	"github.com/kunstack/protoc-gen-go-flags/flags"
)

func main() {
	fmt.Println("Hierarchical Flags Organization Demo")
	fmt.Println("====================================")
	fmt.Println()

	// Demonstrate different delimiter styles
	demonstrateStyles()

	fmt.Println()
	fmt.Println("Multi-Level Prefix Example")
	fmt.Println("--------------------------")
	demonstrateMultiLevel()

	fmt.Println()
	fmt.Println("Practical Use Case: Multiple Services")
	fmt.Println("-------------------------------------")
	demonstrateMultipleServices()
}

func demonstrateStyles() {
	config := &proto.ServiceConfig{}
	config.SetDefaults()

	styles := []struct {
		name    string
		options []flags.Option
	}{
		{
			name:    "No Prefix (Default)",
			options: nil,
		},
		{
			name:    "Dot Delimiter (Default)",
			options: []flags.Option{flags.WithPrefix("service")},
		},
		{
			name: "Dash Delimiter",
			options: []flags.Option{
				flags.WithPrefix("service"),
				flags.WithDelimiter("-"),
			},
		},
		{
			name: "Underscore Delimiter",
			options: []flags.Option{
				flags.WithPrefix("service"),
				flags.WithDelimiter("_"),
			},
		},
		{
			name: "Colon Delimiter",
			options: []flags.Option{
				flags.WithPrefix("service"),
				flags.WithDelimiter(":"),
			},
		},
	}

	for i, style := range styles {
		fmt.Printf("%d. %s:\n", i+1, style.name)

		fs := pflag.NewFlagSet(style.name, pflag.ContinueOnError)
		config.AddFlags(fs, style.options...)

		// Show generated flag names
		fs.VisitAll(func(f *pflag.Flag) {
			fmt.Printf("   --%s\n", f.Name)
		})
		fmt.Println()
	}
}

func demonstrateMultiLevel() {
	config := &proto.ServiceConfig{}
	config.SetDefaults()

	fs := pflag.NewFlagSet("multi-level", pflag.ContinueOnError)
	config.AddFlags(fs,
		flags.WithPrefix("myapp", "backend", "api"),
		flags.WithDelimiter("."))

	fmt.Println("Prefix: myapp.backend.api")
	fmt.Println("Generated flags:")
	fs.VisitAll(func(f *pflag.Flag) {
		fmt.Printf("  --%s\n", f.Name)
	})
}

func demonstrateMultipleServices() {
	// API service
	apiConfig := &proto.ServiceConfig{}
	apiConfig.SetDefaults()
	apiConfig.Host = "api.example.com"
	apiConfig.Port = 8080

	// Database service
	dbConfig := &proto.ServiceConfig{}
	dbConfig.SetDefaults()
	dbConfig.Host = "db.example.com"
	dbConfig.Port = 5432

	// Cache service
	cacheConfig := &proto.ServiceConfig{}
	cacheConfig.SetDefaults()
	cacheConfig.Host = "cache.example.com"
	cacheConfig.Port = 6379

	// Create flag set with all services
	fs := pflag.NewFlagSet("services", pflag.ContinueOnError)
	apiConfig.AddFlags(fs, flags.WithPrefix("api"), flags.WithDelimiter("-"))
	dbConfig.AddFlags(fs, flags.WithPrefix("db"), flags.WithDelimiter("-"))
	cacheConfig.AddFlags(fs, flags.WithPrefix("cache"), flags.WithDelimiter("-"))

	// Parse arguments
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Display configuration
	fmt.Println("\nAPI Service:")
	fmt.Printf("  Host: %s\n", apiConfig.Host)
	fmt.Printf("  Port: %d\n", apiConfig.Port)
	fmt.Printf("  Timeout: %v\n", apiConfig.Timeout.AsDuration())

	fmt.Println("\nDatabase Service:")
	fmt.Printf("  Host: %s\n", dbConfig.Host)
	fmt.Printf("  Port: %d\n", dbConfig.Port)
	fmt.Printf("  Max Connections: %d\n", dbConfig.MaxConnections)

	fmt.Println("\nCache Service:")
	fmt.Printf("  Host: %s\n", cacheConfig.Host)
	fmt.Printf("  Port: %d\n", cacheConfig.Port)
	fmt.Printf("  TLS: %v\n", cacheConfig.TlsEnabled)

	fmt.Println("\nExample usage:")
	fmt.Println("  ./bin/hierarchical --api-host api.prod.com --api-port 443 --api-tls")
	fmt.Println("  ./bin/hierarchical --db-host db.prod.com --db-max-connections 200")
	fmt.Println("  ./bin/hierarchical --cache-host redis.prod.com --cache-port 6380")
}
