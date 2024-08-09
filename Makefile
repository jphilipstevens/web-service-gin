# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./bin/web-service-gin
MAIN_PATH=./main.go
VERSION_FILE=version.txt

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

build_with_new_version:
	./scripts/version.sh patch 
	$(GOBUILD) -ldflags "-X example/web-service-gin/app/version.Version=$(shell cat $(VERSION_FILE))" -o $(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./app/...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

deps:
	$(GOGET) -v -t -d ./...

seed:
	$(GOCMD) run ./scripts/seed.go

.PHONY: all build test coverage clean run deps seed
