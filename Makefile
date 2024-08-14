.PHONY: all build clean test deps build-linux

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=nvml-exporter
BINARY_UNIX=$(BINARY_NAME)_unix

# Directories
SRC_DIR=./cmd

# Default target executed when no arguments are given to make
all: build

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)/main.go

# Run tests
test:
	$(GOTEST) -v ./...

# Clean the build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Install dependencies
deps:
	$(GOGET) -u github.com/prometheus/client_golang/prometheus
	$(GOGET) -u github.com/NVIDIA/go-nvml

# Cross compilation for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(SRC_DIR)/main.go

