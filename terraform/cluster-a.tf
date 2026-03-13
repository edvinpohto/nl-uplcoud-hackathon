# Cluster A — Victim cluster
# Runs victim-app and Pumba DaemonSet for network disruption
module "cluster_a" {
  source = "./modules/cluster"

  cluster_name = "${var.prefix}-cluster-a"
  zone         = var.zone
  network_uuid = module.networking.network_a_uuid
  node_count   = var.cluster_a_node_count
  node_plan    = var.node_plan

  # Allow Cluster B nodes to reach the control plane over SDN
  control_plane_ip_filter = ["0.0.0.0/0"]

  labels = {
    role    = "victim"
    project = "chaos-monkey"
  }

  depends_on = [module.networking]
}
