# Terraform deployment targets AWS by default. Set local_aws=true to target an
# AWS-compatible Docker emulator such as Floki/LocalStack instead.
provider "aws" {
  region = var.region

  # AWS-compatible local emulators accept arbitrary non-empty credentials.
  access_key = var.local_aws ? "test" : null
  secret_key = var.local_aws ? "test" : null

  skip_credentials_validation = var.local_aws
  skip_metadata_api_check     = var.local_aws
  skip_requesting_account_id  = var.local_aws
  skip_region_validation      = var.local_aws

  dynamic "endpoints" {
    for_each = var.local_aws ? [1] : []

    content {
      apigateway = var.local_aws_endpoint
      iam        = var.local_aws_endpoint
      lambda     = var.local_aws_endpoint
      logs       = var.local_aws_endpoint
      sts        = var.local_aws_endpoint
    }
  }
}
