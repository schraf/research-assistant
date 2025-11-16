# Variables
BINARY_NAME=gemini-email
BUILD_DIR=bin
MAIN_PATH=./cmd
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Default target
.PHONY: all
all: vet build

# Build the application
.PHONY: build
build: vet test
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@if [ -f .env ]; then \
		echo "Loading environment from .env file..."; \
		export $$(cat .env | xargs) && go run $(MAIN_PATH); \
	else \
		echo "No .env file found, running without environment variables..."; \
		go run $(MAIN_PATH); \
	fi

# Vet the code
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...
	@echo "Vet complete"

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf logs
	rm -f .env
	@echo "Clean complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Terraform targets
.PHONY: terraform-init
terraform-init:
	@echo "Initializing Terraform..."
	cd terraform && terraform init
	@echo "Terraform initialized"

.PHONY: terraform-plan
terraform-plan:
	@echo "Planning Terraform changes..."
	cd terraform && terraform plan
	@echo "Terraform plan complete"

.PHONY: terraform-apply
terraform-apply:
	@echo "Applying Terraform changes..."
	cd terraform && terraform apply
	@echo "Terraform apply complete"

.PHONY: terraform-destroy
terraform-destroy:
	@echo "Destroying Terraform resources..."
	cd terraform && terraform destroy
	@echo "Terraform destroy complete"

.PHONY: terraform-validate
terraform-validate:
	@echo "Validating Terraform configuration..."
	cd terraform && terraform validate
	@echo "Terraform validation complete"

.PHONY: terraform-fmt
terraform-fmt:
	@echo "Formatting Terraform files..."
	cd terraform && terraform fmt -recursive
	@echo "Terraform formatting complete"

.PHONY: terraform-output
terraform-output:
	@echo "Showing Terraform outputs..."
	cd terraform && terraform output
	@echo "Terraform outputs complete"

# Setup infrastructure and run application
.PHONY: setup-infra
setup-infra: terraform-init terraform-apply
	@echo "Infrastructure setup complete"

# Clean everything including Terraform state
.PHONY: clean-all
clean-all: clean terraform-destroy
	@echo "Cleaning all artifacts and infrastructure..."

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Run vet and build"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  vet          - Vet the code"
	@echo "  fmt          - Format the code"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo ""
	@echo "Terraform targets:"
	@echo "  terraform-init      - Initialize Terraform"
	@echo "  terraform-plan      - Plan Terraform changes"
	@echo "  terraform-apply     - Apply Terraform changes (creates .env file)"
	@echo "  terraform-destroy   - Destroy Terraform resources"
	@echo "  terraform-validate  - Validate Terraform configuration"
	@echo "  terraform-fmt       - Format Terraform files"
	@echo "  terraform-output    - Show Terraform outputs"
	@echo "  setup-infra         - Setup infrastructure (init + apply)"
	@echo "  clean-all           - Clean everything including infrastructure"
	@echo ""
	@echo "  help         - Show this help message"
