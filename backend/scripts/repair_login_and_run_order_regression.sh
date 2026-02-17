#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REGRESSION_SCRIPT="${SCRIPT_DIR}/run_order_mysql8_regression.sh"

DMH_INTEGRATION_BASE_URL="${DMH_INTEGRATION_BASE_URL:-http://127.0.0.1:8889}"
DMH_TEST_ADMIN_USERNAME="${DMH_TEST_ADMIN_USERNAME:-admin}"
DMH_TEST_ADMIN_PASSWORD="${DMH_TEST_ADMIN_PASSWORD:-123456}"
MYSQL_CONTAINER="${MYSQL_CONTAINER:-mysql8}"
API_CONTAINER="${API_CONTAINER:-dmh-api}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-Admin168}"
DMH_DB_NAME="${DMH_DB_NAME:-dmh}"
MAX_WAIT_SECONDS="${MAX_WAIT_SECONDS:-30}"

DEFAULT_BCRYPT_HASH='$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm'

usage() {
  cat <<'EOF'
Usage:
  backend/scripts/repair_login_and_run_order_regression.sh [base-url]

What it does:
  1) Check login using DMH_TEST_ADMIN_USERNAME/DMH_TEST_ADMIN_PASSWORD.
  2) If login fails, repair known local test users' bcrypt password hash.
  3) Restart dmh-api and wait for login to recover.
  4) Run backend/scripts/run_order_mysql8_regression.sh.

Environment variables:
  DMH_INTEGRATION_BASE_URL   Default: http://127.0.0.1:8889
  DMH_TEST_ADMIN_USERNAME    Default: admin
  DMH_TEST_ADMIN_PASSWORD    Default: 123456
  MYSQL_CONTAINER            Default: mysql8
  API_CONTAINER              Default: dmh-api
  MYSQL_ROOT_PASSWORD        Default: Admin168
  DMH_DB_NAME                Default: dmh
  MAX_WAIT_SECONDS           Default: 30

Notes:
  - Repair step updates users(admin, brand_manager) to default test password hash.
  - This script is for local/integration environment only.
EOF
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[repair-order-regression] ERROR: missing command: $1"
    exit 1
  fi
}

ensure_running_container() {
  local name="$1"
  if ! docker ps --format '{{.Names}}' | grep -qx "$name"; then
    echo "[repair-order-regression] ERROR: container not running: $name"
    exit 1
  fi
}

LOGIN_BODY_FILE="$(mktemp)"
trap 'rm -f "$LOGIN_BODY_FILE"' EXIT

login_http_code() {
  curl -sS -o "$LOGIN_BODY_FILE" -w '%{http_code}' \
    -X POST "${DMH_INTEGRATION_BASE_URL}/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"${DMH_TEST_ADMIN_USERNAME}\",\"password\":\"${DMH_TEST_ADMIN_PASSWORD}\"}"
}

login_ok() {
  local code
  code="$(login_http_code)"
  [ "$code" = "200" ]
}

repair_password_hashes() {
  echo "[repair-order-regression] repairing local test user hashes in ${MYSQL_CONTAINER}.${DMH_DB_NAME}"
  docker exec "$MYSQL_CONTAINER" mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -D "${DMH_DB_NAME}" \
    -e "UPDATE users SET password='${DEFAULT_BCRYPT_HASH}' WHERE username IN ('admin', 'brand_manager');"
}

wait_login_recovered() {
  local elapsed=0
  while [ "$elapsed" -lt "$MAX_WAIT_SECONDS" ]; do
    if login_ok; then
      return 0
    fi
    sleep 2
    elapsed=$((elapsed + 2))
  done
  return 1
}

main() {
  if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
    usage
    exit 0
  fi

  if [ "$#" -ge 1 ] && [ -n "${1:-}" ]; then
    DMH_INTEGRATION_BASE_URL="$1"
  fi

  require_command curl
  require_command docker

  if [ ! -x "$REGRESSION_SCRIPT" ]; then
    echo "[repair-order-regression] ERROR: regression script not executable: $REGRESSION_SCRIPT"
    exit 1
  fi

  echo "[repair-order-regression] base URL: ${DMH_INTEGRATION_BASE_URL}"
  echo "[repair-order-regression] login user: ${DMH_TEST_ADMIN_USERNAME}"

  if login_ok; then
    echo "[repair-order-regression] login check passed, run regression directly"
  else
    echo "[repair-order-regression] login check failed, body: $(cat "$LOGIN_BODY_FILE")"

    ensure_running_container "$MYSQL_CONTAINER"
    ensure_running_container "$API_CONTAINER"

    repair_password_hashes

    echo "[repair-order-regression] restarting ${API_CONTAINER}"
    docker restart "$API_CONTAINER" >/dev/null

    echo "[repair-order-regression] waiting for login recovery (max ${MAX_WAIT_SECONDS}s)"
    if ! wait_login_recovered; then
      echo "[repair-order-regression] ERROR: login still failing after repair and restart"
      echo "[repair-order-regression] latest login response: $(cat "$LOGIN_BODY_FILE")"
      exit 1
    fi
    echo "[repair-order-regression] login recovered"
  fi

  DMH_INTEGRATION_BASE_URL="$DMH_INTEGRATION_BASE_URL" \
  DMH_TEST_ADMIN_USERNAME="$DMH_TEST_ADMIN_USERNAME" \
  DMH_TEST_ADMIN_PASSWORD="$DMH_TEST_ADMIN_PASSWORD" \
    "$REGRESSION_SCRIPT"
}

main "$@"
