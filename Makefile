.PHONY: all setup generate analyze populate-d1

# Default target
all: build test generate analyze

# Variables
GO_APP_NAME := halal-random-strings
GO_APP_PATH := ./cmd/$(GO_APP_NAME)
OUTPUT_FILE := output.txt
DB_NAME := my-halal-strings-db # Default D1 database name

# Setup target: Installs dependencies and sets up D1
setup:
	@echo "Setting up development environment..."
	@echo "1. Installing Go dependencies..."
	go mod tidy
	@echo "2. Ensuring Cloudflare Wrangler CLI is installed (npm install -g wrangler)..."
	@if ! command -v wrangler &> /dev/null; then \
		echo "Wrangler not found. Please install it globally: npm install -g wrangler"; \
		exit 1; \
	fi
	@echo "3. Authenticating Wrangler (wrangler login)..."
	@echo "Please follow the browser prompts to authenticate."
	wrangler login
	@echo "4. Creating D1 database '$(DB_NAME)' if it doesn't exist..."
	@echo "You might be prompted to confirm creation."
	wrangler d1 create $(DB_NAME) || true # '|| true' to prevent error if already exists
	@echo "5. Applying D1 schema..."
	wrangler d1 execute $(DB_NAME) --file db/schema.sql
	@echo "Setup complete. You can now run 'make generate', 'make analyze', etc."

# Build target: Compiles the Go application
build:
	@echo "Building Go application..."
	go build -o $(GO_APP_NAME) $(GO_APP_PATH)
	@echo "Build complete: ./$(GO_APP_NAME)"

# Test target: Runs Go tests
test:
	@echo "Running Go tests..."
	go test ./...

# Generate target: Runs the Go application to generate strings
generate:
	@echo "Generating strings to $(OUTPUT_FILE)..."
	@./$(GO_APP_NAME) -r 1000 > $(OUTPUT_FILE)
	@echo "Generated 1000 strings."

# Analyze target: Runs the Python script to analyze uniqueness
analyze:
	@echo "Analyzing uniqueness from $(OUTPUT_FILE)..."
	python3 scripts/analyze.py $(OUTPUT_FILE)

# Populate D1 target: Runs the Python script to feed output.txt to D1
populate-d1:
	@echo "Populating D1 database '$(DB_NAME)' from $(OUTPUT_FILE)..."
	python3 scripts/populate_d1.py $(OUTPUT_FILE) $(DB_NAME)

# Clean target: Removes generated files
clean:
	@echo "Cleaning up generated files..."
	rm -f $(GO_APP_NAME) $(OUTPUT_FILE)
	@echo "Cleanup complete."
