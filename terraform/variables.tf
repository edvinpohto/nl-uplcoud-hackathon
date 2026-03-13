variable "upcloud_token" {
  description = "UpCloud API token"
  type        = string
  sensitive   = true
}

variable "zone" {
  description = "UpCloud zone for clusters"
  type        = string
  default     = "de-fra1"
}

variable "prefix" {
  description = "Resource name prefix"
  type        = string
  default     = "chaos-monkey"
}

variable "network_cidr" {
  description = "CIDR for the shared SDN network"
  type        = string
  default     = "10.10.0.0/22"
}

variable "cluster_a_node_count" {
  description = "Number of worker nodes in Cluster A (victim)"
  type        = number
  default     = 3
}

variable "cluster_b_node_count" {
  description = "Number of worker nodes in Cluster B (chaos + monitor)"
  type        = number
  default     = 2
}

variable "node_plan" {
  description = "UpCloud server plan for worker nodes"
  type        = string
  default     = "2xCPU-4GB"
}

variable "docker_registry" {
  description = "Docker registry for images (e.g. registry.example.com/chaos-monkey)"
  type        = string
  default     = "docker.io/chaosmonkey"
}

variable "image_tag" {
  description = "Docker image tag"
  type        = string
  default     = "latest"
}
