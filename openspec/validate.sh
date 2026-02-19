#!/bin/bash

# OpenSpec Validation Script
# This script validates the OpenSpec changes for order and RBAC tests

echo "=== OpenSpec Validation ==="
echo "Validating order and RBAC test changes..."

# Check if required files exist
REQUIRED_FILES=(
    "/opt/code/dmh/backend/api/internal/logic/order/order_logic_test.go"
    "/opt/code/dmh/backend/api/internal/middleware/permission_middleware_test.go"
    "/opt/code/dmh/backend/api/internal/logic/order/e2e_regression_script.sh"
    "/opt/code/dmh/backend/api/internal/logic/order/e2e_results_template.md"
)

MISSING_FILES=0

for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "✓ $file exists"
    else
        echo "✗ $file is missing"
        MISSING_FILES=$((MISSING_FILES + 1))
    fi
done

# Check test results
echo -e "\n=== Test Results Validation ==="
TEST_RESULTS_DIR="/opt/code/dmh/backend/api/internal/logic/order/test_results"
if [ -d "$TEST_RESULTS_DIR" ]; then
    echo "✓ Test results directory exists"
    RESULT_FILES=$(find "$TEST_RESULTS_DIR" -name "*.md" -o -name "*.txt" | wc -l)
    echo "Found $RESULT_FILES result files"
else
    echo "✗ Test results directory is missing"
    MISSING_FILES=$((MISSING_FILES + 1))
fi

# Summary
echo -e "\n=== Validation Summary ==="
if [ $MISSING_FILES -eq 0 ]; then
    echo "✓ All validation checks passed"
    exit 0
else
    echo "✗ $MISSING_FILES validation checks failed"
    exit 1
fi