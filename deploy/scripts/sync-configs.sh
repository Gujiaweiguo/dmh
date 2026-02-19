#!/bin/bash
# sync-configs.sh - 从项目目录同步配置到统一管理目录
# 用法: ./sync-configs.sh [--dry-run]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
PROJECT_ROOT="/opt/code/dmh"
CONFIG_DIR="/opt/module/dmh/configs"
DRY_RUN=false

# 解析参数
if [[ "$1" == "--dry-run" ]]; then
    DRY_RUN=true
    echo -e "${YELLOW}[DRY-RUN] 仅显示将要执行的操作${NC}"
fi

# 源文件和目标映射
declare -A SYNC_MAP=(
    ["$PROJECT_ROOT/deploy/dmh-api.yaml"]="$CONFIG_DIR/dmh-api.yaml"
    ["$PROJECT_ROOT/deploy/nginx/conf.d/default.conf"]="$CONFIG_DIR/nginx/conf.d/default.conf"
)

echo "=========================================="
echo "  DMH 配置同步工具"
echo "=========================================="
echo ""

# 确保目标目录存在
mkdir -p "$CONFIG_DIR/nginx/conf.d"
mkdir -p "$CONFIG_DIR/frontend"
mkdir -p "$CONFIG_DIR/backup"

# 同步计数
SYNC_COUNT=0
SKIP_COUNT=0

for src in "${!SYNC_MAP[@]}"; do
    dest="${SYNC_MAP[$src]}"
    
    if [[ ! -f "$src" ]]; then
        echo -e "${RED}✗ 源文件不存在: $src${NC}"
        continue
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${YELLOW}  将复制: $src${NC}"
        echo -e "${YELLOW}       → $dest${NC}"
        ((SYNC_COUNT++))
        continue
    fi
    
    # 执行复制
    cp "$src" "$dest"
    chmod 644 "$dest"
    echo -e "${GREEN}✓ 已同步: $(basename $src)${NC}"
    ((SYNC_COUNT++))
done

echo ""
echo "=========================================="
echo -e "${GREEN}同步完成: $SYNC_COUNT 个文件${NC}"
echo "=========================================="
