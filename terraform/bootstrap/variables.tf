variable "upcloud_token" {
  description = "UpCloud API token"
  type        = string
  sensitive   = true
}

variable "object_storage_region" {
  description = "UpCloud Object Storage region"
  type        = string
  default     = "europe-1"
}
