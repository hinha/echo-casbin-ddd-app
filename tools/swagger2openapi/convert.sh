#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"

# Default paths
SWAGGER_JSON="$PROJECT_ROOT/docs/swagger.json"
OUTPUT_FILE="$PROJECT_ROOT/echo-casbin-ddd-app.openapi.json"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--swagger)
            SWAGGER_JSON="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Run the conversion tool
cd "$SCRIPT_DIR"
go run main.go --swagger "$SWAGGER_JSON" --output "$OUTPUT_FILE"

# Check if conversion was successful
if [ $? -eq 0 ]; then
    echo "Conversion completed successfully!"
    echo "OpenAPI file created at: $OUTPUT_FILE"
else
    echo "Conversion failed. Please check the error messages above."
    exit 1
fi