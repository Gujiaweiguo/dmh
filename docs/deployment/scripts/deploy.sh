#!/bin/bash

# DMH系统部署脚本
# 使用方法: ./deploy.sh [环境] [版本]
# 示例: ./deploy.sh production v1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ $# -lt 1 ]; then
    log_error "使用方法: $0 <环境> [版本]"
    log_info "环境选项: development, staging, production"
    exit 1
fi

ENVIRONMENT=$1
VERSION=${2:-"latest"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
DEPLOY_DIR="/opt/dmh"

log_info "开始部署DMH系统"
log_info "环境: $ENVIRONMENT"
log_info "版本: $VERSION"
log_info "项目根目录: $PROJECT_ROOT"

# 检查环境
check_environment() {
    log_info "检查部署环境..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    # 检查Git
    if ! command -v git &> /dev/null; then
        log_error "Git未安装，请先安装Git"
        exit 1
    fi
    
    log_success "环境检查通过"
}

# 创建部署目录
create_directories() {
    log_info "创建部署目录..."
    
    sudo mkdir -p $DEPLOY_DIR/{config,logs,data,backup,ssl}
    sudo mkdir -p $DEPLOY_DIR/data/{uploads,cache}
    sudo chown -R $USER:$USER $DEPLOY_DIR
    
    log_success "目录创建完成"
}

# 备份当前版本
backup_current() {
    if [ -d "$DEPLOY_DIR/current" ]; then
        log_info "备份当前版本..."
        
        BACKUP_DIR="$DEPLOY_DIR/backup/$(date +%Y%m%d_%H%M%S)"
        mkdir -p $BACKUP_DIR
        
        # 备份配置文件
        cp -r $DEPLOY_DIR/config $BACKUP_DIR/
        
        # 备份数据库
        if docker ps | grep -q dmh-mysql; then
            log_info "备份数据库..."
            docker exec dmh-mysql mysqldump -u root -p$MYSQL_ROOT_PASSWORD dmh > $BACKUP_DIR/database.sql
        fi
        
        log_success "备份完成: $BACKUP_DIR"
    fi
}

# 拉取代码
pull_code() {
    log_info "拉取最新代码..."
    
    if [ ! -d "$DEPLOY_DIR/current" ]; then
        git clone https://github.com/your-org/dmh-system.git $DEPLOY_DIR/current
    else
        cd $DEPLOY_DIR/current
        git fetch origin
        if [ "$VERSION" = "latest" ]; then
            git checkout main
            git pull origin main
        else
            git checkout $VERSION
        fi
    fi
    
    log_success "代码拉取完成"
}

# 构建前端
build_frontend() {
    log_info "构建前端项目..."
    
    # 构建管理后台
    cd $DEPLOY_DIR/current/frontend-admin
    if [ ! -d "node_modules" ]; then
        npm install
    fi
    npm run build
    
    # 构建H5前端
    cd $DEPLOY_DIR/current/frontend-h5
    if [ ! -d "node_modules" ]; then
        npm install
    fi
    npm run build
    
    log_success "前端构建完成"
}

# 配置环境文件
setup_config() {
    log_info "配置环境文件..."
    
    # 复制环境配置
    if [ ! -f "$DEPLOY_DIR/.env" ]; then
        cp $DEPLOY_DIR/current/docs/deployment/.env.example $DEPLOY_DIR/.env
        log_warning "请编辑 $DEPLOY_DIR/.env 文件配置环境变量"
        read -p "按回车键继续..."
    fi
    
    # 复制Docker Compose配置
    cp $DEPLOY_DIR/current/docs/deployment/docker-compose.yml $DEPLOY_DIR/
    
    # 复制应用配置
    if [ ! -f "$DEPLOY_DIR/config/config.yaml" ]; then
        cp $DEPLOY_DIR/current/backend/config/config.yaml.example $DEPLOY_DIR/config/config.yaml
        log_warning "请编辑 $DEPLOY_DIR/config/config.yaml 文件配置应用参数"
    fi
    
    # 复制Nginx配置
    mkdir -p $DEPLOY_DIR/config/nginx/sites
    cp $DEPLOY_DIR/current/docs/deployment/config/nginx.conf $DEPLOY_DIR/config/
    cp $DEPLOY_DIR/current/docs/deployment/config/nginx/sites/* $DEPLOY_DIR/config/nginx/sites/
    
    log_success "配置文件设置完成"
}

# 设置SSL证书
setup_ssl() {
    log_info "设置SSL证书..."
    
    if [ ! -f "$DEPLOY_DIR/ssl/dmh.crt" ]; then
        log_warning "SSL证书不存在，生成自签名证书用于测试"
        
        # 生成自签名证书
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout $DEPLOY_DIR/ssl/dmh.key \
            -out $DEPLOY_DIR/ssl/dmh.crt \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=DMH/CN=*.dmh.com"
        
        log_warning "生产环境请使用正式SSL证书"
    fi
    
    log_success "SSL证书设置完成"
}

# 部署服务
deploy_services() {
    log_info "部署服务..."
    
    cd $DEPLOY_DIR
    
    # 加载环境变量
    export $(cat .env | grep -v '^#' | xargs)
    
    # 停止现有服务
    if docker-compose ps | grep -q Up; then
        log_info "停止现有服务..."
        docker-compose down
    fi
    
    # 构建并启动服务
    log_info "构建并启动服务..."
    docker-compose up -d --build
    
    log_success "服务部署完成"
}

# 等待服务启动
wait_for_services() {
    log_info "等待服务启动..."
    
    # 等待数据库启动
    log_info "等待数据库启动..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T mysql mysqladmin ping -h localhost --silent; then
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "数据库启动超时"
        exit 1
    fi
    
    # 等待API服务启动
    log_info "等待API服务启动..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -f http://localhost:8080/api/v1/health &> /dev/null; then
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "API服务启动超时"
        exit 1
    fi
    
    log_success "所有服务启动完成"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    
    # 检查是否需要初始化
    if docker-compose exec -T mysql mysql -u dmh_user -p$MYSQL_PASSWORD dmh -e "SELECT COUNT(*) FROM users;" &> /dev/null; then
        log_info "数据库已初始化，跳过"
        return
    fi
    
    # 执行初始化脚本
    docker-compose exec -T mysql mysql -u dmh_user -p$MYSQL_PASSWORD dmh < $DEPLOY_DIR/current/backend/scripts/init.sql
    
    log_success "数据库初始化完成"
}

# 运行测试
run_tests() {
    log_info "运行部署测试..."
    
    # 测试API健康检查
    if ! curl -f http://localhost:8080/api/v1/health; then
        log_error "API健康检查失败"
        return 1
    fi
    
    # 测试前端访问
    if ! curl -f http://localhost/; then
        log_error "前端访问测试失败"
        return 1
    fi
    
    # 测试数据库连接
    if ! docker-compose exec -T mysql mysql -u dmh_user -p$MYSQL_PASSWORD dmh -e "SELECT 1;"; then
        log_error "数据库连接测试失败"
        return 1
    fi
    
    log_success "所有测试通过"
}

# 清理旧版本
cleanup() {
    log_info "清理旧版本..."
    
    # 清理旧的Docker镜像
    docker image prune -f
    
    # 清理旧的备份文件（保留最近7天）
    find $DEPLOY_DIR/backup -type d -mtime +7 -exec rm -rf {} + 2>/dev/null || true
    
    log_success "清理完成"
}

# 显示部署信息
show_info() {
    log_success "部署完成！"
    echo
    echo "访问地址:"
    echo "  管理后台: http://localhost (或配置的域名)"
    echo "  H5前端: http://localhost/h5"
    echo "  API接口: http://localhost:8080"
    echo
    echo "服务状态:"
    docker-compose ps
    echo
    echo "日志查看:"
    echo "  docker-compose logs -f [服务名]"
    echo
    echo "服务管理:"
    echo "  启动: docker-compose up -d"
    echo "  停止: docker-compose down"
    echo "  重启: docker-compose restart"
}

# 主流程
main() {
    check_environment
    create_directories
    backup_current
    pull_code
    build_frontend
    setup_config
    setup_ssl
    deploy_services
    wait_for_services
    init_database
    
    if run_tests; then
        cleanup
        show_info
    else
        log_error "部署测试失败，请检查日志"
        exit 1
    fi
}

# 错误处理
trap 'log_error "部署过程中发生错误，请检查日志"; exit 1' ERR

# 执行主流程
main

log_success "DMH系统部署完成！"