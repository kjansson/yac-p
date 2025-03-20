resource "aws_prometheus_workspace" "this" {
  count = var.create_amp_workspace ? 1 : 0
  alias = format("%s-amp", var.name_prefix)
  tags  = var.tags
}