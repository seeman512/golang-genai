#!/usr/bin/env bash
# Shared build and AWS/SAM deployment environment.
set -u

export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"
export AWS_REGION="${AWS_REGION:-$AWS_DEFAULT_REGION}"
export AWS_PAGER=""

case "$(uname -m)" in
  arm64|aarch64) export GOARCH_LAMBDA=arm64 LAMBDA_ARCH=arm64 ;;
  *)             export GOARCH_LAMBDA=amd64 LAMBDA_ARCH=x86_64 ;;
esac

export PRESENTATION_FN="lesson-12-presentation"
export WEBHOOK_FN="lesson-12-presentation-webhook"
export AVATAR_FN="lesson-12-user-avatar"
export USER_EVENTS_FN="lesson-12-user-events-webhook"
export API_NAME="lesson-12-presentation-api"
export STAGE="dev"
export ROLE_NAME="lesson-12-lambda-role"

# Terraform consumes these values automatically. Existing TF_VAR_* values
# remain overridable by callers.
export TF_VAR_region="${TF_VAR_region:-$AWS_DEFAULT_REGION}"
export TF_VAR_lambda_architecture="${TF_VAR_lambda_architecture:-$LAMBDA_ARCH}"
export TF_VAR_stage_name="${TF_VAR_stage_name:-$STAGE}"
export TF_VAR_lambda_role_name="${TF_VAR_lambda_role_name:-$ROLE_NAME}"
