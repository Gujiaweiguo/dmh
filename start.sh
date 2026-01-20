#!/bin/bash

# ============================================
# DMH 数字营销中台 - 启动脚本
# ============================================

set -e

echo "========================================"
echo "DMH 数字营销中台启动脚本"
echo "========================================"

# 项目根目录
PROJECT_ROOT="/opt/code/DMH"
cd "$PROJECT_ROOT"

# 检查 MySQL 是否运行
echo ""
echo "检查 MySQL 服务..."
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL 未安装，请先安装 MySQL"
    exit 1
fi

# 数据库配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-123456}"
DB_NAME="${DB_NAME:-dmh}"

# ============================================
# 1. 初始化数据库
# ============================================
echo ""
echo "========================================"
echo "步骤 1/4: 初始化数据库"
echo "========================================"

# 创建数据库并执行初始化脚本
echo "执行数据库初始化脚本..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$PROJECT_ROOT/backend/scripts/init.sql"

# 执行分销商表迁移
echo "执行分销商表迁移..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$PROJECT_ROOT/backend/migrations/20250120_create_distributor_tables_final.sql"

# 执行完整测试数据
echo "导入完整测试数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$PROJECT_ROOT/backend/scripts/seed_complete_test_data.sql"

# 执行分销商测试数据
echo "导入分销商系统测试数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$PROJECT_ROOT/backend/scripts/seed_distributor_test_data.sql"

echo "✓ 数据库初始化完成"

# ============================================
# 2. 启动后端 API
# ============================================
echo ""
echo "========================================"
echo "步骤 2/4: 启动后端 API"
echo "========================================"

cd "$PROJECT_ROOT/backend/api"

# 检查是否已安装依赖
if [ ! -d "vendor" ]; then
    echo "安装 Go 依赖..."
    go mod download
fi

echo "启动后端 API 服务..."
echo "API 服务地址: http://localhost:8888"
go run dmh.go > /tmp/dmh-backend.log 2>&1 &
BACKEND_PID=$!
echo "后端 PID: $BACKEND_PID"

# 等待后端启动
echo "等待后端启动..."
sleep 5

# 检查后端是否启动成功
if ps -p $BACKEND_PID > /dev/null; then
    echo "✓ 后端 API 启动成功"
else
    echo "❌ 后端 API 启动失败，查看日志: /tmp/dmh-backend.log"
    cat /tmp/dmh-backend.log
    exit 1
fi

# ============================================
# 3. 启动前端管理后台
# ============================================
echo ""
echo "========================================"
echo "步骤 3/4: 启动前端管理后台"
echo "========================================"

cd "$PROJECT_ROOT/frontend-admin"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

echo "启动前端管理后台服务..."
echo "管理后台地址: http://localhost:5173"
npm run dev > /tmp/dmh-admin.log 2>&1 &
ADMIN_PID=$!
echo "前端管理后台 PID: $ADMIN_PID"

# 等待前端启动
echo "等待前端管理后台启动..."
sleep 5

# ============================================
# 4. 启动 H5 端
# ============================================
echo ""
echo "========================================"
echo "步骤 4/4: 启动 H5 端"
echo "========================================"

cd "$PROJECT_ROOT/frontend-h5"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "安装 H5 依赖..."
    npm install
fi

echo "启动 H5 端服务..."
echo "H5 端地址: http://localhost:3000"
npm start > /tmp/dmh-h5.log 2>&1 &
H5_PID=$!
echo "H5 端 PID: $H5_PID"

# 等待 H5 启动
echo "等待 H5 端启动..."
sleep 3

# ============================================
# 输出系统信息
# ============================================
echo ""
echo "========================================"
echo "✓ 系统启动完成！"
echo "========================================"
echo ""
echo "系统访问地址："
echo "----------------------------------------"
echo "后端 API：     http://localhost:8888"
echo "前端管理后台：  http://localhost:5173"
echo "H5 端：       http://localhost:3000"
echo ""
echo "========================================"
echo "测试账号信息："
echo "========================================"
echo ""
echo "【平台管理员】"
echo "用户名: admin"
echo "密码:   123456"
echo "权限:   系统所有权限"
echo ""
echo "【品牌管理员】"
echo "用户名: brand_manager"
echo "密码:   123456"
echo "权限:   管理品牌A的业务"
echo ""
echo "【一级分销商】"
echo "用户名: distributor1"
echo "密码:   123456"
echo "余额:   ¥500.00"
echo "下级:   2人"
echo ""
echo "【二级分销商】"
echo "用户名: distributor2"
echo "密码:   123456"
echo "余额:   ¥200.00"
echo "下级:   1人"
echo ""
echo "【三级分销商】"
echo "用户名: distributor3"
echo "密码:   123456"
echo "余额:   ¥100.00"
echo "下级:   0人"
echo ""
echo "【普通用户1】"
echo "用户名: participant1"
echo "密码:   123456"
echo "余额:   ¥0.00"
echo ""
echo "【普通用户2】"
echo "用户名: participant2"
echo "密码:   123456"
echo "余额:   ¥0.00"
echo ""
echo "========================================"
echo "分销链结构："
echo "========================================"
echo "一级: distributor1 (推荐人ID: 无)"
echo "  ├─ 二级: distributor2 (推荐人ID: distributor1)"
echo "  │    └─ 三级: distributor3 (推荐人ID: distributor2)"
echo "  └─ 待添加..."
echo ""
echo "========================================"
echo "提现申请状态："
echo "========================================"
echo "distributor1:"
echo "  - 待审核: ¥50.00 (微信)"
echo "  - 已通过: ¥100.00 (微信)"
echo ""
echo "distributor2:"
echo "  - 待审核: ¥50.00 (支付宝)"
echo "  - 已拒绝: ¥300.00 (金额超限)"
echo ""
echo "========================================"
echo "订单数据："
echo "========================================"
echo "订单10: ¥100.00 (通过distributor1推荐，一级分销)"
echo "  → distributor1 获得 ¥10.00"
echo ""
echo "订单11: ¥200.00 (通过distributor3推荐，三级分销)"
echo "  → distributor1 获得 ¥20.00"
echo "  → distributor2 获得 ¥10.00"
echo "  → distributor3 获得 ¥6.00"
echo ""
echo "订单12: ¥150.00 (通过distributor2推荐，二级分销)"
echo "  → distributor1 获得 ¥15.00"
echo "  → distributor2 获得 ¥7.50"
echo ""
echo "订单13: ¥80.00 (无推荐，无分销奖励)"
echo ""
echo "========================================"
echo "进程信息："
echo "========================================"
echo "后端 PID: $BACKEND_PID"
echo "前端管理后台 PID: $ADMIN_PID"
echo "H5 端 PID: $H5_PID"
echo ""
echo "查看日志："
echo "  后端:     tail -f /tmp/dmh-backend.log"
echo "  前端管理后台: tail -f /tmp/dmh-admin.log"
echo "  H5 端:    tail -f /tmp/dmh-h5.log"
echo ""
echo "停止系统："
echo "  kill $BACKEND_PID $ADMIN_PID $H5_PID"
echo ""
echo "========================================"
echo ""
