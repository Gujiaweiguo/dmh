#!/usr/bin/env bash
#
# Integration Test Runner
# Standardized environment setup and test execution
#
# Usage:
#   ./run_integration_tests.sh                    # Run all integration tests
#   ./run_integration_tests.sh -run TestOrder     # Run specific tests
#
# Environment variables:
#   DMH_INTEGRATION_BASE_URL   - API base URL (default: http://localhost:8889)
#   DMH_TEST_ADMIN_USERNAME    - Admin username (default: admin)
#   DMH_TEST_ADMIN_PASSWORD    - Admin password (default: 123456)
#   MYSQL_TEST_HOST            - MySQL host (default: 127.0.0.1)
#   MYSQL_TEST_PORT            - MySQL port (default: 3306)
#   REDIS_TEST_HOST            - Redis address (default: localhost:6379)

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default configuration
export DMH_INTEGRATION_BASE_URL="${DMH_INTEGRATION_BASE_URL:-http://localhost:8889}"
export DMH_TEST_ADMIN_USERNAME="${DMH_TEST_ADMIN_USERNAME:-admin}"
export DMH_TEST_ADMIN_PASSWORD="${DMH_TEST_ADMIN_PASSWORD:-123456}"
export MYSQL_TEST_HOST="${MYSQL_TEST_HOST:-127.0.0.1}"
export MYSQL_TEST_PORT="${MYSQL_TEST_PORT:-3306}"
export MYSQL_TEST_USER="${MYSQL_TEST_USER:-root}"
export MYSQL_TEST_PASSWORD="${MYSQL_TEST_PASSWORD:-Admin168}"
export MYSQL_TEST_DB="${MYSQL_TEST_DB:-dmh}"
export REDIS_TEST_HOST="${REDIS_TEST_HOST:-localhost:6379}"

# Timeout constants
API_READY_TIMEOUT=20
MYSQL_READY_TIMEOUT=30

echo ""
echo "========================================"
echo "  DMH Integration Test Runner"
echo "========================================"
echo ""
echo "Environment:"
echo "  API URL:      ${DMH_INTEGRATION_BASE_URL}"
echo "  Admin User:   ${DMH_TEST_ADMIN_USERNAME}"
echo "  MySQL Host:   ${MYSQL_TEST_HOST}:${MYSQL_TEST_PORT}"
echo "  Redis Host:   ${REDIS_TEST_HOST}"
echo ""

# Function to check API readiness
check_api() {
    echo -n "Checking API readiness... "
    
    for i in $(seq 1 $API_READY_TIMEOUT); do
        code=$(curl -s -o /dev/null -w "%{http_code}" "${DMH_INTEGRATION_BASE_URL}/api/v1/auth/login" 2>/dev/null || echo "000")
        if [ "$code" != "000" ]; then
            echo -e "${GREEN}OK${NC} (status: $code)"
            return 0
        fi
        sleep 0.5
    done
    
    echo -e "${RED}FAILED${NC}"
    echo ""
    echo "API at ${DMH_INTEGRATION_BASE_URL} is not responding."
    echo ""
    echo "Troubleshooting:"
    echo "  1. Start the API: cd backend && go run ./api/dmh.go -f ./api/etc/dmh-api.yaml"
    echo "  2. Or use Docker: cd deploy && docker compose -f docker-compose-simple.yml up -d"
    echo ""
    return 1
}

# Function to check MySQL readiness
check_mysql() {
    echo -n "Checking MySQL readiness... "
    
    for i in $(seq 1 $MYSQL_READY_TIMEOUT); do
        if mysqladmin ping -h"${MYSQL_TEST_HOST}" -P"${MYSQL_TEST_PORT}" -u"${MYSQL_TEST_USER}" -p"${MYSQL_TEST_PASSWORD}" --silent 2>/dev/null; then
            echo -e "${GREEN}OK${NC}"
            return 0
        fi
        sleep 1
    done
    
    echo -e "${YELLOW}SKIPPED${NC} (MySQL not required for all tests)"
    return 0
}

# Function to check Redis availability
check_redis() {
    echo -n "Checking Redis availability... "
    
    if timeout 2 bash -c "echo PING | nc -w 1 ${REDIS_TEST_HOST%:*} ${REDIS_TEST_HOST#*:} 2>/dev/null" | grep -q PONG; then
        echo -e "${GREEN}OK${NC}"
        export REDIS_AVAILABLE=true
    else
        echo -e "${YELLOW}NOT AVAILABLE${NC} (optional, tests will use memory fallback)"
        export REDIS_AVAILABLE=false
    fi
    return 0
}

# Function to analyze test results
analyze_results() {
    local output_file="$1"
    
    # Count results
    local passed=$(grep -c "^--- PASS" "$output_file" || echo 0)
    local failed=$(grep -c "^--- FAIL" "$output_file" || echo 0)
    local skipped=$(grep -c "^--- SKIP" "$output_file" || echo 0)
    local total=$((passed + failed + skipped))
    
    echo ""
    echo "========================================"
    echo "  Test Results Summary"
    echo "========================================"
    echo ""
    echo "  Total:   $total"
    echo -e "  ${GREEN}Passed:  $passed${NC}"
    echo -e "  ${RED}Failed:  $failed${NC}"
    echo -e "  ${YELLOW}Skipped: $skipped${NC}"
    
    if [ "$total" -gt 0 ]; then
        local execution_rate=$(( (passed + failed) * 100 / total ))
        echo ""
        echo "  Execution Rate: ${execution_rate}%"
    fi
    
    # Show skip reasons if any
    if [ "$skipped" -gt 0 ]; then
        echo ""
        echo "Skip Reasons:"
        grep -oE 'SKIP_REASON:[^|]+' "$output_file" | sort | uniq -c | while read count reason; do
            echo "  $count $reason"
        done
    fi
    
    echo ""
    
    # Return non-zero if any tests failed
    if [ "$failed" -gt 0 ]; then
        return 1
    fi
    
    # Return non-zero if tests were skipped (CI should fail on skips)
    if [ "$skipped" -gt 0 ]; then
        echo -e "${YELLOW}WARNING: Some tests were skipped. This may indicate environment issues.${NC}"
        return 2
    fi
    
    return 0
}

# Main execution
main() {
    # Run environment checks
    check_api || exit 1
    check_mysql
    check_redis
    
    echo ""
    echo "Running integration tests..."
    echo ""
    
    # Create temp file for output
    OUTPUT_FILE=$(mktemp)
    trap 'rm -f "$OUTPUT_FILE"' EXIT
    
    # Run tests with provided args (or default to all integration tests)
    TEST_ARGS="${*:-./test/integration/...}"
    
    if ! go test $TEST_ARGS -v -count=1 -timeout 10m 2>&1 | tee "$OUTPUT_FILE"; then
        # Test command failed, but we still want to analyze results
        :
    fi
    
    # Analyze results
    analyze_results "$OUTPUT_FILE"
}

main "$@"
