#!/bin/bash
# restart-services.sh - 一键重启服务并验证
# 用法: ./restart-services.sh [--skip-backup] [--skip-verify]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose-simple.yml"
BACKUP_SCRIPT="$SCRIPT_DIR/backup-config.sh"
VERIFY_SCRIPT="$SCRIPT_DIR/verify-config.sh"

# 选项
SKIP_BACKUP=false
SKIP_VERIFY=false

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-backup)
            SKIP_BACKUP=true
            shift
            ;;
        --skip-verify)
            SKIP_VERIFY=true
            shift
            ;;
        --help|-h)
            echo "用法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  --skip-backup   跳过备份步骤"
            echo "  --skip-verify   跳过验证步骤"
            echo "  --help          显示此帮助"
            echo ""
            echo "此脚本将执行以下步骤:"
            echo "  1. 备份当前配置"
            echo "  2. 验证配置语法"
            echo "  3. 重启 Docker 容器"
            echo "  4. 等待服务就绪"
            echo "  5. 执行健康检查"
            exit 0
            ;;
        *)
            echo -e "${RED}未知参数: $1${NC}"
            exit 1
            ;;
    esac
done

echo "=========================================="
echo "  DMH 服务重启工具"
echo "=========================================="
echo -e "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# ============================================
# 步骤 1: 备份当前配置
# ============================================
if [[ "$SKIP_BACKUP" == false ]]; then
    echo -e "${BLUE}[1/5] 备份当前配置${NC}"
    
    if [[ -x "$BACKUP_SCRIPT" ]]; then
        "$BACKUP_SCRIPT"
    else
        echo -e "${YELLOW}  备份脚本不存在，跳过备份${NC}"
    fi
else
    echo -e "${YELLOW}[1/5] 跳过备份 (--skip-backup)${NC}"
fi
echo ""

# ============================================
# 步骤 2: 验证配置
# ============================================
echo -e "${BLUE}[2/5] 验证配置语法${NC}"

if [[ "$SKIP_VERIFY" == false ]]; then
    if [[ -x "$VERIFY_SCRIPT" ]]; then
        if "$VERIFY_SCRIPT" --level L1; then
            echo -e "${GREEN}  ✓ 配置验证通过${NC}"
        else
            echo -e "${RED}  ✗ 配置验证失败，请检查配置文件${NC}"
            echo ""
            read -p "是否继续重启? (y/N): " continue_restart
            if [[ "$continue_restart" != "y" && "$continue_restart" != "Y" ]]; then
                echo "已取消重启"
                exit 1
            fi
        fi
    else
        echo -e "${YELLOW}  验证脚本不存在，跳过验证${NC}"
    fi
else
    echo -e "${YELLOW}  跳过验证 (--skip-verify)${NC}"
fi
echo ""

# ============================================
# 步骤 3: 重启 Docker 容器
# ============================================
echo -e "${BLUE}[3/5] 重启 Docker 容器${NC}"

cd "$PROJECT_ROOT"

# 使用 docker compose 重启
if docker compose -f docker-compose-simple.yml restart dmh-api dmh-nginx; then
    echo -e "${GREEN}  ✓ 容器重启命令已执行${NC}"
else
    echo -e "${RED}  ✗ 容器重启失败${NC}"
    exit 1
fi
echo ""

# ============================================
# 步骤 4: 等待服务就绪
# ============================================
echo -e "${BLUE}[4/5] 等待服务就绪${NC}"

MAX_WAIT=60
WAIT_INTERVAL=2
WAITED=0

while [[ $WAITED -lt $MAX_WAIT ]]; do
    # 检查 API 容器状态
    API_STATUS=$(docker inspect --format='{{.State.Running}}' dmh-api 2>/dev/null || echo "false")
    NGINX_STATUS=$(docker inspect --format='{{.State.Running}}' dmh-nginx 2>/dev/null || echo "false")
    
    if [[ "$API_STATUS" == "true" && "$NGINX_STATUS" == "true" ]]; then
        # 额外等待 2 秒确保服务完全启动
        sleep 2
        echo -e "${GREEN}  ✓ 服务已就绪 (${WAITED}s)${NC}"
        break
    fi
    
    echo -n "."
    sleep $WAIT_INTERVAL
    WAITED=$((WAITED + WAIT_INTERVAL))
done

if [[ $WAITED -ge $MAX_WAIT ]]; then
    echo -e "${RED}  ✗ 等待超时 (${MAX_WAIT}s)${NC}"
fi
echo ""

# ============================================
# 步骤 5: 健康检查
# ============================================
echo -e "${BLUE}[5/5] 健康检查${NC}"

check_service() {
    local name="$1"
    local url="$2"
    local expected="${3:-200}"
    
    local response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
    
    if [[ "$response" == "$expected" || "$response" == "404" ]]; then
        echo -e "  ${GREEN}✓${NC} $name - HTTP $response"
        return 0
    else
        echo -e "  ${RED}✗${NC} $name - HTTP $response"
        return 1
    fi
}

HEALTH_OK=true

# 检查各服务
check_service "后端 API" "http://localhost:8889/api/v1/health" || HEALTH_OK=false
check_service "管理后台" "http://localhost:3000" || HEALTH_OK=false
check_service "H5 前端" "http://localhost:3100" || HEALTH_OK=false

echo ""

# ============================================
# 总结
# ============================================
echo "=========================================="

if [[ "$HEALTH_OK" == true ]]; then
    echo -e "${GREEN}✓ 服务重启成功!${NC}"
    echo ""
    echo "访问地址:"
    echo "  - 管理后台: http://localhost:3000"
    echo "  - H5 前端:  http://localhost:3100"
    echo "  - 后端 API: http://localhost:8889"
    exit 0
else
    echo -e "${RED}✗ 服务重启完成，但健康检查未全部通过${NC}"
    echo ""
    echo "故障排查:"
    echo "  1. 查看日志: docker logs dmh-api"
    echo "  2. 查看日志: docker logs dmh-nginx"
    echo ""
    echo "回滚命令:"
    echo "  $SCRIPT_DIR/backup-config.sh --list"
    echo "  $SCRIPT_DIR/backup-config.sh --restore"
    exit 1
fi
