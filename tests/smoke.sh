#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"
go run ./cmd/supplier-skill version
go test ./internal/cli ./internal/mcpclient ./internal/installer ./internal/workbuddy
