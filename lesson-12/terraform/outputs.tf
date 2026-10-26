output "api_id" {
  description = "API Gateway REST API identifier."
  value       = aws_api_gateway_rest_api.api.id
}

output "api_base_url" {
  description = "Base URL for the deployed API in AWS."
  value       = "https://${aws_api_gateway_rest_api.api.id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_stage.stage.stage_name}"
}

output "webhook_url" {
  description = "Generic webhook URL."
  value       = "https://${aws_api_gateway_rest_api.api.id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_stage.stage.stage_name}/webhook"
}

output "user_events_url" {
  description = "Registration welcome webhook URL."
  value       = "https://${aws_api_gateway_rest_api.api.id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_stage.stage.stage_name}/user-events"
}

output "lambda_function_names" {
  description = "Lambda function names managed by this configuration."
  value       = { for key, function in aws_lambda_function.function : key => function.function_name }
}
