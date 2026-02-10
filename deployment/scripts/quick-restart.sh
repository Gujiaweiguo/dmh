#!/bin/bash
# DMH 容器快速重启脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$DEPLOYMENT_DIR/docker-compose-simple.yml"

log "快速重启 DMH 容器..."
cd "$DEPLOYMENT_DIR"
docker compose -f "$COMPOSE_FILE" restart

sleep 5

# 检查状态
info "检查容器状态..."
docker compose -f "$COMPOSE_FILE" ps

log "重启完成！"
info "查看日志: docker compose logs -f"
