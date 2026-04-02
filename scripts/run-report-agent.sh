#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
[ -f ./shellforge ] || go build -o shellforge ./cmd/shellforge
exec ./shellforge report "${1:-.}"
