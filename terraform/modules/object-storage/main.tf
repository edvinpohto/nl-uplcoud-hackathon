terraform {
  required_providers {
    upcloud = {
      source  = "UpCloudLtd/upcloud"
      version = "~> 5.0"
    }
  }
}

resource "upcloud_managed_object_storage" "this" {
  name              = var.name
  region            = var.region
  configured_status = "started"

  labels = var.labels
}

resource "upcloud_managed_object_storage_user" "this" {
  service_uuid = upcloud_managed_object_storage.this.id
  username     = "${var.name}-user"
}

resource "upcloud_managed_object_storage_user_access_key" "this" {
  service_uuid = upcloud_managed_object_storage.this.id
  username     = upcloud_managed_object_storage_user.this.username
  status       = "Active"
}

resource "upcloud_managed_object_storage_bucket" "this" {
  for_each     = toset(var.buckets)
  service_uuid = upcloud_managed_object_storage.this.id
  name         = each.value
}
