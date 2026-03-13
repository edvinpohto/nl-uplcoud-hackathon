variable "prefix" {
  description = "Name prefix for resources"
  type        = string
  default     = "chaos-monkey"
}

variable "zone" {
  description = "UpCloud zone"
  type        = string
}

variable "network_cidr_a" {
  description = "CIDR for Cluster A's private network"
  type        = string
  default     = "10.10.1.0/24"
}

variable "network_cidr_b" {
  description = "CIDR for Cluster B's private network"
  type        = string
  default     = "10.10.2.0/24"
}
