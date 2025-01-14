# Generates a resource name based on the standardized prefix and additional parameters.
output "resource_name" {
  value = provider::dx::resource_name("dx-d-itn", "app", "blob_private_endpoint", 1)
}