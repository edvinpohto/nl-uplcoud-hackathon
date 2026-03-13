output "cluster_a_id" {
  description = "UUID of Cluster A (victim)"
  value       = module.cluster_a.cluster_id
}

output "cluster_a_name" {
  description = "Name of Cluster A"
  value       = module.cluster_a.cluster_name
}

output "cluster_b_id" {
  description = "UUID of Cluster B (chaos)"
  value       = module.cluster_b.cluster_id
}

output "cluster_b_name" {
  description = "Name of Cluster B"
  value       = module.cluster_b.cluster_name
}

output "network_a_uuid" {
  description = "UUID of Cluster A's network"
  value       = module.networking.network_a_uuid
}

output "network_b_uuid" {
  description = "UUID of Cluster B's network"
  value       = module.networking.network_b_uuid
}
