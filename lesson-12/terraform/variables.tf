variable "region" {
  description = "AWS region for the deployment."
  type        = string
  default     = "us-east-1"
}

variable "lambda_architecture" {
  description = "Lambda architecture matching the compiled Go binaries."
  type        = string
  default     = "x86_64"

  validation {
    condition     = contains(["x86_64", "arm64"], var.lambda_architecture)
    error_message = "lambda_architecture must be x86_64 or arm64."
  }
}

variable "stage_name" {
  description = "API Gateway deployment stage."
  type        = string
  default     = "Prod"
}

variable "api_name" {
  description = "API Gateway REST API name."
  type        = string
  default     = "lesson-12-presentation-api"
}

variable "lambda_role_name" {
  description = "Execution role shared by the lesson Lambdas."
  type        = string
  default     = "lesson-12-lambda-role"
}
