#!/usr/bin/env bash
set -euo pipefail

ADMIN_WEB_DIR="/home/hwd/app/ai-listen/web/apps/admin-web"
USER_WEB_DIR="/home/hwd/app/ai-listen/web/apps/user-web"
WEB_ROOT="/var/www/html"

ADMIN_TARGET="${WEB_ROOT}/admin-web"
USER_TARGET="${WEB_ROOT}/user-web"

build_app() {
  local app_dir="$1"
  local app_name="$2"

  echo "[build] ${app_name}"
  cd "${app_dir}"
  npm ci
  npm run build

  if [[ ! -d "${app_dir}/dist" ]]; then
    echo "${app_name} dist not found: ${app_dir}/dist"
    exit 1
  fi
}

deploy_app() {
  local src_dist="$1"
  local target_dir="$2"

  local ts
  ts="$(date +%Y%m%d%H%M%S)"
  local tmp_dir="${target_dir}.tmp.${ts}"

  echo "[deploy] ${src_dist} -> ${target_dir}"

  sudo mkdir -p "${WEB_ROOT}"
  sudo rm -rf "${tmp_dir}"
  sudo mkdir -p "${tmp_dir}"
  sudo cp -a "${src_dist}/." "${tmp_dir}/"

  # 原子切换：先备份旧目录，再替换
  if sudo test -d "${target_dir}"; then
    sudo rm -rf "${target_dir}.bak"
    sudo mv "${target_dir}" "${target_dir}.bak"
  fi

  sudo mv "${tmp_dir}" "${target_dir}"
  sudo rm -rf "${target_dir}.bak"
}

build_app "${ADMIN_WEB_DIR}" "admin-web"
build_app "${USER_WEB_DIR}" "user-web"

deploy_app "${ADMIN_WEB_DIR}/dist" "${ADMIN_TARGET}"
deploy_app "${USER_WEB_DIR}/dist" "${USER_TARGET}"

echo "[nginx] check and reload"
sudo nginx -t
sudo systemctl reload nginx

echo "Deploy finished successfully."
