# Create the service account
resource "google_service_account" "gemini_assistant" {
  account_id   = "gemini-assistant"
  display_name = "Gemini Assistant Service Account"
  description  = "Service account for Research Assistant Cloud Run service"
}

# Grant the service account permission to use Gemini API
resource "google_project_iam_member" "gemini_assistant_aiplatform" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.gemini_assistant.email}"
}

# Grant Cloud Build service account permission to push images to Artifact Registry
data "google_project" "project" {
  project_id = var.project_id
}

resource "google_artifact_registry_repository_iam_member" "cloudbuild_writer" {
  location   = var.region
  repository = google_artifact_registry_repository.research_assistant.repository_id
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}

# Grant Cloud Build service account permission to deploy to Cloud Run
# Note: This uses project-level IAM to avoid circular dependencies
resource "google_project_iam_member" "cloudbuild_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}

