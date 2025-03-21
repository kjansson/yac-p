resource "aws_ssm_parameter" "config" {
  count      = var.config_storage_type == "ssm" ? 1 : 0
  name        = var.config_path == "" ? format("%s-yace-config", var.name_prefix) : var.config_path
  description = "YACE discovery config"
  type        = "String"
  value       = file(var.config_file_local_path)

  tags = var.tags
}