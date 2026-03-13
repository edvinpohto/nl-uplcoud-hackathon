output "network_a_uuid" {
  description = "UUID of Cluster A's private network"
  value       = upcloud_network.cluster_a.id
}

output "network_b_uuid" {
  description = "UUID of Cluster B's private network"
  value       = upcloud_network.cluster_b.id
}

output "router_uuid" {
  description = "UUID of the shared router"
  value       = upcloud_router.chaos.id
}
