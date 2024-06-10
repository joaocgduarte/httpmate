BIN_DIR := ./bin
APP_NAME := httpmate

# Default target
all: build

# Build the Go application
build:
	go build -o $(BIN_DIR)/$(APP_NAME) .

# Clean the build
clean:
	rm -f $(BIN_DIR)/$(APP_NAME)

# Phony targets
.PHONY: all build clean
