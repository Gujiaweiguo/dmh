#!/bin/bash

# DMH 项目管理脚本
# 用法: ./dmh.sh [start|init|stop|restart|status|logs]

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
MYSQL_HOST="${MYSQL_HOST:-172.17.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASS="${MYSQL_PASS:-#Admin168}"
MYSQL_DB="${MYSQL_DB:-dmh}"
DOCKER_MYSQL_CONTAINER="${DOCKER_MYSQL_CONTAINER:-mysql8}"

# 实际找到的容器名（运行时更新）
FOUND_MYSQL_CONTAINER=""

# 工具函数
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

print_usage() {
    echo "用法: ./dmh.sh [命令]"
    echo ""
    echo "命令:"
    echo "  init      - 初始化环境（首次运行，包含数据库初始化）"
    echo "  start     - 启动所有服务（日常使用）"
    echo "  stop      - 停止所有服务"
    echo "  restart   - 重启所有服务"
    echo "  status    - 查看服务状态"
    echo "  logs      - 查看服务日志"
    echo ""
    echo "示例:"
    echo "  ./dmh.sh init     # 首次运行"
    echo "  ./dmh.sh start    # 日常启动"
    echo "  ./dmh.sh stop     # 停止服务"
}

# 启动 Docker daemon (WSL2 环境)
start_docker_daemon() {
    # 检查 dockerd 是否运行
    if pgrep -x dockerd > /dev/null 2>&1; then
        return 0
    fi
    
    echo -e "${YELLOW}启动 Docker daemon...${NC}"
    
    # 启动 dockerd
    sudo dockerd > /tmp/dockerd.log 2>&1 &
    
    # 等待 docker.sock 创建
    local count=0
    while [ ! -S /var/run/docker.sock ] && [ $count -lt 30 ]; do
        sleep 1
        count=$((count + 1))
    done
    
    if [ -S /var/run/docker.sock ]; then
        # 修改权限
        sudo chmod 666 /var/run/docker.sock 2>/dev/null
        echo -e "${GREEN}✓ Docker daemon 已启动${NC}"
        return 0
    else
        echo -e "${RED}✗ Docker daemon 启动失败${NC}"
        return 1
    fi
}

# 检查 Docker MySQL
check_docker_mysql() {
    echo -e "${YELLOW}检查 Docker MySQL...${NC}"
    
    if ! command_exists docker; then
        echo -e "${RED}✗ Docker 未安装${NC}"
        return 1
    fi
    
    # 检查 Docker 是否运行
    if ! docker info > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠ Docker 未运行${NC}"
        
        # 检测是否在 WSL 环境
        if grep -qi microsoft /proc/version 2>/dev/null; then
            echo -e "${YELLOW}检测到 WSL 环境，尝试启动 Docker daemon${NC}"
            if ! start_docker_daemon; then
                return 1
            fi
        else
            # 尝试启动 Docker (非 WSL 环境)
            echo -e "${YELLOW}正在尝试启动 Docker...${NC}"
            if sudo systemctl start docker 2>/dev/null || sudo service docker start 2>/dev/null; then
                echo -e "${GREEN}✓ Docker 已启动${NC}"
                sleep 3
            else
                echo -e "${RED}✗ 无法启动 Docker${NC}"
                echo -e "${YELLOW}请手动启动 Docker${NC}"
                return 1
            fi
        fi
        
        # 再次检查
        if ! docker info > /dev/null 2>&1; then
            echo -e "${RED}✗ Docker 仍未运行${NC}"
            return 1
        fi
    fi
    
    echo -e "${GREEN}✓ Docker 运行中${NC}"
    
    # 查找 MySQL 容器 - 优先使用配置的名字
    local container
    container=$(docker ps -a --filter "name=^${DOCKER_MYSQL_CONTAINER}$" --format "{{.Names}}" 2>/dev/null | head -1)
    
    if [ -z "$container" ]; then
        # 模糊匹配包含 mysql 的容器
        container=$(docker ps -a --format "{{.Names}}" 2>/dev/null | grep -i mysql | head -1)
    fi
    
    if [ -z "$container" ]; then
        # 按镜像查找
        container=$(docker ps -a --filter "ancestor=mysql" --format "{{.Names}}" 2>/dev/null | head -1)
    fi
    
    if [ -n "$container" ]; then
        # 保存找到的容器名
        FOUND_MYSQL_CONTAINER="$container"
        
        local status
        status=$(docker inspect -f '{{.State.Status}}' "$container" 2>/dev/null)
        
        if [ "$status" != "running" ]; then
            echo -e "${YELLOW}启动 MySQL 容器: $container${NC}"
            docker start "$container"
            sleep 5
            
            # 再次检查状态
            status=$(docker inspect -f '{{.State.Status}}' "$container" 2>/dev/null)
            if [ "$status" == "running" ]; then
                echo -e "${GREEN}✓ MySQL 容器已启动: $container${NC}"
            else
                echo -e "${RED}✗ MySQL 容器启动失败${NC}"
                return 1
            fi
        else
            echo -e "${GREEN}✓ MySQL 容器运行中: $container${NC}"
        fi
        
        return 0
    fi
    
    echo -e "${YELLOW}⚠ 未找到 MySQL 容器${NC}"
    return 1
}

# 创建 MySQL 容器
create_mysql_container() {
    echo -e "${YELLOW}创建 MySQL 容器...${NC}"
    
    docker run -d \
        --name "$DOCKER_MYSQL_CONTAINER" \
        -p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD="$MYSQL_PASS" \
        -e MYSQL_DATABASE="$MYSQL_DB" \
        -v mysql_data:/var/lib/mysql \
        mysql:8.0
    
    # 保存容器名
    FOUND_MYSQL_CONTAINER="$DOCKER_MYSQL_CONTAINER"
    
    echo -e "${YELLOW}等待 MySQL 启动...${NC}"
    sleep 30
    
    echo -e "${GREEN}✓ MySQL 容器创建成功${NC}"
}

# 初始化数据库
init_database() {
    echo -e "${YELLOW}初始化数据库...${NC}"
    
    # 使用实际找到的容器名，如果没有则使用默认配置
    local container="${FOUND_MYSQL_CONTAINER:-$DOCKER_MYSQL_CONTAINER}"
    
    # 基础表结构脚本
    local scripts=(
        "backend/scripts/init.sql"
        "backend/scripts/create_member_system_tables.sql"
        "backend/scripts/create_external_sync_tables.sql"
        "backend/scripts/create_page_configs_table.sql"
    )
    
    for script in "${scripts[@]}"; do
        if [ -f "$script" ]; then
            echo -e "${BLUE}导入: $script${NC}"
            if ! docker exec -i "$container" mysql -u "$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" < "$script" 2>/dev/null; then
                echo -e "${YELLOW}⚠ 导入 $script 时出现警告（可能表已存在）${NC}"
            fi
        fi
    done
    
    # 询问是否导入测试数据
    echo -e "${YELLOW}是否导入测试数据（品牌、活动、会员）？(y/n)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        local test_scripts=(
            "backend/scripts/test_data.sql"
            "backend/scripts/seed_member_campaign_data.sql"
            "backend/scripts/seed_complete_test_data.sql"
        )
        
        for script in "${test_scripts[@]}"; do
            if [ -f "$script" ]; then
                echo -e "${BLUE}导入测试数据: $script${NC}"
                if ! docker exec -i "$container" mysql -u "$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" < "$script" 2>/dev/null; then
                    echo -e "${YELLOW}⚠ 导入 $script 时出现警告${NC}"
                fi
            fi
        done
        echo -e "${GREEN}✓ 测试数据导入完成${NC}"
    fi
    
    echo -e "${GREEN}✓ 数据库初始化完成${NC}"
}

