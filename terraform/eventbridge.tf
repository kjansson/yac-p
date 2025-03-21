module "scheduler" {
  source = "terraform-aws-modules/eventbridge/aws"

  create_bus = false

  rules = {
    format("%s-cron", var.name_prefix) = {
      description         = format("%s-cron invoker rule", var.name_prefix)
      schedule_expression = "rate(5 minutes)"
    }
  }

  targets = {
    crons = [
      {
        name  = format("%s-lambda-target", var.name_prefix)
        arn   = aws_lambda_function.this.arn
        input = jsonencode({ "job" : "cron-by-rate" })
      }
    ]
  }
  tags = var.tags
}