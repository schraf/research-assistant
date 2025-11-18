# Research Assistant

An AI-powered research assistant that automatically conducts research on topics, generates structured reports, publishes them to Telegraph, and sends email notifications. Built with Go and deployed on Google Cloud Platform.

## Overview

The Research Assistant is a cloud-native application that leverages Google Gemini AI to perform comprehensive research on any given topic. It uses a serverless architecture with Google Cloud Run for both the API server and asynchronous job processing.

### Features

- **AI-Powered Research**: Uses Google Gemini AI to conduct in-depth research on topics
- **Structured Reports**: Generates well-formatted research reports with sections and paragraphs
- **Telegraph Publishing**: Automatically publishes research reports to Telegraph
- **Email Notifications**: Sends email notifications when research is complete
- **Asynchronous Processing**: Research jobs are processed asynchronously via Cloud Run Jobs
- **Token-Based Authentication**: Secure API access using bearer token authentication
- **Infrastructure as Code**: Terraform configuration for easy deployment

## Architecture

The application consists of two main components:

1. **Server** (`cmd/server`): HTTP API service that accepts research requests
   - Runs as a Cloud Run Service
   - Validates authentication tokens
   - Queues research jobs to Cloud Run Jobs

2. **Worker** (`cmd/worker`): Background job processor that performs the research
   - Runs as a Cloud Run Job
   - Processes research requests asynchronously
   - Generates research reports using Gemini AI
   - Publishes reports to Telegraph
   - Sends email notifications

### Research Process

1. **Planning**: Generates a research plan with subtopics and questions
2. **Knowledge Gathering**: Iteratively gathers knowledge on each subtopic
3. **Analysis**: Analyzes collected knowledge to identify gaps
4. **Synthesis**: Synthesizes all research into a structured report
5. **Publishing**: Posts the report to Telegraph
6. **Notification**: Sends an email with the report URL

## Prerequisites

- Go 1.24 or later
- Google Cloud Platform account with billing enabled
- `gcloud` CLI installed and authenticated
- Terraform (for infrastructure provisioning)
- Docker or Podman (for local container builds, optional)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/schraf/research-assistant.git
cd research-assistant
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Configure Environment Variables

Create a `.env` file in the project root with the following variables:

```bash
# Google Cloud
GOOGLE_CLOUD_PROJECT=your-project-id
GOOGLE_API_KEY=your-gemini-api-key

# Telegraph
TELEGRAPH_API_KEY=your-telegraph-api-key
TELEGRAPH_AUTHOR_NAME=Your Name

# Authentication
AUTH_SECRET=your-secret-key
AUTH_TOKEN_MESSAGES=message1,message2

# Email (SMTP)
MAIL_SMTP_SERVER=smtp.gmail.com
MAIL_SMTP_PORT=587
MAIL_SENDER_EMAIL=your-email@gmail.com
MAIL_SENDER_PASSWORD=your-app-password
MAIL_RECIPIENT_EMAIL=recipient@example.com

# Cloud Run (for production)
CLOUD_RUN_JOB_NAME=research-worker
CLOUD_RUN_JOB_REGION=us-central1
```

### 4. Generate Auth Token

Before making API requests, generate an authentication token:

```bash
make auth-token
```

Or manually:

```bash
export AUTH_SECRET="your-secret"
export AUTH_TOKEN_MESSAGES="message1,message2"
./bin/genauthtoken -seed "$AUTH_SECRET"
```

Use one of the generated tokens in the `Authorization: Bearer` header.

### 5. Generate Telegraph Token (if needed)

If you need to create a new Telegraph account:

```bash
make telegraph-token
```

### 6. Build and Run Locally

```bash
# Build all binaries
make build

# Run the server
make run
```

The server will start on `http://localhost:8080` (or the port specified by `PORT` environment variable).

