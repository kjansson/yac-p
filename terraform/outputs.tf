output "amp_workspace_endpoint" {
  description = "The endpoint of the Amazon Managed Prometheus workspace."
  value       = aws_prometheus_workspace.this[0].prometheus_endpoint
}
