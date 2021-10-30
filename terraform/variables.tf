variable "project_id" {
  type        = string
  description = "The project ID to manage the Cloud SQL resources"
}

variable "region" {
  type        = string
  description = "The region of the Cloud SQL resources"
  default     = "us-central1"
}

variable "db_user" {
  type        = string
  description = "The db user for login"
}

variable "db_password" {
  type        = string
  description = "The db password for login"
}

variable "db_name" {
  type        = string
  description = "The name for pgsql resources"
}
