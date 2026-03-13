output "cluster_id" {
  description = "UUID of the Kubernetes cluster"
  value       = upcloud_kubernetes_cluster.this.id
}

output "cluster_name" {
  description = "Name of the Kubernetes cluster"
  value       = upcloud_kubernetes_cluster.this.name
}

output "node_group_name" {
  description = "Name of the node group"
  value       = upcloud_kubernetes_node_group.this.name
}
