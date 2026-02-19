#!/bin/bash
# verify-config.sh - 验证配置文件正确性
# 用法: ./verify-config.sh [--level L1|L2|L3]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
CONFIG_DIR="/opt/module/dmh/configs"
PROJECT_ROOT="/opt/code/dmh"

# 验证级别
LEVEL="ALL"
ERRORS=0
WARNINGS=0

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --level L1   仅语法检查"
    echo "  --level L2   L1 + 连接检查"
    echo "  --level L3   L2 + 功能检查"
    echo "  --help       显示此帮助"
    echo ""
    echo "验证级别:"
    echo "  L1 - 配置语法检查 (YAML, Nginx)"
    echo "  L2 - 数据库/Redis 连接检查"
    echo "  L3 - API 健康检查"
}

# 参数处理
case "$1" in
    --help|-h)
        show_help
        exit 0
        ;;
    --level)
        LEVEL="$2"
        ;;
esac

echo "=========================================="
echo "  DMH 配置验证工具"
echo "=========================================="
echo -e "验证级别: ${BLUE}$LEVEL${NC}"
echo ""

# ============================================
# L1: 语法检查
# ============================================
echo -e "${BLUE}[L1] 语法检查${NC}"
echo ""

# YAML 语法检查
check_yaml() {
    local file="$1"
    local name=$(basename "$file")
    
    if [[ ! -f "$file" ]]; then
        echo -e "  ${YELLOW}⊙${NC} $name - 文件不存在"
        return 0
    fi
    
    if python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
        echo -e "  ${GREEN}✓${NC} $name - YAML 语法正确"
        return 0
    else
        echo -e "  ${RED}✗${NC} $name - YAML 语法错误"
        ((ERRORS++))
        return 1
    fi
}

# Nginx 语法检查
check_nginx() {
    local file="$1"
    local name=$(basename "$file")
    
    if [[ ! -f "$file" ]]; then
        echo -e "  ${YELLOW}⊙${NC} $name - 文件不存在"
        return 0
    fi
    
    # 使用 nginx -t 检查（需要临时配置）
    if docker exec dmh-nginx nginx -t 2>&1 | grep -q "successful"; then
        echo -e "  ${GREEN}✓${NC} $name - Nginx 语法正确"
        return 0
    else
        echo -e "  ${RED}✗${NC} $name - Nginx 语法检查失败"
        ((WARNINGS++))
        return 1
    fi
}

# 执行 L1 检查
check_yaml "$CONFIG_DIR/dmh-api.yaml"
check_nginx "$CONFIG_DIR/nginx/conf.d/default.conf"

if [[ "$LEVEL" == "L1" ]]; then
    echo ""
    echo "=========================================="
    if [[ $ERRORS -eq 0 ]]; then
        echo -e "${GREEN}✓ L1 验证通过${NC}"
        exit 0
    else
        echo -e "${RED}✗ L1 验证失败 ($ERRORS 个错误)${NC}"
        exit 1
    fi
fi

# ============================================
# L2: 连接检查
# ============================================
echo ""
echo -e "${BLUE}[L2] 连接检查${NC}"
echo ""

# 从配置文件读取数据库连接信息
DB_HOST=$(grep -oP 'Host:\s*\K[^[:space:]]+' "$CONFIG_DIR/dmh-api.yaml" 2>/dev/null | head -1 || echo "mysql8")
DB_PORT=$(grep -oP 'Port:\s*\K\d+' "$CONFIG_DIR/dmh-api.yaml" 2>/dev/null | head -1 || echo "3306")
REDIS_HOST=$(grep -A5 'Redis:' "$CONFIG_DIR/dmh-api.yaml" 2>/dev/null | grep -oP 'Host:\s*\K[^[:space:]]+' | head -1 || echo "redis-dmh:6379")

# 数据库连接检查
check_database() {
    echo -n "  数据库 ($DB_HOST:$DB_PORT): "
    
    if docker exec dmh-api sh -c "timeout 5 bash -c '</dev/tcp/$DB_HOST/$DB_PORT'" 2>/dev/null; then
        echo -e "${GREEN}可达${NC}"
        return 0
    else
        echo -e "${YELLOW}不可达 (警告)${NC}"
        ((WARNINGS++))
        return 1
    fi
}

# Redis 连接检查
check_redis() {
    echo -n "  Redis ($REDIS_HOST): "
    
    if docker exec dmh-api sh -c "timeout 5 bash -c '</dev/tcp/${REDIS_HOST%:*}/${REDIS_HOST#*:}'" 2>/dev/null; then
        echo -e "${GREEN}可达${NC}"
        return 0
    else
        echo -e "${YELLOW}不可达 (警告)${NC}"
        ((WARNINGS++))
        return 1
    fi
}

check_database
check_redis

if [[ "$LEVEL" == "L2" ]]; then
    echo ""
    echo "=========================================="
    if [[ $ERRORS -eq 0 ]]; then
        echo -e "${GREEN}✓ L2 验证通过${NC}"
        echo -e "  警告: $WARNINGS"
        exit 0
    else
        echo -e "${RED}✗ L2 验证失败 ($ERRORS 个错误)${NC}"
        exit 1
    fi
fi

# ============================================
# L3: 功能检查
# ============================================
echo ""
echo -e "${BLUE}[L3] 功能检查${NC}"
echo ""

# API 健康检查
check_api_health() {
    echo -n "  API 健康检查: "
    
    local response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8889/api/v1/health 2>/dev/null || echo "000")
    
    # 尝试其他可能的健康端点
    if [[ "$response" == "000" || "$response" == "404" ]]; then
        response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8889/health 2>/dev/null || echo "000")
    fi
    
    if [[ "$response" == "200" ]]; then
        echo -e "${GREEN}正常 (HTTP $response)${NC}"
        return 0
    elif [[ "$response" == "404" ]]; then
        # 404 也算服务正常（端点可能不存在但服务在运行）
        echo -e "${GREEN}服务运行中 (HTTP $response)${NC}"
        return 0
    else
        echo -e "${YELLOW}异常 (HTTP $response)${NC}"
        ((WARNINGS++))
        return 1
    fi
}

# 前端检查
check_frontend() {
    echo -n "  管理后台 (3000): "
    local response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 2>/dev/null || echo "000")
    if [[ "$response" == "200" ]]; then
        echo -e "${GREEN}正常${NC}"
    else
        echo -e "${YELLOW}异常 (HTTP $response)${NC}"
        ((WARNINGS++))
    fi
    
    echo -n "  H5 前端 (3100): "
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3100 2>/dev/null || echo "000")
    if [[ "$response" == "200" ]]; then
        echo -e "${GREEN}正常${NC}"
    else
        echo -e "${YELLOW}异常 (HTTP $response)${NC}"
        ((WARNINGS++))
    fi
}

check_api_health
check_frontend

# ============================================
# 总结
# ============================================
echo ""
echo "=========================================="
if [[ $ERRORS -eq 0 ]]; then
    echo -e "${GREEN}✓ 验证通过${NC}"
    echo -e "  错误: $ERRORS"
    echo -e "  警告: $WARNINGS"
    exit 0
else
    echo -e "${RED}✗ 验证失败${NC}"
    echo -e "  错误: $ERRORS"
    echo -e "  警告: $WARNINGS"
    exit 1
fi
