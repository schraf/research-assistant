variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-central1"
}

variable "google_api_key" {
  description = "API key for accessing Google Gemini"
  type        = string
  sensitive   = true
}

variable "telegraph_api_key" {
  description = "API key for posting to Telegra.ph"
  type        = string
  sensitive   = true
}

variable "telegraph_author_name" {
  description = "Author name for Telegraph articles"
  type        = string
  default     = ""
}

variable "auth_secret" {
  description = "Secret key for generating and validating auth tokens"
  type        = string
  sensitive   = true
}

variable "auth_token_messages" {
  description = "Comma-separated list of messages that can be used to generate valid auth tokens"
  type        = string
  sensitive   = true
}

variable "smtp_hostname" {
  description = "Hostname for a smtp server"
  type        = string
  default     = "smtp.gmail.com"
}

variable "smtp_port" {
  description = "Port number for a smtp server"
  type        = string
  default     = "587"
}

variable "mail_sender_email" {
  description = "The email address for the sender"
  type        = string
  sensitive   = true
}

variable "mail_sender_password" {
  description = "The email password for the sender"
  type        = string
  sensitive   = true
}

variable "mail_recipient_email" {
  description = "The recipient email address"
  type        = string
}

variable "initial_image" {
  description = "Initial container image to use (defaults to placeholder). Leave empty to use placeholder image."
  type        = string
  default     = ""
}
