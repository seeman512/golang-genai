#!/usr/bin/env bash
# Deploy the lesson through Terraform/OpenTofu.
# Usage: scripts/deploy-tf.sh [plan|apply|destroy]
set -euo pipefail
cd "$(dirname "$0")/.."
source scripts/env.sh

TF=$(command -v terraform || command -v tofu || true)
[[ -n "$TF" ]] || { echo "Terraform >= 1.5 or OpenTofu is required" >&2; exit 1; }

mode=${1:-plan}
state_args=()
plan_args=()
if [[ -n "${TF_STATE_FILE:-}" ]]; then
  state_args+=("-state=${TF_STATE_FILE}")
fi
if [[ -n "${TF_PLAN_FILE:-}" ]]; then
  plan_args+=("-out=${TF_PLAN_FILE}")
fi

case "$mode" in
  plan|apply)
    for package in \
      build/presentation.zip \
      build/presentation-webhook.zip \
      build/user-avatar.zip \
      build/user-events-webhook.zip; do
      [[ -f "$package" ]] || { echo "Missing $package; run make build first" >&2; exit 1; }
    done
    ;;
  destroy) ;;
  *) echo "Usage: $0 [plan|apply|destroy]" >&2; exit 2 ;;
esac

"$TF" -chdir=terraform init -input=false
"$TF" -chdir=terraform validate

case "$mode" in
  plan)
    "$TF" -chdir=terraform plan -input=false "${state_args[@]}" "${plan_args[@]}"
    echo
    if [[ -n "${TF_PLAN_FILE:-}" ]]; then
      echo "Plan saved to terraform/${TF_PLAN_FILE}. Review it before applying."
    else
      echo "Plan only. Review it, then run: make deploy-tf"
    fi
    ;;
  apply)
    "$TF" -chdir=terraform apply -input=false -auto-approve "${state_args[@]}"
    echo
    "$TF" -chdir=terraform output
    ;;
  destroy)
    "$TF" -chdir=terraform destroy -input=false -auto-approve "${state_args[@]}"
    ;;
esac
