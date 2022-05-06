PROJECT_NAME=vfbdb

.PHONY: default all build help fmt vet install uninstall

default: build

help:
	@echo "Build targets:"
	@echo "  all      Run fmt vet build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality targets:"
	@echo "  fmt  Format files with go fmt."
	@echo "  vet  Examine go sources with go vet."
	@echo "Test targets:"
	@echo "  test  Run go test."
	@echo "Other targets:"
	@echo "  help  Print help message."

all: fmt vet build

build:
	go build -v -o $(PROJECT_NAME) ./cmd/vfbdb

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...
	
install:
	cp $(PROJECT_NAME) /usr/bin

uninstall:
	rm /usr/bin/$(PROJECT_NAME)
