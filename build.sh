#!/usr/bin/env bash

set -euo pipefail

go mod tidy

PLATFORMS=(
  "windows amd64 .exe"
  "windows arm64 .exe"
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

mkdir -p bin

for entry in "${PLATFORMS[@]}"; do
  read -r GOOS GOARCH SUFFIX <<< "$entry"
  SUFFIX="${SUFFIX:-}"  # Default to empty string if not set
  echo "Building for $GOOS/$GOARCH..."
  export GOOS
  export GOARCH
  go build -o "bin/process-reporter-${GOOS}-${GOARCH}${SUFFIX}"
done

echo "All builds complete."
