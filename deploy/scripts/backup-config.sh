#!/bin/bash
# backup-config.sh - 备份当前配置
# 用法: ./backup-config.sh [--list]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
CONFIG_DIR="/opt/module/dmh/configs"
BACKUP_DIR="$CONFIG_DIR/backup"
MAX_BACKUPS=10

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --list      列出所有备份"
    echo "  --restore   恢复最近的备份"
    echo "  --help      显示此帮助"
    echo ""
    echo "不带参数运行将创建新备份"
}

# 列出备份
list_backups() {
    echo -e "${BLUE}现有备份:${NC}"
    echo ""
    
    if [[ ! -d "$BACKUP_DIR" ]] || [[ -z $(ls -A "$BACKUP_DIR" 2>/dev/null) ]]; then
        echo "  暂无备份"
        return
    fi
    
    local count=0
    for backup in $(ls -1dt "$BACKUP_DIR"/*/ 2>/dev/null); do
        ((count++))
        local name=$(basename "$backup")
        local size=$(du -sh "$backup" 2>/dev/null | cut -f1)
        local files=$(find "$backup" -type f | wc -l)
        echo -e "  $count. ${GREEN}$name${NC} ($size, $files 文件)"
    done
    
    echo ""
    echo "共 $count 个备份"
}

# 恢复备份
restore_backup() {
    local latest=$(ls -1dt "$BACKUP_DIR"/*/ 2>/dev/null | head -1)
    
    if [[ -z "$latest" ]]; then
        echo -e "${RED}错误: 没有可恢复的备份${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}将从以下备份恢复: $(basename $latest)${NC}"
    read -p "确认恢复? (y/N): " confirm
    
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        echo "已取消"
        exit 0
    fi
    
    # 恢复配置
    cp -r "$latest"* "$CONFIG_DIR/"
    echo -e "${GREEN}✓ 配置已恢复${NC}"
}

# 参数处理
case "$1" in
    --help|-h)
        show_help
        exit 0
        ;;
    --list|-l)
        list_backups
        exit 0
        ;;
    --restore|-r)
        restore_backup
        exit 0
        ;;
esac

# 创建备份
echo "=========================================="
echo "  DMH 配置备份工具"
echo "=========================================="
echo ""

# 创建带时间戳的备份目录
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_PATH="$BACKUP_DIR/$TIMESTAMP"

mkdir -p "$BACKUP_PATH"

# 备份配置文件
backup_count=0

# 备份 dmh-api.yaml
if [[ -f "$CONFIG_DIR/dmh-api.yaml" ]]; then
    cp "$CONFIG_DIR/dmh-api.yaml" "$BACKUP_PATH/"
    echo -e "${GREEN}✓${NC} dmh-api.yaml"
    ((backup_count++))
fi

# 备份 nginx 配置
if [[ -d "$CONFIG_DIR/nginx" ]]; then
    mkdir -p "$BACKUP_PATH/nginx"
    cp -r "$CONFIG_DIR/nginx/"* "$BACKUP_PATH/nginx/" 2>/dev/null || true
    echo -e "${GREEN}✓${NC} nginx/"
    ((backup_count++))
fi

# 备份前端配置
if [[ -d "$CONFIG_DIR/frontend" ]]; then
    mkdir -p "$BACKUP_PATH/frontend"
    cp -r "$CONFIG_DIR/frontend/"* "$BACKUP_PATH/frontend/" 2>/dev/null || true
    echo -e "${GREEN}✓${NC} frontend/"
    ((backup_count++))
fi

# 清理旧备份
backup_total=$(ls -1d "$BACKUP_DIR"/*/ 2>/dev/null | wc -l)
if [[ $backup_total -gt $MAX_BACKUPS ]]; then
    echo ""
    echo -e "${YELLOW}清理旧备份 (保留最近 $MAX_BACKUPS 个)...${NC}"
    ls -1dt "$BACKUP_DIR"/*/ | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -rf
    echo -e "${GREEN}✓ 已清理 $((backup_total - MAX_BACKUPS)) 个旧备份${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}备份完成!${NC}"
echo "  位置: $BACKUP_PATH"
echo "  文件: $backup_count 项"
echo "=========================================="