### 7. Make a Research Request

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/research?topic=artificial%20intelligence"
```

## Project Structure

```
.
├── cmd/
│   ├── server/              # HTTP API server
│   ├── worker/              # Background job processor
│   ├── genauthtoken/        # Auth token generator tool
│   └── gentelegraphtoken/   # Telegraph token generator tool
├── internal/
│   ├── auth/               # Authentication and token validation
│   ├── gemini/             # Google Gemini AI client
│   ├── mail/               # Email sending functionality
│   ├── models/             # Data models and interfaces
│   ├── researcher/         # Core research logic
│   ├── service/            # HTTP service handlers
│   ├── telegraph/          # Telegraph API client
│   ├── utils/              # Utility functions
│   └── worker/             # Worker job processing
├── terraform/              # Infrastructure as Code
├── bin/                    # Build output directory
├── logs/                   # Log files
├── Dockerfile              # Container image definition
├── cloudbuild.yaml         # Cloud Build configuration
├── Makefile                # Build and deployment commands
└── go.mod                  # Go module dependencies
```

## Development

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Code Linting

```bash
make vet
```

### Building Binaries

```bash
make build
```

This builds all binaries to the `bin/` directory:
- `bin/server` - HTTP API server
- `bin/worker` - Background job processor
- `bin/genauthtoken` - Auth token generator
- `bin/gentelegraphtoken` - Telegraph token generator

## Deployment

The application is designed to run on Google Cloud Platform. See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions.

### Quick Deployment

1. **Configure Terraform variables**:
   ```bash
   cp terraform/terraform.tfvars.example terraform/terraform.tfvars
   # Edit terraform/terraform.tfvars with your values
   ```

2. **Provision infrastructure**:
   ```bash
   make setup-infra
   ```

3. **Deploy application**:
   ```bash
   PROJECT_ID=your-project-id make deploy
   ```

This will:
- Build the container image
- Push it to Artifact Registry
- Deploy to Cloud Run Service and Cloud Run Job

## API Reference

### POST /research

Initiates a research job for the given topic.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Query Parameters:**
- `topic` (required): The research topic

**Response:**
```json
{
  "success": true,
  "request_id": "uuid",
  "message": "Research request queued"
}
```

**Example:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "https://your-service-url/research?topic=quantum%20computing"
```

## Environment Variables

### Required for Server

- `PORT`: Server port (default: 8080)
- `AUTH_SECRET`: Secret key for token generation
- `AUTH_TOKEN_MESSAGES`: Comma-separated list of valid auth messages
- `CLOUD_RUN_JOB_NAME`: Name of the Cloud Run Job
- `CLOUD_RUN_JOB_REGION`: Region of the Cloud Run Job
- `GOOGLE_CLOUD_PROJECT`: GCP project ID

### Required for Worker

- `GOOGLE_API_KEY`: Google Gemini API key
- `TELEGRAPH_API_KEY`: Telegraph API access token
- `TELEGRAPH_AUTHOR_NAME`: Author name for Telegraph articles
- `MAIL_SMTP_SERVER`: SMTP server hostname
- `MAIL_SMTP_PORT`: SMTP server port
- `MAIL_SENDER_EMAIL`: Email address for sending notifications
- `MAIL_SENDER_PASSWORD`: Email password (use app-specific password for Gmail)
- `MAIL_RECIPIENT_EMAIL`: Email address to receive notifications

## Makefile Commands

### Development Commands

- `make all` - Run vet and build
- `make build` - Build all binaries
- `make run` - Run the server locally
- `make test` - Run tests
- `make fmt` - Format code
- `make vet` - Vet code
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies

### Utility Commands

- `make auth-token` - Generate authentication tokens
- `make telegraph-token` - Generate Telegraph API token

### Terraform Commands

- `make terraform-init` - Initialize Terraform
- `make terraform-plan` - Plan Terraform changes
- `make terraform-apply` - Apply Terraform changes
- `make terraform-destroy` - Destroy infrastructure
- `make terraform-validate` - Validate Terraform configuration
- `make terraform-fmt` - Format Terraform files
- `make terraform-output` - Show Terraform outputs
- `make setup-infra` - Setup infrastructure (init + apply)

### Deployment Commands

- `make container-build` - Build container image locally (requires PROJECT_ID)
- `make container-push` - Build and push container image (requires PROJECT_ID)
- `make gcloud-build` - Build and deploy using Cloud Build (requires PROJECT_ID)
- `make deploy` - Full deployment using Cloud Build (requires PROJECT_ID)

Run `make help` to see all available commands.

## Dependencies

- [Google Gemini AI](https://ai.google.dev/) - AI research capabilities
- [Telegraph API](https://telegra.ph/api) - Publishing platform
- [Google Cloud Run](https://cloud.google.com/run) - Serverless compute
- [Terraform](https://www.terraform.io/) - Infrastructure provisioning

## License

See [LICENSE](LICENSE).

## Support

For issues and questions, please open an issue on GitHub.
