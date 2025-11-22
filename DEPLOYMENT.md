# Deployment Guide

This guide explains how to build and deploy the Research Assistant application to Google Cloud Platform. The deployment includes:
- Google Cloud Run Service (API server)
- Google Cloud Run Job (background worker)
- Google Cloud API Gateway (public API endpoint)

## Prerequisites

1. **Google Cloud Project**: You need a GCP project with billing enabled
2. **gcloud CLI**: Install and authenticate with `gcloud auth login`
3. **Docker or Podman**: For local container builds (optional, only needed for Option B deployment)
4. **Terraform**: For infrastructure provisioning

## Step 1: Configure Terraform Variables

1. Copy the example terraform variables file:
   ```bash
   cp terraform/terraform.tfvars.example terraform/terraform.tfvars
   ```

2. Edit `terraform/terraform.tfvars` and fill in all required values:
   - `project_id`: Your GCP project ID
   - `region`: GCP region (defaults to us-central1)
   - `google_api_key`: Your Google Gemini AI API key
   - `telegraph_api_key`: Your Telegraph API key
   - `telegraph_author_name`: Author name for Telegraph articles
   - `auth_secret`: Secret key for token generation
   - `auth_token_messages`: Comma-separated list of valid auth messages
   - `smtp_hostname`: SMTP server hostname (default: smtp.gmail.com)
   - `smtp_port`: SMTP server port (default: 587)
   - `mail_sender_email`: Sender email address
   - `mail_sender_password`: Sender email password (use app-specific password for Gmail)
   - `mail_recipient_email`: Recipient email address

## Step 2: Provision Infrastructure

1. Initialize Terraform:
   ```bash
   make terraform-init
   ```

2. Review the planned changes:
   ```bash
   make terraform-plan
   ```

3. Apply the infrastructure:
   ```bash
   make terraform-apply
   ```

This will create:
- Artifact Registry repository for container images
- Cloud Run service (initially with a placeholder "hello world" image)
- Cloud Run Job for processing research requests
- API Gateway API, API Config, and Gateway
- Service accounts with necessary permissions
- Required GCP APIs enabled (including API Gateway APIs)

**Note:** The Cloud Run service will be created with a placeholder image (`gcr.io/cloudrun/hello`) since the actual application image doesn't exist yet. The first deployment in Step 3 will replace this placeholder with your actual application image.

**Important:** The Cloud Run service is not directly accessible from the public internet. All traffic must go through the API Gateway, which is created automatically by Terraform.

## Step 3: Build and Deploy Container Image

You have two options for building and deploying:

### Option A: Using Cloud Build (Recommended)

This is the easiest method as it handles authentication automatically. This will build your application image and replace the placeholder image in the Cloud Run service:

```bash
PROJECT_ID=your-project-id make deploy
```

Or with a custom region:
```bash
PROJECT_ID=your-project-id REGION=us-central1 make deploy
```

**Note:** This step replaces the placeholder image created by Terraform with your actual application image.

### Option B: Local Container Build and Push

1. Authenticate with Artifact Registry (works for both Docker and Podman):
   ```bash
   gcloud auth configure-docker us-central1-docker.pkg.dev
   ```

2. Build and push the image:
   ```bash
   PROJECT_ID=your-project-id make container-push
   ```

3. Update the Cloud Run service to use the new image (replace `us-central1` with your region if different):
   ```bash
   gcloud run services update research-assistant \
     --image us-central1-docker.pkg.dev/your-project-id/research-assistant/research-assistant:latest \
     --region us-central1
   ```

4. Update the Cloud Run Job to use the new image (replace `us-central1` with your region if different):
   ```bash
   gcloud run jobs update research-worker \
     --image us-central1-docker.pkg.dev/your-project-id/research-assistant/research-assistant:latest \
     --region us-central1
   ```

## Step 4: Verify Deployment

1. Get the API Gateway URL (this is the public endpoint you should use):
   ```bash
   make terraform-output
   ```

   Or directly:
   ```bash
   cd terraform && terraform output api_gateway_url
   ```

   **Note:** The Cloud Run service URL is also available via `terraform output cloud_run_service_url`, but it's not directly accessible from the internet. All production traffic should use the API Gateway URL.

