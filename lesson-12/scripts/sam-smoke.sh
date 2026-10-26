#!/usr/bin/env bash
# Smoke-test the SAM-built functions and local API Gateway emulator.
set -euo pipefail
cd "$(dirname "$0")/.."
source scripts/env.sh

command -v sam >/dev/null || { echo "AWS SAM CLI is required" >&2; exit 1; }
command -v curl >/dev/null || { echo "curl is required" >&2; exit 1; }

sam build --template-file sam/template.yaml \
  --parameter-overrides "ParameterKey=LambdaArchitecture,ParameterValue=${LAMBDA_ARCH}"

avatar_output=build/sam-avatar.json
sam local invoke UserAvatarFunction \
  --template .aws-sam/build/template.yaml \
  --event events/user-avatar.json >"$avatar_output"
grep -q '"initials": "AL"' "$avatar_output" || grep -q '"initials":"AL"' "$avatar_output"
echo "OK: user-avatar Lambda"

mkdir -p build
sam local start-api \
  --template .aws-sam/build/template.yaml \
  --host 127.0.0.1 \
  --port 3000 >build/sam-api.log 2>&1 &
sam_pid=$!
cleanup() {
  kill "$sam_pid" >/dev/null 2>&1 || true
  wait "$sam_pid" >/dev/null 2>&1 || true
}
trap cleanup EXIT

for _ in $(seq 1 60); do
  if curl -sS -o /dev/null http://127.0.0.1:3000/ >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

url=http://127.0.0.1:3000/user-events
status=$(curl -sS -o build/sam-user-events-response.json -w '%{http_code}' \
  -X POST -H 'Content-Type: application/json' \
  --data-binary @events/user-registered-body.json "$url")
[[ "$status" == 200 ]]
grep -q 'Welcome to Mystic Square' build/sam-user-events-response.json
echo "OK: user-events webhook (HTTP ${status})"

status=$(curl -sS -o /dev/null -w '%{http_code}' \
  -X POST -H 'Content-Type: application/json' \
  --data-binary @events/user-registered-invalid.json "$url")
[[ "$status" == 400 ]]
echo "OK: invalid user event (HTTP ${status})"
