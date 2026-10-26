locals {
  lambda_functions = {
    presentation = {
      function_name = "lesson-12-presentation"
      package       = "${path.module}/../build/presentation.zip"
    }
    presentation_webhook = {
      function_name = "lesson-12-presentation-webhook"
      package       = "${path.module}/../build/presentation-webhook.zip"
    }
    user_avatar = {
      function_name = "lesson-12-user-avatar"
      package       = "${path.module}/../build/user-avatar.zip"
    }
    user_events_webhook = {
      function_name = "lesson-12-user-events-webhook"
      package       = "${path.module}/../build/user-events-webhook.zip"
    }
  }

  api_routes = {
    webhook       = "presentation_webhook"
    "user-events" = "user_events_webhook"
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = var.lambda_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

# The functions only need permissions to write their execution logs. They do
# not access Mystic Square's database or user data directly in Stage 1.
resource "aws_iam_role_policy" "lambda_logs" {
  name = "lesson-12-lambda-logs"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "CreateLogGroup"
        Effect   = "Allow"
        Action   = ["logs:CreateLogGroup"]
        Resource = "*"
      },
      {
        Sid      = "WriteExecutionLogs"
        Effect   = "Allow"
        Action   = ["logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "arn:aws:logs:*:*:log-group:/aws/lambda/lesson-12-*:*"
      }
    ]
  })
}

resource "aws_lambda_function" "function" {
  for_each = local.lambda_functions

  function_name = each.value.function_name
  role          = aws_iam_role.lambda_exec.arn
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = [var.lambda_architecture]

  filename         = each.value.package
  source_code_hash = filebase64sha256(each.value.package)

  timeout     = 10
  memory_size = 256

  depends_on = [aws_iam_role_policy.lambda_logs]
}

resource "aws_api_gateway_rest_api" "api" {
  name = var.api_name
}

resource "aws_api_gateway_resource" "route" {
  for_each = local.api_routes

  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = each.key
}

resource "aws_api_gateway_method" "post" {
  for_each = local.api_routes

  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.route[each.key].id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "post" {
  for_each = local.api_routes

  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.route[each.key].id
  http_method             = aws_api_gateway_method.post[each.key].http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.function[each.value].arn}/invocations"
}

resource "aws_lambda_permission" "api_gateway" {
  for_each = local.api_routes

  statement_id  = "AllowApiGateway${replace(each.key, "-", "")}"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function[each.value].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/POST/${each.key}"
}

resource "aws_api_gateway_deployment" "api" {
  rest_api_id = aws_api_gateway_rest_api.api.id

  triggers = {
    redeployment = sha1(jsonencode([
      for route in sort(keys(local.api_routes)) : {
        path        = route
        resource_id = aws_api_gateway_resource.route[route].id
        method_id   = aws_api_gateway_method.post[route].id
        integration = aws_api_gateway_integration.post[route].id
      }
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [aws_api_gateway_integration.post, aws_lambda_permission.api_gateway]
}

resource "aws_api_gateway_stage" "stage" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  deployment_id = aws_api_gateway_deployment.api.id
  stage_name    = var.stage_name
}