# 检查端口
check_port() {
    local port=$1
    if command_exists lsof && lsof -Pi :"$port" -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${RED}✗ 端口 $port 已被占用${NC}"
        return 1
    fi
    return 0
}

# 启动后端
start_backend() {
    echo -e "${YELLOW}启动后端服务...${NC}"
    
    if ! command_exists go; then
        echo -e "${RED}✗ Go 未安装${NC}"
        return 1
    fi
    
    mkdir -p logs
    
    cd backend
    nohup go run api/dmh.go -f api/etc/dmh-api.yaml > ../logs/backend.log 2>&1 &
    echo $! > ../logs/backend.pid
    cd ..
    
    echo -e "${GREEN}✓ 后端启动 (PID: $(cat logs/backend.pid))${NC}"
    echo -e "   ${BLUE}http://localhost:8889${NC}"
}

# 启动前端
start_frontend() {
    local name=$1
    local dir=$2
    local port=$3
    
    echo -e "${YELLOW}启动 $name...${NC}"
    
    if ! command_exists npm; then
        echo -e "${RED}✗ npm 未安装${NC}"
        return 1
    fi
    
    cd "$dir"
    
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    nohup npm run dev > ../logs/"${name}".log 2>&1 &
    echo $! > ../logs/"${name}".pid
    cd ..
    
    echo -e "${GREEN}✓ $name 启动 (PID: $(cat logs/${name}.pid))${NC}"
    echo -e "   ${BLUE}http://localhost:$port${NC}"
}

# 停止服务
stop_service() {
    local name=$1
    local pidfile="logs/${name}.pid"
    
    if [ -f "$pidfile" ]; then
        local pid
        pid=$(cat "$pidfile")
        if ps -p "$pid" > /dev/null 2>&1; then
            kill "$pid" 2>/dev/null || kill -9 "$pid" 2>/dev/null
            echo -e "${GREEN}✓ 已停止 $name (PID: $pid)${NC}"
        fi
        rm -f "$pidfile"
    fi
}

