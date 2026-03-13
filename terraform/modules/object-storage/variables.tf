variable "name" {
  description = "Name of the object storage instance"
  type        = string
}

variable "region" {
  description = "UpCloud Object Storage region"
  type        = string
  default     = "europe-1"
}

variable "buckets" {
  description = "List of bucket names to create"
  type        = list(string)
  default     = []
}

variable "labels" {
  description = "Labels to apply to the object storage"
  type        = map(string)
  default     = {}
}
