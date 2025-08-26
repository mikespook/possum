#!/bin/bash

# Run all tests with verbose output
go test -v ./...

# Print test coverage
echo -e "\nTest coverage:"
go test -cover ./...