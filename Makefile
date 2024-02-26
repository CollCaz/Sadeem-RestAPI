
# Simple Makefile for a Go project

# Build the application
all: build

build:
	@cp active.* cmd/api/
	@go build -o bin/sadeemAPI cmd/api/main.go

# Run the application
run:
	@cp active.* cmd/api/
	@go run cmd/api/main.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -r -f bin

up:
	@migrate -source file://DB/Migrations/ -database "$$DATABASE_URL" up

down:
	@migrate -source file://DB/Migrations/ -database "$$DATABASE_URL" down

help:
	@printf "build:builds the app in ./bin/\nrun: runs the application\nclean: removes the bin directory\nup: applies up migrations\ndown: applies down migrations\n"

.PHONY: all build run test clean
		
