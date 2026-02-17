#!/bin/bash

# ============================================
# DMH 活动高级功能部署脚本
# ============================================

set -e

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BACKEND_DIR="$PROJECT_ROOT/backend"
DOCKER_COMPOSE="$PROJECT_ROOT/deployment/docker-compose-dmh.yml"
LOG_FILE="$PROJECT_ROOT/deployment/logs/deploy-$(date +%Y%m%d_%H%M%S).log"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} ERROR: $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} WARNING: $1" | tee -a "$LOG_FILE"
}

# 检查前置条件
check_prerequisites() {
    log "检查部署前置条件..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查 docker-compose
    if ! command -v docker-compose &> /dev/null && ! command -v docker &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go 1.23+"
        exit 1
    fi
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装，请先安装 Node.js 20+"
        exit 1
    fi
    
    log "✓ 前置条件检查通过"
}

# 数据库迁移
run_database_migrations() {
    log "执行数据库迁移..."
    
    # 检查 MySQL 连接
    log "检查 MySQL 连接..."
    if ! docker exec mysql8 mysql -uroot -p'Admin168' dmh -e "SELECT 1" &> /dev/null; then
        log_error "MySQL 连接失败"
        exit 1
    fi
    
    # 执行迁移脚本
    for migration in $BACKEND_DIR/migrations/202*.sql; do
        if [ -f "$migration" ]; then
            log "执行迁移脚本: $(basename $migration)"
            docker exec -i mysql8 mysql -uroot -p'Admin168' dmh < "$migration"
        fi
    done
    
    log "✓ 数据库迁移完成"
}

# 构建后端
build_backend() {
    log "构建后端服务..."
    cd "$BACKEND_DIR"
    
    # 清理旧构建
    rm -rf bin/dmh-api
    
    # 编译
    log "编译后端二进制文件..."
    go build -o bin/dmh-api api/dmh.go
    
    if [ ! -f "bin/dmh-api" ]; then
        log_error "后端构建失败"
        exit 1
    fi
    
    log "✓ 后端构建完成"
}

# 构建前端
build_frontend() {
    log "构建前端服务..."
    
    # H5 前端
    log "构建 H5 前端..."
    cd "$PROJECT_ROOT/frontend-h5"
    npm install
    npm run build
    
    if [ ! -d "dist" ]; then
        log_error "H5 前端构建失败"
        exit 1
    fi
    
    # 管理后台前端
    log "构建管理后台前端..."
    cd "$PROJECT_ROOT/frontend-admin"
    npm install
    npm run build
    
    if [ ! -d "dist" ]; then
        log_error "管理后台前端构建失败"
        exit 1
    fi
    
    log "✓ 前端构建完成"
}

# 备份现有版本
backup_current_version() {
    log "备份当前版本..."
    
    BACKUP_DIR="$PROJECT_ROOT/deployment/backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # 备份后端
    if [ -f "$BACKEND_DIR/bin/dmh-api" ]; then
        cp "$BACKEND_DIR/bin/dmh-api" "$BACKUP_DIR/"
        log "✓ 后端已备份"
    fi
    
    # 备份前端
    if [ -d "$PROJECT_ROOT/frontend-h5/dist" ]; then
        cp -r "$PROJECT_ROOT/frontend-h5/dist" "$BACKUP_DIR/h5-dist"
        log "✓ H5 前端已备份"
    fi
    
    if [ -d "$PROJECT_ROOT/frontend-admin/dist" ]; then
        cp -r "$PROJECT_ROOT/frontend-admin/dist" "$BACKUP_DIR/admin-dist"
        log "✓ 管理后台前端已备份"
    fi
    
    log "✓ 备份完成: $BACKUP_DIR"
    echo "$BACKUP_DIR" > "$PROJECT_ROOT/deployment/LAST_BACKUP.txt"
}

# 部署服务
deploy_services() {
    log "部署服务..."

    cd "$PROJECT_ROOT/deployment"

    # 停止现有服务
    log "停止现有服务..."
    docker compose -f "$DOCKER_COMPOSE" down

    # 构建并启动新服务
    log "启动新服务..."
    docker compose -f "$DOCKER_COMPOSE" build --no-cache
    docker compose -f "$DOCKER_COMPOSE" up -d

    log "✓ 服务部署完成"
}

# 健康检查
health_check() {
    log "执行健康检查..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        log "健康检查尝试 $attempt/$max_attempts..."
        
        # 检查后端服务
        if curl -f http://localhost:8889/health &> /dev/null; then
            log "✓ 后端服务正常"
            return 0
        fi
        
        sleep 2
        ((attempt++))
    done
    
    log_error "健康检查失败：服务未在预期时间内启动"
    return 1
}

# 功能验证
verify_deployment() {
    log "执行功能验证..."
    
    # 检查数据库迁移
    log "验证数据库迁移..."
    if docker exec mysql8 mysql -uroot -p'Admin168' dmh -e "SHOW COLUMNS FROM campaigns LIKE 'payment_config'" &> /dev/null; then
        log "✓ payment_config 字段存在"
    else
        log_error "payment_config 字段不存在"
    fi
    
    if docker exec mysql8 mysql -uroot -p'Admin168' dmh -e "SHOW COLUMNS FROM campaigns LIKE 'poster_template_id'" &> /dev/null; then
        log "✓ poster_template_id 字段存在"
    else
        log_error "poster_template_id 字段不存在"
    fi
    
    if docker exec mysql8 mysql -uroot -p'Admin168' dmh -e "SHOW COLUMNS FROM orders LIKE 'verification_status'" &> /dev/null; then
        log "✓ verification_status 字段存在"
    else
        log_error "verification_status 字段不存在"
    fi
    
    if docker exec mysql8 mysql -uroot -p'Admin168' dmh -e "SELECT COUNT(*) FROM poster_template_configs" &> /dev/null; then
        log "✓ 海报模板表存在"
    else
        log_error "海报模板表不存在"
    fi
    
    log "✓ 功能验证完成"
}

# 清理函数
cleanup() {
    log "清理临时文件..."
    
    # 清理旧备份（保留最近 5 个）
    if [ -d "$PROJECT_ROOT/deployment/backups" ]; then
        ls -t "$PROJECT_ROOT/deployment/backups" | tail -n +6 | xargs rm -rf
    fi
    
    log "✓ 清理完成"
}

# 主函数
main() {
    echo "========================================="
    echo "DMH 活动高级功能部署"
    echo "========================================="
    echo ""
    
    # 检查前置条件
    check_prerequisites
    
    # 备份现有版本
    backup_current_version
    
    # 数据库迁移
    run_database_migrations
    
    # 构建服务
    build_backend
    build_frontend
    
    # 部署服务
    deploy_services
    
    # 健康检查
    if health_check; then
        # 功能验证
        verify_deployment
        
        # 清理
        cleanup
        
        echo ""
        log "========================================="
        log "✓ 部署成功！"
        log "========================================="
        echo ""
        echo "服务访问地址："
        echo "  - 后端 API: http://localhost:8889"
        echo "  - H5 前端: http://localhost:3100"
        echo "  - 管理后台: http://localhost:3000"
        echo ""
        echo "备份位置: $PROJECT_ROOT/deployment/backups/"
        echo "日志文件: $LOG_FILE"
        echo ""
        exit 0
    else
        log_error "部署失败，请检查日志：$LOG_FILE"
        exit 1
    fi
}

# 捕获退出信号
trap cleanup EXIT
trap 'log_error "部署被中断"; exit 1' INT TERM

# 执行主函数
main
