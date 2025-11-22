#!/bin/bash

# Run Go test coverage
echo "Running Go test coverage..."

# Run tests with coverage and generate coverage profile
go tool cover -html=coverage.out

# Check if coverage profile was generated
if [ -f coverage.out ]; then
    echo "Coverage profile generated successfully"
    
    # Display coverage summary
    echo "Coverage summary:"
    go tool cover -func=coverage.out
    
    # Generate HTML coverage report
    echo "Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    
    echo "Coverage report saved to coverage.html"
    echo "Open coverage.html in your browser to view detailed coverage"
else
    echo "Failed to generate coverage profile"
    exit 1
fi
