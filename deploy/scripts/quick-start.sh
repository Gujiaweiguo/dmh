#!/bin/bash
# DMH 容器化部署快速启动脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# 配置
DEPLOYMENT_DIR="/opt/code/dmh/deployment"
COMPOSE_FILE="$DEPLOYMENT_DIR/docker-compose-simple.yml"

log "DMH 容器化部署快速启动..."
log "========================================="

# 步骤1: 检查文件
info "检查配置文件..."
if [ ! -f "$COMPOSE_FILE" ]; then
    error "配置文件不存在: $COMPOSE_FILE"
    exit 1
fi

 # 步骤2: 停止现有容器
if docker ps | grep -q "dmh-nginx\|dmh-api"; then
    log "停止现有容器..."
    docker compose -f "$COMPOSE_FILE" down
else
    info "没有运行中的容器"
fi

 # 步骤3: 启动容器
log "启动容器（可能需要2-5分钟安装依赖）..."
docker compose -f "$COMPOSE_FILE" up -d

log "========================================="
log "容器已启动！"
log "========================================="
info ""
info "服务访问地址："
info "  管理后台:      http://localhost:3000"
info "  H5 前端:       http://localhost:3100"
info "  后端 API:      http://localhost:8889"
info ""
info "查看日志："
info "  docker compose -f $COMPOSE_FILE logs -f"
info ""
warn "首次启动需要2-5分钟安装依赖，请耐心等待..."
info "查看安装进度："
info "  docker compose -f $COMPOSE_FILE logs dmh-nginx | grep apk"
info "  docker compose -f $COMPOSE_FILE logs dmh-api | grep apk"
info ""
warn "========================================="
