
# Simple Makefile for a Go project

# Build the application
all: build

build:
	@cp active.* cmd/api/
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@cp active.* cmd/api/
	@go run cmd/api/main.go

# Test the application
test:
	@cp active.* cmd/api/
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

up:
	@ migrate -source file://DB/Migrations/ -database "$DATABASE_URL" up

down:
	@ migrate -source file://DB/Migrations/ -database "$DATABASE_URL" down

# Extract all goi18n message structs
extract:
	@goi18n extract

# Merge goi18n files for translation
merge:
	@goi18n merge active.en.toml active.ar.toml

translate: extract merge

.PHONY: all build run test clean
		
