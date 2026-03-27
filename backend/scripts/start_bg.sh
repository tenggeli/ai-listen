#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="${HOME}/log"
TIMESTAMP="$(date +"%Y%m%d%H%M%S")"
LOG_FILE="${LOG_DIR}/backend-${TIMESTAMP}.log"

mkdir -p "${LOG_DIR}"

cd "${ROOT_DIR}"

HTTP_ADDR="${HTTP_ADDR:-:18080}" nohup go run ./cmd/server >"${LOG_FILE}" 2>&1 &
PID=$!

echo "backend started"
echo "pid=${PID}"
echo "addr=${HTTP_ADDR:-:18080}"
echo "log=${LOG_FILE}"
