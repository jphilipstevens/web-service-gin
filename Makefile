# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./bin/web-service-gin
MAIN_PATH=./example/main.go
VERSION_FILE=version.txt

all: build ## Default: Ensure library compiles (safe check)

# 🛠️ Build
build: ## Build the entire library (no binary output)
	go build ./...

build-example: ## Build the example binary
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

build_with_new_version: ## Patch version + build binary with embedded version string
	./scripts/version.sh patch 
	$(GOBUILD) -ldflags "-X github.com/jphilipstevens/web-service-gin/app/version.Version=$(shell cat $(VERSION_FILE))" -o $(BINARY_NAME) -v $(MAIN_PATH)

# 🧪 Test & Coverage
test: ## Run unit tests
	$(GOTEST) -v ./app/...

coverage: ## Run tests with coverage output in HTML
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 🧹 Clean up
clean: ## Remove binary and coverage output files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

# 🚀 Run
run: ## Build and run the example binary
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# 📦 Dependency management
deps: ## Download Go module dependencies
	$(GOGET) -v -t -d ./...

# 🌱 Data Seeding
seed: ## Run the Go-based seeder script
	$(GOCMD) run ./scripts/seed.go

# 📘 Help Menu
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "🛠  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: all build build-example build_with_new_version \
        test coverage clean run deps seed help
