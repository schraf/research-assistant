terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
  ])

  project = var.project_id
  service = each.value

  disable_dependent_services = false
}

# Create Artifact Registry repository for container images
resource "google_artifact_registry_repository" "research_assistant" {
  location      = var.region
  repository_id = "research-assistant"
  description   = "Docker repository for Research Assistant container images"
  format        = "DOCKER"

  depends_on = [google_project_service.required_apis]
}

# Cloud Run service
resource "google_cloud_run_v2_service" "research_assistant" {
  name     = "research-assistant"
  location = var.region

  template {
    service_account = google_service_account.gemini_assistant.email

    containers {
      name  = "research-assistant"
      # Use placeholder image initially - will be replaced by first deployment
      image = var.initial_image != "" ? var.initial_image : "gcr.io/cloudrun/hello"

      ports {
        container_port = 8080
      }

      env {
        name  = "GOOGLE_API_KEY"
        value = var.google_api_key
      }

      env {
        name  = "TELEGRAPH_API_KEY"
        value = var.telegraph_api_key
      }

      env {
        name  = "TELEGRAPH_AUTHOR_NAME"
        value = var.telegraph_author_name
      }

      env {
        name  = "AUTH_SECRET"
        value = var.auth_secret
      }

      env {
        name  = "AUTH_TOKEN_MESSAGES"
        value = var.auth_token_messages
      }

      env {
        name  = "MAIL_SMTP_SERVER"
        value = var.smtp_hostname
      }

      env {
        name  = "MAIL_SMTP_PORT"
        value = var.smtp_port
      }

      env {
        name  = "MAIL_SENDER_EMAIL"
        value = var.mail_sender_email
      }

      env {
        name  = "MAIL_SENDER_PASSWORD"
        value = var.mail_sender_password
      }

      env {
        name  = "MAIL_RECIPIENT_EMAIL"
        value = var.mail_recipient_email
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
    }

    timeout = "1800s"
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.research_assistant,
  ]
}

# Allow unauthenticated access to the Cloud Run service
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  location = google_cloud_run_v2_service.research_assistant.location
  name     = google_cloud_run_v2_service.research_assistant.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
