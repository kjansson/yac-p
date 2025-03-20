resource "aws_ssm_parameter" "config" {
  name        = format("%s-yace-config", var.name_prefix)
  description = "YACE discovery config"
  type        = "String"
  value       = file(var.yace_config_file_path)

  tags = var.tags
}