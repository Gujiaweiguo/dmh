#!/bin/bash
# DMH 生产环境目录初始化脚本

set -e

DATA_DIR="/opt/data"
LOGS_DIR="/opt/logs"

echo "=== DMH 生产环境目录初始化 ==="
echo ""

# 创建数据目录
echo "创建数据目录..."
mkdir -p "$DATA_DIR/mysql"
mkdir -p "$DATA_DIR/redis/dmh"
mkdir -p "$DATA_DIR/uploads/dmh"

# 创建日志目录
echo "创建日志目录..."
mkdir -p "$LOGS_DIR/mysql"
mkdir -p "$LOGS_DIR/dmh-api"
mkdir -p "$LOGS_DIR/nginx"

# 设置权限
if command -v chown &> /dev/null; then
    chown -R 999:999 "$DATA_DIR/mysql" 2>/dev/null || true
    chown -R 999:999 "$LOGS_DIR/mysql" 2>/dev/null || true
fi
chmod -R 755 "$DATA_DIR/mysql" 2>/dev/null || true
chmod -R 755 "$LOGS_DIR/mysql" 2>/dev/null || true
chmod -R 755 "$DATA_DIR/redis/dmh" 2>/dev/null || true
chmod -R 777 "$DATA_DIR/uploads/dmh" 2>/dev/null || true
chmod -R 777 "$LOGS_DIR/dmh-api" 2>/dev/null || true
chmod -R 777 "$LOGS_DIR/nginx" 2>/dev/null || true

echo ""
echo "=== 目录结构 ==="
echo ""
echo "数据目录: $DATA_DIR"
ls -la "$DATA_DIR"
echo ""
echo "日志目录: $LOGS_DIR"
ls -la "$LOGS_DIR"
echo ""
echo "=== 初始化完成 ==="
