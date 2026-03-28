#!/usr/bin/env bash
set -euo pipefail

ADMIN_WEB_DIR="/home/hwd/app/ai-listen/web/apps/admin-web"
USER_WEB_DIR="/home/hwd/app/ai-listen/web/apps/user-web"
WEB_ROOT="/var/www/html"

ADMIN_TARGET="${WEB_ROOT}/admin-web"
USER_TARGET="${WEB_ROOT}/user-web"

echo "[1/4] Build admin-web"
cd "${ADMIN_WEB_DIR}"
npm ci
npm run build

echo "[2/4] Build user-web"
cd "${USER_WEB_DIR}"
npm ci
npm run build

echo "[3/4] Deploy dist files (requires sudo)"
if [[ ! -d "${ADMIN_WEB_DIR}/dist" ]]; then
  echo "admin-web dist not found: ${ADMIN_WEB_DIR}/dist"
  exit 1
fi
if [[ ! -d "${USER_WEB_DIR}/dist" ]]; then
  echo "user-web dist not found: ${USER_WEB_DIR}/dist"
  exit 1
fi

echo mkdir -p "${ADMIN_TARGET}"
sudo mkdir -p "${ADMIN_TARGET}"
sudo mkdir -p "${USER_TARGET}"
echo rm -rf "${ADMIN_TARGET:?}/"*
sudo rm -rf "${ADMIN_TARGET:?}/"*
sudo rm -rf "${USER_TARGET:?}/"*
echo cp -a "${ADMIN_WEB_DIR}/dist/." "${ADMIN_TARGET}/"
sudo cp -a "${ADMIN_WEB_DIR}/dist/." "${ADMIN_TARGET}/"
sudo cp -a "${USER_WEB_DIR}/dist/." "${USER_TARGET}/"

echo "[4/4] Validate and reload nginx"
sudo nginx -t
sudo systemctl reload nginx

echo "Deploy finished successfully."
