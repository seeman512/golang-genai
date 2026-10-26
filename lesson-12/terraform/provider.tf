# Terraform deployment for a real AWS account. Local Lambda testing uses SAM;
# this provider is intentionally not configured for an AWS emulator.
provider "aws" {
  region = var.region
}
