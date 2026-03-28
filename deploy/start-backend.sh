#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="ai-listen"
BACKEND_DIR="/home/hwd/app/ai-listen/backend"
LOG_DIR="${HOME}/log"
TIMESTAMP="$(date +%Y%m%d%H%M%S)"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}-${TIMESTAMP}.log"
PID_FILE="${LOG_DIR}/${PROJECT_NAME}.pid"

mkdir -p "${LOG_DIR}"

stop_running_service() {
  if [[ -f "${PID_FILE}" ]]; then
    local old_pid
    old_pid="$(cat "${PID_FILE}")"
    if [[ -n "${old_pid}" ]] && kill -0 "${old_pid}" 2>/dev/null; then
      echo "Stopping existing backend process: ${old_pid}"
      kill "${old_pid}" 2>/dev/null || true

      for _ in {1..10}; do
        if ! kill -0 "${old_pid}" 2>/dev/null; then
          break
        fi
        sleep 1
      done

      if kill -0 "${old_pid}" 2>/dev/null; then
        echo "Force killing process: ${old_pid}"
        kill -9 "${old_pid}" 2>/dev/null || true
      fi
    fi
  fi

  rm -f "${PID_FILE}"
}

stop_running_service

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
