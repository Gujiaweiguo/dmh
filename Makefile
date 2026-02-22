# DMH - Digital Marketing Hub
# 常用命令 Makefile

.PHONY: help up down restart logs test build clean test-integration test-e2e test-e2e-headless test-order-regression backend-coverage admin-coverage h5-coverage spec-validate-all full-regression

# 默认显示帮助
help:
	@echo "DMH - Digital Marketing Hub"
	@echo ""
	@echo "常用命令:"
	@echo "  make up          - 启动所有服务 (Docker Compose)"
	@echo "  make down        - 停止所有服务"
	@echo "  make restart     - 重启所有服务"
	@echo "  make logs        - 查看服务日志"
	@echo "  make ps          - 查看运行状态"
	@echo "  make test        - 运行所有测试"
	@echo "  make test-backend - 运行后端测试"
	@echo "  make test-admin  - 运行管理后台测试"
	@echo "  make test-h5     - 运行 H5 测试"
	@echo "  make test-order-regression - 运行订单专项回归"
	@echo "  make test-e2e-headless - 运行双端 E2E（headless）"
	@echo "  make backend-coverage - 后端覆盖率门禁 (>=76%)"
	@echo "  make admin-coverage - Admin 覆盖率门禁 (>=70%)"
	@echo "  make h5-coverage - H5 覆盖率门禁 (>=44%)"
	@echo "  make full-regression - 本地一键全量回归"
	@echo "  make build       - 构建前端生产包"
	@echo "  make clean       - 清理临时文件"
	@echo "  make db-migrate  - 运行数据库迁移"
	@echo "  make db-backup   - 备份数据库"
	@echo "  make update      - 更新代码并重启"

# Docker Compose 命令
COMPOSE = docker compose -f deploy/docker-compose-simple.yml

up:
	$(COMPOSE) up -d
	@echo "✓ 服务已启动"
	@echo "  后端 API: http://localhost:8889"
	@echo "  管理后台: http://localhost:3000"
	@echo "  H5 前端: http://localhost:3100"

down:
	$(COMPOSE) down
	@echo "✓ 服务已停止"

restart:
	$(COMPOSE) restart
	@echo "✓ 服务已重启"

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

# 测试命令
test: test-backend test-admin test-h5

test-backend:
	cd backend && go test -p 1 $$(go list ./... | grep -v -E 'dmh/test/integration|dmh/test/performance') -v

test-integration:
	cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1

test-admin:
	cd frontend-admin && npm run test

test-h5:
	cd frontend-h5 && npm run test

test-order-regression:
	backend/scripts/run_order_mysql8_regression.sh

test-e2e:
	cd frontend-admin && npm run test:e2e
	cd frontend-h5 && npm run test:e2e

test-e2e-headless:
	cd frontend-admin && npm run test:e2e:headless
	cd frontend-h5 && npm run test:e2e:headless

backend-coverage:
	cd backend && go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic
	cd backend && go tool cover -func=coverage.out | awk '/total:/ { gsub("%","",$$3); if ($$3 + 0 >= 76) { print "Backend coverage OK (" $$3 "%)"; exit 0 } else { print "Backend coverage not enough: " $$3 "%"; exit 1 } }'

admin-coverage:
	cd frontend-admin && npm run test:cov 2>&1 | tee /tmp/dmh-admin-coverage.log
	@PCT=$$(grep -E "All files\s*\|" /tmp/dmh-admin-coverage.log | tail -1 | awk -F'|' '{gsub(/ /,"",$$2); print $$2}' | tr -d '%' || true); \
	if echo "$$PCT" | grep -Eq '^[0-9]+(\.[0-9]+)?$$' && [ "$${PCT%.*}" -ge 70 ]; then \
	  echo "Frontend Admin coverage OK ($$PCT%)"; \
	else \
	  echo "Frontend Admin coverage not enough: $$PCT%"; \
	  exit 1; \
	fi

h5-coverage:
	cd frontend-h5 && npm run test:cov 2>&1 | tee /tmp/dmh-h5-coverage.log
	@PCT=$$(grep -E "All files\s*\|" /tmp/dmh-h5-coverage.log | tail -1 | awk -F'|' '{gsub(/ /,"",$$2); print $$2}' | tr -d '%' || true); \
	if echo "$$PCT" | grep -Eq '^[0-9]+(\.[0-9]+)?$$' && [ "$${PCT%.*}" -ge 44 ]; then \
	  echo "Frontend H5 coverage OK ($$PCT%)"; \
	else \
	  echo "Frontend H5 coverage not enough: $$PCT%"; \
	  exit 1; \
	fi

# 构建命令
build:
	cd frontend-admin && npm run build
	cd frontend-h5 && npm run build
	@echo "✓ 前端构建完成"

build-backend:
	cd backend && go build -o dmh-api api/dmh.go
	@echo "✓ 后端构建完成"

# 数据库命令
db-migrate:
	@echo "运行数据库迁移..."
	@read -p "输入迁移文件名 (如: add_user_table): " name; \
	file="backend/migrations/$$(date +%Y%m%d)_$${name}.sql"; \
	echo "-- Migration: $${name}" > $$file; \
	echo "-- Date: $$(date)" >> $$file; \
	echo "" >> $$file; \
	echo "✓ 迁移文件已创建: $$file"

db-backup:
	@mkdir -p backups
	docker exec mysql8 mysqldump -uroot -p'#Admin168' dmh > backups/dmh_backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✓ 数据库已备份到 backups/"

# 清理命令
clean:
	@echo "清理临时文件..."
	@rm -rf frontend-*/dist
	@rm -rf frontend-*/coverage
	@rm -rf frontend-*/test-results
	@rm -rf frontend-*/playwright-report
	@rm -f backend/*.out
	@rm -f *.log
	@docker system prune -f 2>/dev/null || true
	@echo "✓ 清理完成"

# 代码质量
check:
	cd backend && gofmt -d .
	cd frontend-admin && npm run lint
	cd frontend-h5 && npm run lint

fmt:
	cd backend && gofmt -w .
	@echo "✓ Go 代码已格式化"

# 开发快捷命令
dev-backend:
	cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml

dev-admin:
	cd frontend-admin && npm run dev

dev-h5:
	cd frontend-h5 && npm run dev

# OpenSpec 命令
spec-list:
	@cd openspec && openspec list 2>/dev/null || echo "openspec CLI 未安装"

spec-validate:
	@openspec validate --strict --no-interactive 2>/dev/null || echo "openspec CLI 未安装"

spec-validate-all:
	@openspec validate --all --no-interactive

full-regression: test-backend test-integration test-order-regression test-admin test-h5 test-e2e-headless backend-coverage admin-coverage h5-coverage spec-validate-all
	@echo "✓ 全量回归完成"

# 更新命令
update:
	git pull origin main
	$(COMPOSE) pull
	$(COMPOSE) up -d
	@echo "✓ 代码和服务已更新"

# 完整部署
deploy: build
	cp backend/dmh-api deploy/
	$(COMPOSE) restart dmh-api
	@echo "✓ 部署完成"
