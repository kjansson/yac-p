variable "name_prefix" {
  description = "Prefix to apply to all resources."
  type        = string
}

variable "create_amp_workspace" {
  description = "Create an Amazon Managed Prometheus workspace."
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = any
  default     = {}
}

variable "create_lambda_log_group" {
  description = "Create a CloudWatch log group for the Lambda function."
  type        = bool
  default     = true
}

variable "lambda_log_group_retention" {
  description = "Retention period for the Lambda log group"
  type        = number
  default     = 7
}

variable "lambda_log_group_name" {
  description = "Name of the Lambda log group. Setting create_lambda_log_group to true will override this value and use the created log group name."
  type        = string
  default     = ""
}

variable "prometheus_endpoint" {
  description = "The endpoint of the Prometheus workspace. Setting create_amp_workspace to true will override this value and use the created workspace endpoint."
  type        = string
  default     = ""
}

variable "prometheus_region" {
  description = "The region of the Prometheus workspace. Only used for Amazon Managed Prometheus."
  type        = string
}

variable "lambda_image_uri" {
  description = "The URI of the Lambda function image."
  type        = string
}

variable "lambda_image_repository_arn" {
  description = "The ARN of the Lambda function image repository."
  type        = string
}

variable "lambda_log_level" {
  description = "The log level for the Lambda function."
  type        = string
  default     = "INFO"
}

variable "lambda_schedule_rate" {
  description = "The rate at which to invoke the Lambda function (in minutes)."
  type        = number
  default     = 5
}

variable "config_file_local_path" {
  description = "Path to the local YACE config file to upload and use."
  type        = string
}

variable "create_config_file_bucket" {
  description = "Create an S3 bucket to store the YACE config file. Overrides config_bucket."
  type        = bool
  default     = false
}

variable "config_bucket" {
  description = "Name of existing S3 bucket to store the YACE config file. Ignored if create_config_file_bucket is true."
  type        = string
  default     = ""
}

variable "config_storage_type" {
  description = "Service used to store YACE config file for Lambda. Valid options are ssm or s3. Default is ssm."
  type        = string
  default     = "ssm"
  validation {
    condition = can(regex("^(ssm|s3)$", var.config_storage_type))
    error_message = "config_storage_type must be either 'ssm' or 's3'"
  }
}

variable "config_path" {
  description = "Custom path to the YACE config file in the chosen storage service."
  default     = ""
}
