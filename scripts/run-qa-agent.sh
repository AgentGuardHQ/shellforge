#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

if [[ ! -f ./shellforge ]]; then
  echo "[run-qa-agent] Building shellforge..."
  go build -o shellforge ./cmd/shellforge
fi

exec ./shellforge qa "${1:-.}"
