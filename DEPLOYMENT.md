# Deployment Guide

This guide explains how to build and deploy the Research Assistant application to Google Cloud Run.

## Prerequisites

1. **Google Cloud Project**: You need a GCP project with billing enabled
2. **gcloud CLI**: Install and authenticate with `gcloud auth login`
3. **Docker**: For local builds (optional)
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
- Service account with necessary permissions
- Required GCP APIs enabled

**Note:** The Cloud Run service will be created with a placeholder image (`gcr.io/cloudrun/hello`) since the actual application image doesn't exist yet. The first deployment in Step 3 will replace this placeholder with your actual application image.

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

### Option B: Local Docker Build and Push

1. Authenticate Docker with Artifact Registry:
   ```bash
   gcloud auth configure-docker us-central1-docker.pkg.dev
   ```

2. Build and push the image:
   ```bash
   PROJECT_ID=your-project-id make docker-push
   ```

3. Update the Cloud Run service to use the new image:
   ```bash
   gcloud run services update research-assistant \
     --image us-central1-docker.pkg.dev/your-project-id/research-assistant/research-assistant:latest \
     --region us-central1
   ```

## Step 4: Verify Deployment

1. Get the service URL:
   ```bash
   make terraform-output
   ```

   Or directly:
   ```bash
   cd terraform && terraform output cloud_run_service_url
   ```

2. Test the service:
   ```bash
   # Generate an auth token first (see below)
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "https://YOUR-SERVICE-URL/research?topic=test"
   ```

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

- `PORT`: Set to 8080
- `TELEGRAPH_API_KEY`: From terraform variable
- `TELEGRAPH_AUTHOR_NAME`: From terraform variable
- `AUTH_SECRET`: From terraform variable
- `AUTH_TOKEN_MESSAGES`: From terraform variable
- `MAIL_SMTP_SERVER`: From terraform variable
- `MAIL_SMTP_PORT`: From terraform variable
- `MAIL_SENDER_EMAIL`: From terraform variable
- `MAIL_SENDER_PASSWORD`: From terraform variable
- `MAIL_RECIPIENT_EMAIL`: From terraform variable

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

## Troubleshooting

### View Cloud Run Logs

```bash
gcloud run services logs read research-assistant --region us-central1
```

### Check Service Status

```bash
gcloud run services describe research-assistant --region us-central1
```

### Update Environment Variables

Edit `terraform/terraform.tfvars` and run:
```bash
make terraform-apply
```

This will update the Cloud Run service with new environment variables.

## Cleanup

To destroy all infrastructure:

```bash
make terraform-destroy
```

Or:
```bash
cd terraform && terraform destroy
```
