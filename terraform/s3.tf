resource "random_string" "this" {
  count   = var.config_storage_type == "s3" && var.create_config_file_bucket ? 1 : 0
  length  = 6
  special = false
  upper   = false
}

resource "aws_s3_bucket" "this" {
  count         = var.config_storage_type == "s3" && var.create_config_file_bucket ? 1 : 0
  bucket        = format("%s-yace-config-%s", var.name_prefix, random_string.this[0].result)
  force_destroy = true
}

resource "aws_s3_object" "config" {
    count  = var.config_storage_type == "s3" ? 1 : 0
    key = var.config_path == "" ? format("%s-yace-config/config.yaml", var.name_prefix) : var.config_path
    bucket = var.create_config_file_bucket ? aws_s3_bucket.this[0].bucket : var.config_bucket
    source = var.config_file_local_path
}