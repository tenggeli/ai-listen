#!/usr/bin/env bash
set -euo pipefail

ADMIN_WEB_DIR="/home/hwd/app/ai-listen/web/apps/admin-web"
USER_WEB_DIR="/home/hwd/app/ai-listen/web/apps/user-web"
WEB_ROOT="/var/www/html"
STAGE_ROOT="${HOME}/webdeploy"
STAGE_TS="$(date +%Y%m%d%H%M%S)"

ADMIN_TARGET="${WEB_ROOT}/admin-web"
USER_TARGET="${WEB_ROOT}/user-web"
ADMIN_STAGE="${STAGE_ROOT}/admin-web-${STAGE_TS}"
USER_STAGE="${STAGE_ROOT}/user-web-${STAGE_TS}"

build_and_stage_app() {
  local app_dir="$1"
  local app_name="$2"
  local stage_dir="$3"

  echo "[build] ${app_name} and stage to ${stage_dir}"
  cd "${app_dir}"
  npm ci
  npm run build

  if [[ ! -d "${app_dir}/dist" ]]; then
    echo "${app_name} dist not found: ${app_dir}/dist"
    exit 1
  fi

  mkdir -p "${STAGE_ROOT}"
  rm -rf "${stage_dir}"
  mkdir -p "${stage_dir}"
  cp -a "${app_dir}/dist/." "${stage_dir}/"
}

deploy_app() {
  local src_stage="$1"
  local target_dir="$2"

  local tmp_dir="${target_dir}.tmp.${STAGE_TS}"

  echo "[deploy] ${src_stage} -> ${target_dir}"

  sudo mkdir -p "${WEB_ROOT}"
  sudo rm -rf "${tmp_dir}"
  sudo mkdir -p "${tmp_dir}"
  sudo cp -a "${src_stage}/." "${tmp_dir}/"

  # 原子切换：先备份旧目录，再替换
  if sudo test -d "${target_dir}"; then
    sudo rm -rf "${target_dir}.bak"
    sudo mv "${target_dir}" "${target_dir}.bak"
  fi

  sudo mv "${tmp_dir}" "${target_dir}"
  sudo rm -rf "${target_dir}.bak"
}

build_and_stage_app "${ADMIN_WEB_DIR}" "admin-web" "${ADMIN_STAGE}"
build_and_stage_app "${USER_WEB_DIR}" "user-web" "${USER_STAGE}"

deploy_app "${ADMIN_STAGE}" "${ADMIN_TARGET}"
deploy_app "${USER_STAGE}" "${USER_TARGET}"

echo "[nginx] check and reload"
sudo nginx -t
sudo systemctl reload nginx

echo "Deploy finished successfully."
