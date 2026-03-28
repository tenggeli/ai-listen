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
sudo rm -rf "${ADMIN_TARGET}"
sudo rm -rf "${USER_TARGET}"
sudo mv "${ADMIN_WEB_DIR}/dist/*" "${ADMIN_TARGET}"
sudo mv "${USER_WEB_DIR}/dist/*" "${USER_TARGET}"

echo "[4/4] Validate and reload nginx"
sudo nginx -t
sudo systemctl reload nginx

echo "Deploy finished successfully."
