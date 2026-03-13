variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "zone" {
  description = "UpCloud zone"
  type        = string
}

variable "network_uuid" {
  description = "UUID of the private SDN network to attach to"
  type        = string
}

variable "node_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "node_plan" {
  description = "UpCloud server plan for worker nodes"
  type        = string
  default     = "2xCPU-4GB"
}

variable "control_plane_ip_filter" {
  description = "IP filter for control plane access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "labels" {
  description = "Labels to apply to the cluster"
  type        = map(string)
  default     = {}
}
