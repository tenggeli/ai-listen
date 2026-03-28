#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="ai-listen"
BACKEND_DIR="/home/hwd/app/ai-listen/backend"
LOG_DIR="${HOME}/log"
TIMESTAMP="$(date +%Y%m%d%H%M%S)"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}-${TIMESTAMP}.log"
PID_FILE="${LOG_DIR}/${PROJECT_NAME}.pid"

mkdir -p "${LOG_DIR}"

cd "${BACKEND_DIR}"

# 优先使用已编译二进制，若不存在则使用 go run
if [[ -x "./bin/server" ]]; then
  nohup ./bin/server > "${LOG_FILE}" 2>&1 &
else
  nohup go run ./cmd/server > "${LOG_FILE}" 2>&1 &
fi

PID=$!
echo "${PID}" > "${PID_FILE}"

echo "Backend started in background"
echo "PID: ${PID}"
echo "Log: ${LOG_FILE}"
echo "PID file: ${PID_FILE}"
