#!/bin/bash
# ============================================
# DMH 测试数据恢复脚本
# ============================================

set -e

# 配置
MYSQL_CONTAINER="mysql8"
MYSQL_USER="root"
MYSQL_PASSWORD="Admin168"
DATABASE="dmh"
BACKUP_DIR="/tmp"
TEST_DATA_SQL="/opt/code/DMH/backend/scripts/dmh_test_data_20260131_final.sql"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查MySQL容器是否运行
if ! docker ps | grep -q $MYSQL_CONTAINER; then
    log_error "MySQL容器 $MYSQL_CONTAINER 未运行"
    exit 1
fi

log_info "MySQL容器 $MYSQL_CONTAINER 运行正常"

# 检查测试数据文件是否存在
if [ ! -f "$TEST_DATA_SQL" ]; then
    log_error "测试数据文件不存在: $TEST_DATA_SQL"
    exit 1
fi

log_info "测试数据文件存在: $TEST_DATA_SQL"

# 备份当前数据库
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/dmh_backup_before_restore_${TIMESTAMP}.sql"

log_info "备份当前数据库到: $BACKUP_FILE"
docker exec $MYSQL_CONTAINER mysqldump -u$MYSQL_USER -p"$MYSQL_PASSWORD" --default-character-set=utf8mb4 $DATABASE > $BACKUP_FILE 2>&1 | grep -v "Using a password"

if [ -f "$BACKUP_FILE" ]; then
    log_info "数据库备份完成: $BACKUP_FILE"
else
    log_error "数据库备份失败"
    exit 1
fi

# 重建数据库
log_info "重建数据库 $DATABASE"
docker exec $MYSQL_CONTAINER mysql -u$MYSQL_USER -p"$MYSQL_PASSWORD" -e "DROP DATABASE IF EXISTS $DATABASE; CREATE DATABASE $DATABASE DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>&1 | grep -v "Using a password"

# 导入测试数据
log_info "导入测试数据: $TEST_DATA_SQL"
docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p"$MYSQL_PASSWORD" --default-character-set=utf8mb4 $DATABASE < $TEST_DATA_SQL 2>&1 | grep -v "Using a password"

# 验证数据
log_info "验证数据..."
docker exec $MYSQL_CONTAINER mysql -u$MYSQL_USER -p"$MYSQL_PASSWORD" -D $DATABASE -e "
SELECT 
    (SELECT COUNT(*) FROM users) as users,
    (SELECT COUNT(*) FROM brands) as brands,
    (SELECT COUNT(*) FROM campaigns) as campaigns,
    (SELECT COUNT(*) FROM orders) as orders,
    (SELECT COUNT(*) FROM rewards) as rewards,
    (SELECT COUNT(*) FROM withdrawals) as withdrawals,
    (SELECT COUNT(*) FROM audit_logs) as audit_logs;
" 2>&1 | grep -v "Using a password"

log_info "测试数据恢复完成！"
log_info "备份文件位置: $BACKUP_FILE"

# 提示信息
echo ""
echo "========================================="
echo "  测试账号信息"
echo "========================================="
echo ""
echo "平台管理员:"
echo "  用户名: admin"
echo "  密码: 123456"
echo ""
echo "品牌管理员:"
echo "  用户名: brand_manager / brand_admin"
echo "  密码: 123456"
echo ""
echo "普通用户:"
echo "  用户名: user001 / user002 / user003"
echo "  密码: 123456"
echo ""
echo "分销商:"
echo "  用户名: distributor001 / distributor002"
echo "  密码: 123456"
echo ""
echo "========================================="
echo "  访问地址"
echo "========================================="
echo ""
echo "管理后台: http://localhost:3000"
echo "H5前端: http://localhost:3100"
echo "品牌管理: http://localhost:3100/brand/login"
echo "分销中心: http://localhost:3100/distributor"
echo ""
