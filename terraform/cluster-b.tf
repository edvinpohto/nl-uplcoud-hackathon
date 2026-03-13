# Cluster B — Chaos + Monitoring cluster
# Runs chaos-monkey, health-checker, Prometheus, Grafana
module "cluster_b" {
  source = "./modules/cluster"

  cluster_name = "${var.prefix}-cluster-b"
  zone         = var.zone
  network_uuid = module.networking.network_b_uuid
  node_count   = var.cluster_b_node_count
  node_plan    = var.node_plan

  control_plane_ip_filter = ["0.0.0.0/0"]

  labels = {
    role    = "chaos"
    project = "chaos-monkey"
  }

  depends_on = [module.networking]
}
