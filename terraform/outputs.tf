output "amp_workspace_endpoint" {
  description = "The endpoint of the Amazon Managed Prometheus workspace."
  value       = var.create_amp_workspace ? aws_prometheus_workspace.this[0].prometheus_endpoint : null
}
