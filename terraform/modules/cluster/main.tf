terraform {
  required_providers {
    upcloud = {
      source  = "UpCloudLtd/upcloud"
      version = "~> 5.0"
    }
  }
}

resource "upcloud_kubernetes_cluster" "this" {
  name                  = var.cluster_name
  zone                  = var.zone
  private_node_groups   = true
  network               = var.network_uuid

  control_plane_ip_filter = var.control_plane_ip_filter

  labels = merge(var.labels, {
    managed-by = "terraform"
  })
}

resource "upcloud_kubernetes_node_group" "this" {
  cluster    = upcloud_kubernetes_cluster.this.id
  name       = "${var.cluster_name}-workers"
  node_count = var.node_count
  plan       = var.node_plan

  labels = {
    managed-by = "terraform"
    role       = "worker"
  }
}
