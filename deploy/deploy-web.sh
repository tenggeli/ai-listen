#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

WEB_ROOT="${WEB_ROOT:-/var/www/html}"
STAGE_ROOT="${STAGE_ROOT:-${HOME}/webdeploy}"
STAGE_TS="$(date +%Y%m%d%H%M%S)"

APP_NAMES=("admin-web" "user-web" "provider-web")
APP_DIRS=(
  "${REPO_ROOT}/web/apps/admin-web"
  "${REPO_ROOT}/web/apps/user-web"
  "${REPO_ROOT}/web/apps/provider-web"
)
APP_TARGETS=(
  "${WEB_ROOT}/admin-web"
  "${WEB_ROOT}/user-web"
  "${WEB_ROOT}/provider-web"
)
APP_STAGES=()

build_and_stage_app() {
  local app_dir="$1"
  local app_name="$2"
  local stage_dir="$3"

  if [[ ! -d "${app_dir}" ]]; then
    echo "${app_name} directory not found: ${app_dir}"
    exit 1
  fi

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

for i in "${!APP_NAMES[@]}"; do
  stage_dir="${STAGE_ROOT}/${APP_NAMES[$i]}-${STAGE_TS}"
  build_and_stage_app "${APP_DIRS[$i]}" "${APP_NAMES[$i]}" "${stage_dir}"
  APP_STAGES+=("${stage_dir}")
done

for i in "${!APP_NAMES[@]}"; do
  deploy_app "${APP_STAGES[$i]}" "${APP_TARGETS[$i]}"
done

echo "[nginx] check and reload"
sudo nginx -t
sudo systemctl reload nginx

echo "Deploy finished successfully."
