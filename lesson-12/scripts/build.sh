#!/usr/bin/env bash
# Build OS-only Lambda packages. provided.al2023 expects an executable named
# bootstrap at the root of each zip archive.
set -euo pipefail
cd "$(dirname "$0")/.."
source scripts/env.sh

command -v go >/dev/null || { echo "Go is required" >&2; exit 1; }
command -v zip >/dev/null || { echo "zip is required" >&2; exit 1; }

rm -rf build
mkdir -p build/presentation build/presentation-webhook build/user-avatar build/user-events-webhook

build_function() {
  local name=$1
  echo "Building ${name} for linux/${GOARCH_LAMBDA}"
  GOOS=linux GOARCH="$GOARCH_LAMBDA" CGO_ENABLED=0 \
    go build -tags lambda.norpc -ldflags='-s -w' \
    -o "build/${name}/bootstrap" "./cmd/${name}"
  (cd "build/${name}" && zip -q -X "../${name}.zip" bootstrap)
}

build_function presentation
build_function presentation-webhook
build_function user-avatar
build_function user-events-webhook
ls -lh build/*.zip
