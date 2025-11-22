terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "apigateway.googleapis.com",
    "servicemanagement.googleapis.com",
    "servicecontrol.googleapis.com",
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

  ingress = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    service_account = google_service_account.gemini_assistant.email

    containers {
      name = "research-assistant"
      # Use latest tagged image from Artifact Registry, or custom image if specified
      image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.research_assistant.repository_id}/research-assistant:latest"

      ports {
        container_port = 8080
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
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }

      env {
        name  = "CLOUD_RUN_JOB_NAME"
        value = google_cloud_run_v2_job.research_worker.name
      }

      env {
        name  = "CLOUD_RUN_JOB_REGION"
        value = var.region
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
    }

    timeout = "60s"
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.research_assistant,
  ]
}

# Create Cloud Run Job for processing research requests
resource "google_cloud_run_v2_job" "research_worker" {
  name     = "research-worker"
  location = var.region

  template {
    template {
      service_account = google_service_account.gemini_assistant.email

      containers {
        image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.research_assistant.repository_id}/research-assistant:latest"

        command = ["/worker"]

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

        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
        }

        resources {
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
      }

      timeout     = "3600s"
      max_retries = 2
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_service_account.gemini_assistant,
  ]
}

# API Gateway API configuration
resource "google_api_gateway_api" "research_assistant_api" {
  provider     = google-beta
  api_id       = "research-assistant-api"
  display_name = "Research Assistant API"
  project      = var.project_id

  depends_on = [google_project_service.required_apis]
}

# API Gateway API Config - uses the OpenAPI spec
resource "google_api_gateway_api_config" "research_assistant_api_config" {
  provider    = google-beta
  api         = google_api_gateway_api.research_assistant_api.api_id
  api_config_id = "research-assistant-api-config-${formatdate("YYYYMMDDhhmmss", timestamp())}"

  openapi_documents {
    document {
      path     = "openapi.yaml"
      contents = base64encode(templatefile("${path.module}/../openapi.yaml", {
        CLOUD_RUN_SERVICE_URL = google_cloud_run_v2_service.research_assistant.uri
        API_ID                = "research-assistant-api"
      }))
    }
  }

  gateway_config {
    backend_config {
      google_service_account = google_service_account.api_gateway.email
    }
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    google_api_gateway_api.research_assistant_api,
    google_cloud_run_v2_service.research_assistant,
    google_service_account.api_gateway,
  ]
}

# API Gateway Gateway
resource "google_api_gateway_gateway" "research_assistant_gateway" {
  provider   = google-beta
  region     = var.region
  project    = var.project_id
  api_config = google_api_gateway_api_config.research_assistant_api_config.id
  gateway_id = "research-assistant-gateway"

  depends_on = [google_api_gateway_api_config.research_assistant_api_config]
}

# Grant API Gateway service account permission to invoke Cloud Run service
resource "google_cloud_run_v2_service_iam_member" "api_gateway_invoker" {
  location = google_cloud_run_v2_service.research_assistant.location
  name     = google_cloud_run_v2_service.research_assistant.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.api_gateway.email}"

  depends_on = [google_service_account.api_gateway]
}

