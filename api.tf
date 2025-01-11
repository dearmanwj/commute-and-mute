data "aws_cloudfront_distribution" "cloudfront" {
  id = "EWKILBEV49A1M"
}

resource "aws_apigatewayv2_api" "cam_users_api" {
  name          = "cam-users-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = [ "https://${data.aws_cloudfront_distribution.cloudfront.domain_name}"]
    allow_methods = ["GET", "POST", "PUT", "OPTIONS"]
    allow_headers = ["authorization"]
  }
}

resource "aws_apigatewayv2_integration" "cam_users_integration" {
  api_id           = aws_apigatewayv2_api.cam_users_api.id
  integration_type = "AWS_PROXY"

  description            = "CAM User lambda integration"
  integration_uri        = aws_lambda_function.users_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_route" {
  api_id             = aws_apigatewayv2_api.cam_users_api.id
  route_key          = "GET /users"
  target             = "integrations/${aws_apigatewayv2_integration.cam_users_integration.id}"
  authorizer_id      = aws_apigatewayv2_authorizer.api.id
  authorization_type = "CUSTOM"
}

resource "aws_apigatewayv2_route" "put_route" {
  api_id             = aws_apigatewayv2_api.cam_users_api.id
  route_key          = "PUT /users"
  target             = "integrations/${aws_apigatewayv2_integration.cam_users_integration.id}"
  authorizer_id      = aws_apigatewayv2_authorizer.api.id
  authorization_type = "CUSTOM"
}

resource "aws_apigatewayv2_route" "delete_route" {
  api_id             = aws_apigatewayv2_api.cam_users_api.id
  route_key          = "DELETE /users"
  target             = "integrations/${aws_apigatewayv2_integration.cam_users_integration.id}"
  authorizer_id      = aws_apigatewayv2_authorizer.api.id
  authorization_type = "CUSTOM"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.cam_users_api.id
  name        = "$default"
  auto_deploy = "true"
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api.arn
    format          = "$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] $context.httpMethod $context.resourcePath $context.protocol $context.status $context.responseLength $context.requestId $context.extendedRequestId"
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "API-Gateway-Execution-Logs_${aws_apigatewayv2_api.cam_users_api.id}"
  retention_in_days = 1
}

resource "aws_apigatewayv2_authorizer" "api" {
  api_id                            = aws_apigatewayv2_api.cam_users_api.id
  authorizer_type                   = "REQUEST"
  authorizer_uri                    = aws_lambda_function.authorizer_lambda.invoke_arn
  identity_sources                  = ["$request.header.Authorization"]
  name                              = "cam-api-authorizer"
  authorizer_payload_format_version = "2.0"
  enable_simple_responses           = "true"
}
