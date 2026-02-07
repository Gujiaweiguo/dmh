#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

# Optional first argument overrides base URL.
if [ "$#" -ge 1 ] && [ -n "${1:-}" ]; then
  export DMH_INTEGRATION_BASE_URL="$1"
fi

export DMH_INTEGRATION_BASE_URL="${DMH_INTEGRATION_BASE_URL:-http://localhost:8889}"
export DMH_TEST_ADMIN_USERNAME="${DMH_TEST_ADMIN_USERNAME:-admin}"
export DMH_TEST_ADMIN_PASSWORD="${DMH_TEST_ADMIN_PASSWORD:-123456}"

echo "[order-mysql8-regression] base URL: ${DMH_INTEGRATION_BASE_URL}"
echo "[order-mysql8-regression] admin user: ${DMH_TEST_ADMIN_USERNAME}"

cd "$ROOT_DIR"

OUT_FILE="$(mktemp)"
trap 'rm -f "$OUT_FILE"' EXIT

go test ./test/integration -run 'TestOrderVerifyRoutesAuthGuard|TestOrderCreateDuplicateMessage' -count=1 -v | tee "$OUT_FILE"

if grep -E -q -- '--- SKIP: TestOrderVerifyRoutesAuthGuard|--- SKIP: TestOrderCreateDuplicateMessage' "$OUT_FILE"; then
  echo ""
  echo "[order-mysql8-regression] ERROR: tests were skipped."
  echo "Check API availability and test credentials:"
  echo "  DMH_INTEGRATION_BASE_URL=${DMH_INTEGRATION_BASE_URL}"
  echo "  DMH_TEST_ADMIN_USERNAME=${DMH_TEST_ADMIN_USERNAME}"
  echo "  DMH_TEST_ADMIN_PASSWORD=<your-password>"
  exit 2
fi

echo "[order-mysql8-regression] PASS"
