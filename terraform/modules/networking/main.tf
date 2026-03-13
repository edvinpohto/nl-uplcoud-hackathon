terraform {
  required_providers {
    upcloud = {
      source  = "UpCloudLtd/upcloud"
      version = "~> 5.0"
    }
  }
}

resource "upcloud_router" "chaos" {
  name = "${var.prefix}-router"
}

# Separate network for each cluster — UpCloud only allows one cluster per network.
# Both networks share the same router so nodes can communicate across clusters.

resource "upcloud_network" "cluster_a" {
  name = "${var.prefix}-net-a"
  zone = var.zone

  ip_network {
    address            = var.network_cidr_a
    dhcp               = true
    dhcp_default_route = true
    family             = "IPv4"
    gateway            = cidrhost(var.network_cidr_a, 1)
  }

  router = upcloud_router.chaos.id
}

resource "upcloud_network" "cluster_b" {
  name = "${var.prefix}-net-b"
  zone = var.zone

  ip_network {
    address            = var.network_cidr_b
    dhcp               = true
    dhcp_default_route = true
    family             = "IPv4"
    gateway            = cidrhost(var.network_cidr_b, 1)
  }

  router = upcloud_router.chaos.id
}

resource "upcloud_gateway" "chaos" {
  name     = "${var.prefix}-gateway"
  zone     = var.zone
  features = ["nat"]

  router {
    id = upcloud_router.chaos.id
  }
}
