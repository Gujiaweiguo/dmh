#!/bin/bash

# RBAC权限系统测试运行脚本

set -e

echo "=========================================="
echo "RBAC权限系统测试套件"
echo "=========================================="

# 设置测试环境变量
export GO_ENV=test
export CGO_ENABLED=1

# 进入后端目录
cd "$(dirname "$0")/.."

echo "1. 运行单元测试..."
echo "------------------------------------------"

# 运行安全服务单元测试
echo "运行密码服务测试..."
go test -v ./api/internal/service -run TestPasswordServiceTestSuite -timeout 30s

echo "运行审计服务测试..."
go test -v ./api/internal/service -run TestAuditServiceTestSuite -timeout 30s

echo "运行会话服务测试..."
go test -v ./api/internal/service -run TestSessionServiceTestSuite -timeout 30s

# 运行中间件测试
echo "运行权限中间件测试..."
go test -v ./api/internal/middleware -run TestPermissionMiddlewareTestSuite -timeout 30s

# 运行逻辑层测试
echo "运行登录逻辑测试..."
go test -v ./api/internal/logic/auth -run TestLoginLogicTestSuite -timeout 30s

echo ""
echo "2. 运行集成测试..."
echo "------------------------------------------"

# 运行RBAC集成测试
echo "运行RBAC集成测试..."
go test -v ./test/integration -run TestRBACIntegrationTestSuite -timeout 60s

echo ""
echo "3. 运行性能测试..."
echo "------------------------------------------"

# 运行性能测试
echo "运行RBAC性能测试..."
go test -v ./test/performance -run TestRBACPerformanceTestSuite -timeout 300s

echo ""
echo "4. 运行基准测试..."
echo "------------------------------------------"

# 运行基准测试
echo "运行权限检查基准测试..."
go test -bench=BenchmarkPermissionCheck ./test/performance -benchmem -count=3

echo "运行密码强度检查基准测试..."
go test -bench=BenchmarkPasswordStrengthCheck ./test/performance -benchmem -count=3

echo "运行会话验证基准测试..."
go test -bench=BenchmarkSessionValidation ./test/performance -benchmem -count=3

echo ""
echo "5. 生成测试覆盖率报告..."
echo "------------------------------------------"

# 生成覆盖率报告
echo "生成测试覆盖率报告..."
go test -coverprofile=coverage.out ./api/internal/service ./api/internal/middleware ./api/internal/logic/auth
go tool cover -html=coverage.out -o coverage.html

echo "覆盖率报告已生成: coverage.html"

echo ""
echo "6. 运行竞态条件检测..."
echo "------------------------------------------"

# 运行竞态条件检测
echo "检测竞态条件..."
go test -race -v ./api/internal/service -run TestPasswordServiceTestSuite -timeout 30s
go test -race -v ./api/internal/service -run TestAuditServiceTestSuite -timeout 30s
go test -race -v ./api/internal/service -run TestSessionServiceTestSuite -timeout 30s

echo ""
echo "=========================================="
echo "测试完成！"
echo "=========================================="

# 检查测试结果
if [ $? -eq 0 ]; then
    echo "✅ 所有测试通过"
    exit 0
else
    echo "❌ 部分测试失败"
    exit 1
fi