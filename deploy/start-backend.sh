#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="ai-listen"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BACKEND_DIR="${BACKEND_DIR:-${REPO_ROOT}/backend}"
LOG_DIR="${LOG_DIR:-${HOME}/log}"
PORT="${LISTEN_SERVER_PORT:-$(awk -F= '/^[[:space:]]*server\.port[[:space:]]*=/{gsub(/[[:space:]]*/, "", $2); print $2; exit}' "${HOME}/conf/listenbase.cof" 2>/dev/null || true)}"
PORT="${PORT:-8080}"
START_WAIT_SECONDS="${START_WAIT_SECONDS:-3}"
TIMESTAMP="$(date +%Y%m%d%H%M%S)"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}-${TIMESTAMP}.log"
PID_FILE="${LOG_DIR}/${PROJECT_NAME}.pid"

mkdir -p "${LOG_DIR}"

is_backend_process() {
  local pid="$1"
  local cmd
  cmd="$(ps -p "${pid}" -o command= 2>/dev/null || true)"
  [[ "${cmd}" == *"/ai-listen/backend/"* || "${cmd}" == *"/ai-listen/backend"* || "${cmd}" == *"go run ./cmd/server"* || "${cmd}" == *"/bin/server"* ]]
}

find_listen_pid_by_port() {
  local pid=""
  if command -v lsof >/dev/null 2>&1; then
    pid="$(lsof -tiTCP:"${PORT}" -sTCP:LISTEN 2>/dev/null | head -n 1 || true)"
  elif command -v ss >/dev/null 2>&1; then
    pid="$(ss -ltnp 2>/dev/null | awk -v p=":${PORT}" '$4 ~ p { if (match($0, /pid=[0-9]+/)) { print substr($0, RSTART+4, RLENGTH-4); exit } }')"
  fi
  echo "${pid}"
}

stop_pid_gracefully() {
  local pid="$1"
  if [[ -z "${pid}" ]] || ! kill -0 "${pid}" 2>/dev/null; then
    return
  fi
  kill "${pid}" 2>/dev/null || true
  for _ in {1..10}; do
    if ! kill -0 "${pid}" 2>/dev/null; then
      return
    fi
    sleep 1
  done
  kill -9 "${pid}" 2>/dev/null || true
}

stop_running_service() {
  if [[ -f "${PID_FILE}" ]]; then
    local old_pid
    old_pid="$(cat "${PID_FILE}")"
    if [[ -n "${old_pid}" ]] && kill -0 "${old_pid}" 2>/dev/null; then
      echo "Stopping existing backend process from pid file: ${old_pid}"
      stop_pid_gracefully "${old_pid}"
    fi
  fi

  rm -f "${PID_FILE}"

  local port_pid
  port_pid="$(find_listen_pid_by_port)"
  if [[ -n "${port_pid}" ]] && kill -0 "${port_pid}" 2>/dev/null; then
    if is_backend_process "${port_pid}"; then
      echo "Port ${PORT} is occupied by old backend process ${port_pid}, stopping it"
      stop_pid_gracefully "${port_pid}"
    else
      echo "Port ${PORT} is occupied by another process (PID=${port_pid})."
      echo "Please release the port or start backend with LISTEN_SERVER_PORT=<new_port>."
      exit 1
    fi
  fi
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

sleep "${START_WAIT_SECONDS}"
if ! kill -0 "${PID}" 2>/dev/null; then
  echo "Backend failed to stay running. Recent logs:"
  tail -n 60 "${LOG_FILE}" || true
  rm -f "${PID_FILE}"
  exit 1
fi

port_pid="$(find_listen_pid_by_port)"
if [[ -z "${port_pid}" ]]; then
  echo "Backend process is running (PID=${PID}), but port ${PORT} is not listening yet."
  echo "Check logs: ${LOG_FILE}"
  exit 1
fi

echo "Backend started in background"
echo "PID: ${PID}"
echo "Port: ${PORT}"
echo "Log: ${LOG_FILE}"
echo "PID file: ${PID_FILE}"
