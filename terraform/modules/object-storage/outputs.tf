output "service_uuid" {
  description = "UUID of the object storage service"
  value       = upcloud_managed_object_storage.this.id
}

output "endpoint" {
  description = "Endpoint domain for the object storage"
  value       = "https://${upcloud_managed_object_storage.this.endpoint[0].domain_name}"
}

output "access_key_id" {
  description = "Access key ID"
  value       = upcloud_managed_object_storage_user_access_key.this.access_key_id
  sensitive   = true
}

output "secret_access_key" {
  description = "Secret access key"
  value       = upcloud_managed_object_storage_user_access_key.this.secret_access_key
  sensitive   = true
}
