PROJECT_NAME=wbfbd

default: build

all: fmt vet build

build:
	go build -v -o $(PROJECT_NAME) .

help:
	@echo "Build related targets:"
	@echo "  all      Run fmt vet build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality related targets:"
	@echo "  fmt  Format files with go fmt."
	@echo "  vet  Examine go sources with go vet."
	@echo "Test related targets:"
	@echo "  test                Run all tests."
	@echo "Other targets:"
	@echo "  help  Print help message."

fmt:
	go fmt ./...

vet:
	go vet ./...

install:
	cp $(PROJECT_NAME) /usr/bin

uninstall:
	rm /usr/bin/$(PROJECT_NAME)
