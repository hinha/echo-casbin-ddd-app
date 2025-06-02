package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func main() {
	// Define command line flags
	swaggerFile := flag.String("swagger", "docs/swagger.json", "Path to Swagger JSON file")
	outputFile := flag.String("output", "openapi.json", "Path to output OpenAPI JSON file")
	flag.Parse()

	// Read Swagger JSON file
	swaggerData, err := os.ReadFile(*swaggerFile)
	if err != nil {
		log.Fatalf("Failed to read Swagger file: %v", err)
	}

	// Parse Swagger JSON
	var swagger openapi2.T
	if err := json.Unmarshal(swaggerData, &swagger); err != nil {
		log.Fatalf("Failed to parse Swagger JSON: %v", err)
	}

	// Convert Swagger to OpenAPI
	openapi, err := openapi2conv.ToV3(&swagger)
	if err != nil {
		log.Fatalf("Failed to convert Swagger to OpenAPI: %v", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Write OpenAPI JSON to file
	openapiData, err := json.MarshalIndent(openapi, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI JSON: %v", err)
	}

	if err := os.WriteFile(*outputFile, openapiData, 0644); err != nil {
		log.Fatalf("Failed to write OpenAPI file: %v", err)
	}

	fmt.Printf("Successfully converted Swagger to OpenAPI: %s\n", *outputFile)
}
