PROJECT_NAME=vfbdb

.PHONY: default all build help fmt vet install uninstall

default: build

help:
	@echo "Build targets:"
	@echo "  all      Run fmt vet build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality targets:"
	@echo "  fmt   Format files with go fmt."
	@echo "  lint  Lint files with golangci-lint."
	@echo "Test targets:"
	@echo "  test  Run go test."
	@echo "Other targets:"
	@echo "  help  Print help message."

# Build targets
all: lint fmt build

build:
	go build -v -o $(PROJECT_NAME) ./cmd/vfbdb


# Quality targets
fmt:
	go fmt ./...

lint:
	golangci-lint run


# Test targets
test:
	go test ./...


# Installation targets
install:
	cp $(PROJECT_NAME) /usr/bin

uninstall:
	rm /usr/bin/$(PROJECT_NAME)
