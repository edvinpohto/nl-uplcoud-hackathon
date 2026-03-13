terraform {
  required_providers {
    upcloud = {
      source  = "UpCloudLtd/upcloud"
      version = "~> 5.0"
    }
  }
}

provider "upcloud" {
  token = var.upcloud_token
}

resource "upcloud_managed_object_storage" "tfstate" {
  name              = "chaos-monkey-tfstate"
  region            = var.object_storage_region
  configured_status = "started"

  network {
    family = "IPv4"
    name   = "public"
    type   = "public"
  }

  labels = {
    purpose = "terraform-state"
    project = "chaos-monkey"
  }
}

resource "upcloud_managed_object_storage_user" "tfstate" {
  service_uuid = upcloud_managed_object_storage.tfstate.id
  username     = "tfstate-user"
}

resource "upcloud_managed_object_storage_user_access_key" "tfstate" {
  service_uuid = upcloud_managed_object_storage.tfstate.id
  username     = upcloud_managed_object_storage_user.tfstate.username
  status       = "Active"
}

resource "upcloud_managed_object_storage_bucket" "tfstate" {
  service_uuid = upcloud_managed_object_storage.tfstate.id
  name         = "terraform-state"
}

resource "upcloud_managed_object_storage_user_policy" "tfstate" {
  service_uuid = upcloud_managed_object_storage.tfstate.id
  username     = upcloud_managed_object_storage_user.tfstate.username
  name         = "AmazonS3FullAccess"
}

output "object_storage_endpoint" {
  value = "https://${tolist(upcloud_managed_object_storage.tfstate.endpoint)[0].domain_name}"
}

output "access_key" {
  value     = upcloud_managed_object_storage_user_access_key.tfstate.access_key_id
  sensitive = true
}

output "secret_key" {
  value     = upcloud_managed_object_storage_user_access_key.tfstate.secret_access_key
  sensitive = true
}
