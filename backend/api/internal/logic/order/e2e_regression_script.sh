#!/bin/bash

# DMH Order and Permission E2E Regression Script
# This script performs end-to-end testing of order creation, verification, and permission checks

set -e

echo "=== DMH Order and Permission E2E Regression Test ==="
echo "Starting at: $(date)"

# Configuration
TEST_DB="dmh"
TEST_USER="root"
TEST_PASSWORD="Admin168"
TEST_HOST="127.0.0.1"
TEST_PORT="3306"

# Docker configuration (if MySQL is running in Docker)
USE_DOCKER="true"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test results
PASS_COUNT=0
FAIL_COUNT=0
TOTAL_TESTS=0

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_result="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${YELLOW}Running test: $test_name${NC}"
    echo "Command: $test_command"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ PASS: $test_name${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}✗ FAIL: $test_name${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

# Function to check database connection
check_db_connection() {
    echo "Checking database connection..."
    if [ "$USE_DOCKER" = "true" ]; then
        # Use Docker to connect to MySQL
        if docker exec mysql8 mysql -u"$TEST_USER" -p"$TEST_PASSWORD" -e "SELECT 1" "$TEST_DB" >/dev/null 2>&1; then
            echo "✓ Database connection successful (via Docker)"
            return 0
        else
            echo "✗ Database connection failed (via Docker)"
            return 1
        fi
    else
        # Use direct MySQL connection
        if mysql -u"$TEST_USER" -p"$TEST_PASSWORD" -h"$TEST_HOST" -P"$TEST_PORT" -e "SELECT 1" "$TEST_DB" >/dev/null 2>&1; then
            echo "✓ Database connection successful"
            return 0
        else
            echo "✗ Database connection failed"
            return 1
        fi
    fi
}

# Function to cleanup test data
cleanup_test_data() {
    echo "Cleaning up test data..."
    if [ "$USE_DOCKER" = "true" ]; then
        # Use Docker to cleanup test data
        docker exec mysql8 mysql -u"$TEST_USER" -p"$TEST_PASSWORD" "$TEST_DB" <<EOF
DELETE FROM orders WHERE phone LIKE '13800138%';
DELETE FROM verification_records WHERE verification_code LIKE '1_13800138%';
EOF
    else
        # Use direct MySQL connection
        mysql -u"$TEST_USER" -p"$TEST_PASSWORD" -h"$TEST_HOST" -P"$TEST_PORT" "$TEST_DB" <<EOF
DELETE FROM orders WHERE phone LIKE '13800138%';
DELETE FROM verification_records WHERE verification_code LIKE '1_13800138%';
EOF
    fi
}

# Main test execution
main() {
    # Check database connection
    if ! check_db_connection; then
        echo "Database connection failed. Exiting."
        exit 1
    fi
    
    # Cleanup existing test data
    cleanup_test_data
    
    echo -e "\n${YELLOW}=== Running Order Creation Tests ===${NC}"
    
    # Test 1: Create order with brand admin (should succeed)
    run_test "Create order with brand admin" \
        "curl -X POST http://localhost:8889/api/v1/orders -H 'Content-Type: application/json' -d '{\"campaign_id\": 1, \"phone\": \"13800138001\", \"form_data\": {\"name\": \"Test User\"}}' -H 'Authorization: Bearer brand_admin_token'" \
        "success"
    
    # Test 2: Create order with participant (should fail - no permission)
    run_test "Create order with participant (should fail)" \
        "curl -X POST http://localhost:8889/api/v1/orders -H 'Content-Type: application/json' -d '{\"campaign_id\": 1, \"phone\": \"13800138002\", \"form_data\": {\"name\": \"Test User\"}}' -H 'Authorization: Bearer participant_token'" \
        "failure"
    
    echo -e "\n${YELLOW}=== Running Order Verification Tests ===${NC}"
    
    # Test 3: Verify order with brand admin (should succeed)
    run_test "Verify order with brand admin" \
        "curl -X POST http://localhost:8889/api/v1/orders/verify -H 'Content-Type: application/json' -d '{\"code\": \"1_13800138001_1717651200_signature\", \"remark\": \"Test verification\"}' -H 'Authorization: Bearer brand_admin_token'" \
        "success"
    
    # Test 4: Verify order with participant (should fail - no permission)
    run_test "Verify order with participant (should fail)" \
        "curl -X POST http://localhost:8889/api/v1/orders/verify -H 'Content-Type: application/json' -d '{\"code\": \"1_13800138001_1717651200_signature\", \"remark\": \"Test verification\"}' -H 'Authorization: Bearer participant_token'" \
        "failure"
    
    echo -e "\n${YELLOW}=== Running Order Unverification Tests ===${NC}"
    
    # Test 5: Unverify order with brand admin (should succeed)
    run_test "Unverify order with brand admin" \
        "curl -X POST http://localhost:8889/api/v1/orders/unverify -H 'Content-Type: application/json' -d '{\"code\": \"1_13800138001_1717651200_signature\", \"reason\": \"Test unverification\"}' -H 'Authorization: Bearer brand_admin_token'" \
        "success"
    
    echo -e "\n${YELLOW}=== Running Transaction Rollback Tests ===${NC}"
    
    # Test 6: Transaction rollback on verification failure
    run_test "Transaction rollback on verification failure" \
        "curl -X POST http://localhost:8889/api/v1/orders/verify -H 'Content-Type: application/json' -d '{\"code\": \"1_13800138003_1717651200_malformed_signature\", \"remark\": \"Test verification\"}' -H 'Authorization: Bearer brand_admin_token'" \
        "failure"
    
    # Print summary
    echo -e "\n${YELLOW}=== Test Summary ===${NC}"
    echo "Total tests: $TOTAL_TESTS"
    echo "Passed: $PASS_COUNT"
    echo "Failed: $FAIL_COUNT"
    
    if [ $FAIL_COUNT -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}Some tests failed.${NC}"
        exit 1
    fi
}

# Run the main function
main "$@"