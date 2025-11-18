output "service_account_email" {
  description = "Email of the created service account"
  value       = google_service_account.gemini_assistant.email
}

output "service_account_id" {
  description = "ID of the created service account"
  value       = google_service_account.gemini_assistant.id
}

output "service_account_name" {
  description = "Name of the created service account"
  value       = google_service_account.gemini_assistant.name
}

output "cloud_run_service_url" {
  description = "URL of the deployed Cloud Run service"
  value       = google_cloud_run_v2_service.research_assistant.uri
}

output "artifact_registry_repository" {
  description = "Full path to the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.research_assistant.repository_id}"
}

output "container_image" {
  description = "Full container image path"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.research_assistant.repository_id}/research-assistant:latest"
}

