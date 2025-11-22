# Variables
BUILD_DIR=bin
REGION ?= us-central1

# Default target
.PHONY: all
all: vet build

# Build the application
.PHONY: build
build: vet test
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/gentelegraphtoken ./cmd/gentelegraphtoken
	go build -o $(BUILD_DIR)/server ./cmd/server
	go build -o $(BUILD_DIR)/worker ./cmd/worker

# Run the application
.PHONY: run
run:
	@echo "Running..."
	go run ./cmd/server

# Run tool to generate a telegraph token
.PHONY: telegraph-token
telegraph-token:
	@echo "Generating telegraph token..."
	go run ./cmd/gentelegraphtoken


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

# Container and deployment targets
.PHONY: container-build
container-build:
	@echo "Building container image with podman..."
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "Error: PROJECT_ID environment variable is required"; \
		echo "Usage: PROJECT_ID=your-project-id make container-build"; \
		exit 1; \
	fi
	@podman build -t $(REGION)-docker.pkg.dev/$(PROJECT_ID)/research-assistant/research-assistant:latest .

.PHONY: container-push
container-push: container-build
	@echo "Pushing container image to Artifact Registry..."
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "Error: PROJECT_ID environment variable is required"; \
		echo "Usage: PROJECT_ID=your-project-id make container-push"; \
		exit 1; \
	fi
	@podman push $(REGION)-docker.pkg.dev/$(PROJECT_ID)/research-assistant/research-assistant:latest

.PHONY: gcloud-build
gcloud-build:
	@echo "Building and deploying with Cloud Build..."
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "Error: PROJECT_ID environment variable is required"; \
		echo "Usage: PROJECT_ID=your-project-id make gcloud-build"; \
		exit 1; \
	fi
	@TAG=$$(git rev-parse --short HEAD 2>/dev/null); \
	if [ -z "$$TAG" ]; then \
		TAG=$$(date +%s); \
	fi; \
	echo "Using tag: $$TAG"; \
	gcloud builds submit --config=cloudbuild.yaml \
		--substitutions=_REGION=$(REGION),_REPO_NAME=research-assistant,_TAG=$$TAG \
		--project=$(PROJECT_ID)

.PHONY: deploy
deploy: gcloud-build
	@echo "Deployment complete!"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all             - Run vet and build"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  auth-token      - Generate an API auth token"
	@echo "  telegraph-token - Generate a Telegraph API token"
	@echo "  vet             - Vet the code"
	@echo "  fmt             - Format the code"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo ""
	@echo "Terraform targets:"
	@echo "  terraform-init      - Initialize Terraform"
	@echo "  terraform-plan      - Plan Terraform changes"
	@echo "  terraform-apply     - Apply Terraform changes"
	@echo "  terraform-destroy   - Destroy Terraform resources"
	@echo "  terraform-validate  - Validate Terraform configuration"
	@echo "  terraform-fmt       - Format Terraform files"
	@echo "  terraform-output    - Show Terraform outputs"
	@echo "  setup-infra         - Setup infrastructure (init + apply)"
	@echo "  clean-all           - Clean everything including infrastructure"
	@echo ""
	@echo "Deployment targets:"
	@echo "  container-build     - Build container image locally with podman (requires PROJECT_ID)"
	@echo "  container-push      - Build and push container image with podman (requires PROJECT_ID)"
	@echo "  gcloud-build        - Build and deploy using Cloud Build (requires PROJECT_ID)"
	@echo "  deploy              - Full deployment using Cloud Build (requires PROJECT_ID)"
	@echo ""
	@echo "  help               - Show this help message"