# 初始化命令
cmd_init() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   DMH 初始化${NC}"
    echo -e "${GREEN}========================================${NC}"
    
    # 1. 检查/创建 MySQL
    if ! check_docker_mysql; then
        echo -e "${YELLOW}未找到 MySQL 容器，是否创建？(y/n)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            create_mysql_container
        else
            echo -e "${RED}需要 MySQL 才能继续${NC}"
            exit 1
        fi
    fi
    
    # 2. 初始化数据库
    echo -e "${YELLOW}是否初始化数据库？(y/n)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        init_database
    fi
    
    # 3. 安装依赖
    echo -e "${YELLOW}安装后端依赖...${NC}"
    cd backend && go mod download && cd ..
    
    echo -e "${YELLOW}安装 H5 前端依赖...${NC}"
    cd frontend-h5 && npm install && cd ..
    
    echo -e "${YELLOW}安装管理后台依赖...${NC}"
    cd frontend-admin && npm install && cd ..
    
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   初始化完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo -e "${YELLOW}运行 ./dmh.sh start 启动服务${NC}"
}

# 启动命令
cmd_start() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   DMH 启动${NC}"
    echo -e "${GREEN}========================================${NC}"
    
    # 检查 MySQL
    if ! check_docker_mysql; then
        echo -e "${RED}MySQL 容器不存在${NC}"
        echo -e "${YELLOW}是否创建 MySQL 容器？(y/n)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            create_mysql_container
            echo -e "${YELLOW}是否初始化数据库？(y/n)${NC}"
            read -r response
            if [[ "$response" =~ ^[Yy]$ ]]; then
                init_database
            fi
        else
            echo -e "${RED}需要 MySQL 才能启动服务${NC}"
            echo -e "${YELLOW}提示: 运行 ./dmh.sh init 进行完整初始化${NC}"
            exit 1
        fi
    fi
    
    # 检查端口
    check_port 8889 || exit 1
    check_port 3000 || exit 1
    check_port 3100 || exit 1
    
    # 启动服务
    start_backend
    sleep 3
    start_frontend "h5" "frontend-h5" 3100
    start_frontend "admin" "frontend-admin" 3000
    
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   所有服务已启动${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo -e "${BLUE}后端 API:   http://localhost:8889${NC}"
    echo -e "${BLUE}H5 前端:    http://localhost:3100${NC}"
    echo -e "${BLUE}管理后台:   http://localhost:3000${NC}"
    echo -e ""
    echo -e "${YELLOW}查看日志: ./dmh.sh logs${NC}"
    echo -e "${YELLOW}停止服务: ./dmh.sh stop${NC}"
}

# 停止命令
cmd_stop() {
    echo -e "${YELLOW}停止所有服务...${NC}"
    
    stop_service "backend"
    stop_service "h5"
    stop_service "admin"
    
    # 清理残留进程
    pkill -f "dmh.go" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true
    
    echo -e "${GREEN}所有服务已停止${NC}"
}

# 重启命令
cmd_restart() {
    cmd_stop
    sleep 2
    cmd_start
}

# 状态命令
cmd_status() {
    echo -e "${BLUE}========== 服务状态 ==========${NC}"
    
    check_service() {
        local name=$1
        local port=$2
        local pidfile="logs/${name}.pid"
        
        if [ -f "$pidfile" ]; then
            local pid
            pid=$(cat "$pidfile")
            if ps -p "$pid" > /dev/null 2>&1; then
                echo -e "${GREEN}✓ $name 运行中 (PID: $pid, Port: $port)${NC}"
            else
                echo -e "${RED}✗ $name 已停止${NC}"
            fi
        else
            echo -e "${RED}✗ $name 未启动${NC}"
        fi
    }
    
    check_service "后端" 8889
    check_service "H5" 3100
    check_service "管理后台" 3000
    
    echo ""
    if check_docker_mysql; then
        echo -e "${GREEN}✓ MySQL 运行中${NC}"
    else
        echo -e "${RED}✗ MySQL 未运行${NC}"
    fi
}

# 日志命令
cmd_logs() {
    echo "选择要查看的日志:"
    echo "1) 后端"
    echo "2) H5 前端"
    echo "3) 管理后台"
    echo "4) 全部"
    read -r choice
    
    case $choice in
        1) tail -f logs/backend.log ;;
        2) tail -f logs/h5.log ;;
        3) tail -f logs/admin.log ;;
        4) tail -f logs/*.log ;;
        *) echo "无效选择" ;;
    esac
}

# 主逻辑
case "${1:-}" in
    init)
        cmd_init
        ;;
    start)
        cmd_start
        ;;
    stop)
        cmd_stop
        ;;
    restart)
        cmd_restart
        ;;
    status)
        cmd_status
        ;;
    logs)
        cmd_logs
        ;;
    *)
        print_usage
        exit 1
        ;;
esac