2. Test the service via API Gateway:
   ```bash
   # Generate an auth token first (see below)
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "https://YOUR-GATEWAY-URL/research?topic=test"
   ```

   Replace `YOUR-GATEWAY-URL` with the URL from step 1 (e.g., `https://research-assistant-gateway-xxxxx-uc.a.run.app`).

## Generating Auth Tokens

Before you can call the service, you need to generate an auth token:

```bash
export AUTH_SECRET="your-secret-from-tfvars"
export AUTH_TOKEN_MESSAGES="message1,message2"
./bin/genauthtoken -seed "$AUTH_SECRET"
```

This will output tokens for each message. Use one of these tokens in the `Authorization: Bearer` header.

## Environment Variables in Cloud Run

All environment variables are automatically set by Terraform from your `terraform.tfvars` file. The following variables are configured:

### Cloud Run Service Environment Variables

- `PORT`: Defaults to 8080 (Cloud Run sets this automatically, but the service defaults to 8080 if not set)
- `GOOGLE_API_KEY`: From terraform variable
- `TELEGRAPH_API_KEY`: From terraform variable
- `TELEGRAPH_AUTHOR_NAME`: From terraform variable
- `AUTH_SECRET`: From terraform variable
- `AUTH_TOKEN_MESSAGES`: From terraform variable
- `MAIL_SMTP_SERVER`: From terraform variable
- `MAIL_SMTP_PORT`: From terraform variable
- `MAIL_SENDER_EMAIL`: From terraform variable
- `MAIL_SENDER_PASSWORD`: From terraform variable
- `MAIL_RECIPIENT_EMAIL`: From terraform variable
- `GOOGLE_CLOUD_PROJECT`: From terraform variable (project_id)
- `CLOUD_RUN_JOB_NAME`: Automatically set to the Cloud Run Job name
- `CLOUD_RUN_JOB_REGION`: From terraform variable (region)

### Cloud Run Job Environment Variables

The Cloud Run Job has the same environment variables as the service (except `CLOUD_RUN_JOB_NAME` and `CLOUD_RUN_JOB_REGION` which are only needed by the service).

## Updating the Deployment

To update the application after making code changes:

1. Rebuild and redeploy:
   ```bash
   PROJECT_ID=your-project-id make deploy
   ```

   Cloud Build will automatically:
   - Build the new container image
   - Push it to Artifact Registry
   - Update the Cloud Run service with the new image
   - Update the Cloud Run Job with the new image

## Troubleshooting

### API Gateway vs Cloud Run Service

- **API Gateway URL**: Use this for all production API calls. This is the public-facing endpoint.
- **Cloud Run Service URL**: This is not directly accessible from the internet. It's only accessible via the API Gateway.

To get both URLs:
```bash
terraform output api_gateway_url      # Use this for API calls
terraform output cloud_run_service_url # Internal URL (for reference only)
```

### View API Gateway Logs

```bash
gcloud api-gateway gateways describe research-assistant-gateway --location us-central1
```

### View Cloud Run Service Logs

```bash
gcloud run services logs read research-assistant --region us-central1
```

### View Cloud Run Job Logs

```bash
gcloud run jobs executions list --job research-worker --region us-central1
gcloud run jobs executions logs read EXECUTION_NAME --job research-worker --region us-central1
```

### Check API Gateway Status

```bash
gcloud api-gateway gateways describe research-assistant-gateway --location us-central1
gcloud api-gateway apis describe research-assistant-api
gcloud api-gateway api-configs list --api=research-assistant-api
```

### Check Service Status

```bash
gcloud run services describe research-assistant --region us-central1
```

### Check Job Status

```bash
gcloud run jobs describe research-worker --region us-central1
```

### Update Environment Variables

Edit `terraform/terraform.tfvars` and run:
```bash
make terraform-apply
```

This will update both the Cloud Run service and Cloud Run Job with new environment variables.

**Note:** If you update the OpenAPI specification (`openapi.yaml`), you'll need to run `terraform apply` again to update the API Gateway configuration. The API config uses a timestamp-based ID to ensure updates are properly deployed.

## Cleanup

To destroy all infrastructure:

```bash
make terraform-destroy
```

Or:
```bash
cd terraform && terraform destroy
```
