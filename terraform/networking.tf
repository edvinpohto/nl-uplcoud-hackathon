module "networking" {
  source = "./modules/networking"

  prefix         = var.prefix
  zone           = var.zone
  network_cidr_a = "10.10.1.0/24"
  network_cidr_b = "10.10.2.0/24"
}
