# 高级功能性能测试报告

## 📋 测试概述

本报告记录了 DMH 高级功能的性能测试，包括海报生成、二维码生成、订单核销和并发压力测试。

**测试文件位置**: `backend/test/performance/advanced_features_performance_test.go`

---

## 🎯 性能目标

| 测试项 | 目标值 | 实际状态 |
|--------|--------|----------|
| 海报生成时间 | < 3 秒 | 待测试 |
| 二维码生成时间 | < 500ms | 待测试 |
| 核销接口响应时间 | < 500ms | 待测试 |
| 并发海报生成最大响应时间 | < 10 秒 | 待测试 |

---

## 🧪 测试场景

### 12.1 海报生成性能测试

**测试方法**: 单次请求测试
**API 端点**: `POST /api/v1/campaigns/:id/poster`
**测试内容**:
- 调用海报生成接口
- 测量完整请求响应时间
- 验证生成的海报URL

**预期结果**:
- [ ] 响应时间 < 3 秒
- [ ] 返回有效的海报URL
- [ ] 数据库保存海报记录成功

**实现代码位置**: `backend/test/performance/advanced_features_performance_test.go:104`

---

### 12.2 二维码生成性能测试

**测试方法**: 多次请求取平均值
**API 端点**: `GET /api/v1/campaigns/:id/payment-qrcode`
**测试内容**:
- 连续发起 10 次请求
- 计算平均响应时间
- 验证二维码数据格式

**预期结果**:
- [ ] 平均响应时间 < 500ms
- [ ] 所有请求成功
- [ ] 返回有效的二维码URL

**实现代码位置**: `backend/test/performance/advanced_features_performance_test.go:137`

---

### 12.3 核销接口响应时间测试

**测试方法**: 多次请求取平均值
**API 端点**: `POST /api/v1/orders/verify`
**测试内容**:
- 创建测试订单
- 连续发起 10 次核销请求
- 计算平均响应时间

**预期结果**:
- [ ] 平均响应时间 < 500ms
- [ ] 核销操作成功
- [ ] 订单状态正确更新

**实现代码位置**: `backend/test/performance/advanced_features_performance_test.go:177`

---

### 12.4 并发海报生成压力测试

**测试方法**: 并发请求测试
**API 端点**: `POST /api/v1/campaigns/:id/poster`
**测试内容**:
- 同时发起 20 个并发请求
- 测量每个请求的响应时间
- 统计成功率和性能指标

**预期结果**:
- [ ] 最大响应时间 < 10 秒
- [ ] 成功率 > 95%
- [ ] 无内存泄漏或资源耗尽

**实现代码位置**: `backend/test/performance/advanced_features_performance_test.go:250`

---

## 🚀 如何运行测试

### 前置条件

1. **启动 MySQL 数据库**:
```bash
docker start mysql8
# 确认数据库运行在 127.0.0.1:3306
```

2. **运行数据库迁移**:
```bash
cd backend
mysql -h127.0.0.1 -uroot -pAdmin168 dmh < migrations/20250124_add_advanced_features.sql
```

3. **确保 Redis 可用**:
```bash
# 使用 redis-dataease (端口 16379)
# 或启动本地 Redis: docker run -d -p 6379:6379 redis:7
```

4. **编译后端服务**:
```bash
cd backend
go build -o dmh-test api/dmh.go
```

5. **启动后端服务**:
```bash
./dmh-test -f api/etc/dmh-api.yaml &
# 服务将在 http://localhost:8889 启动
```

### 运行性能测试

**运行所有测试**:
```bash
cd backend/test/performance
go test -v -run TestAdvancedFeaturesPerformanceTestSuite
```

**运行单个测试**:
```bash
# 只测试海报生成
go test -v -run Test_12_1_PosterGenerationPerformance

# 只测试二维码生成
go test -v -run Test_12_2_PaymentQRCodePerformance

# 只测试核销接口
go test -v -run Test_12_3_OrderVerifyPerformance

# 只测试并发压力
go test -v -run Test_12_4_ConcurrentPosterStressTest
```

### 测试输出示例

```
=== RUN   TestAdvancedFeaturesPerformanceTestSuite
=== RUN   TestAdvancedFeaturesPerformanceTestSuite.SetupSuite
    advanced_features_performance_test.go:61: ✓ 登录成功
    advanced_features_performance_test.go:101: ✓ 测试活动创建成功，ID: 1
=== RUN   TestAdvancedFeaturesPerformanceTestSuite.Test_12_1_PosterGenerationPerformance
    advanced_features_performance_test.go:106: 测试场景 12.1: 海报生成性能测试（目标 < 3秒）
    advanced_features_performance_test.go:125: 海报生成耗时: 2.345s
    advanced_features_performance_test.go:126: 响应状态码: 200
    advanced_features_performance_test.go:129: ✓ 海报生成成功
--- PASS: TestAdvancedFeaturesPerformanceTestSuite.Test_12_1_PosterGenerationPerformance (2.35s)
=== RUN   TestAdvancedFeaturesPerformanceTestSuite.Test_12_2_PaymentQRCodePerformance
    advanced_features_performance_test.go:139: 测试场景 12.2: 二维码生成性能测试（目标 < 500ms）
    advanced_features_performance_test.go:164: 请求 1 成功，耗时: 234ms
    advanced_features_performance_test.go:164: 请求 2 成功，耗时: 198ms
    ...
    advanced_features_performance_test.go:170: ✓ 平均耗时: 216ms (10/10 次成功)
--- PASS: TestAdvancedFeaturesPerformanceTestSuite.Test_12_2_PaymentQRCodePerformance (2.15s)
```

---

## 📊 性能指标监控

### 关键指标

| 指标 | 监控方法 | 目标值 |
|------|----------|--------|
| **海报生成时间** | 测试日志 | < 3 秒 |
| **二维码生成时间** | 测试日志 | < 500ms |
| **核销接口响应时间** | 测试日志 | < 500ms |
| **内存使用** | 系统监控 | < 500MB |
| **CPU 使用率** | 系统监控 | < 80% |
| **数据库查询时间** | 慢查询日志 | < 100ms |
| **缓存命中率** | Redis 监控 | > 90% |

### 实时监控命令

```bash
# 查看后端服务内存和CPU使用
ps aux | grep dmh-test

# 查看数据库连接数
mysql -h127.0.0.1 -uroot -pAdmin168 -e "SHOW PROCESSLIST" dmh

# 查看 Redis 内存使用
redis-cli -p 16379 INFO memory

# 查看慢查询日志
tail -f /var/log/mysql/slow-query.log
```

---

## 🐛 常见问题排查

### 问题1: 数据库连接失败

**错误信息**: `dial tcp 127.0.0.1:3306: connect: connection refused`

**解决方案**:
```bash
# 检查 MySQL 是否运行
docker ps | grep mysql8

# 启动 MySQL
docker start mysql8

# 检查端口
netstat -tuln | grep 3306
```

### 问题2: Redis 连接失败

**错误信息**: `dial tcp: lookup redis-dmh: no such host`

**解决方案**:
```bash
# 修改配置文件使用本地 Redis
# 在 backend/api/etc/dmh-api.yaml 中:
Redis:
  Host: 127.0.0.1:6379  # 改为本地Redis
```

### 问题3: 登录失败

**错误信息**: `token is empty` 或 `用户名或密码错误`

**解决方案**:
```bash
# 检查数据库中是否有测试数据
mysql -h127.0.0.1 -uroot -pAdmin168 dmh -e "SELECT * FROM users;"

# 如果没有数据，运行初始化脚本
mysql -h127.0.0.1 -uroot -pAdmin168 dmh < backend/migrations/insert_test_data.sql
```

### 问题4: 海报生成失败

**错误信息**: `活动不存在` 或 `保存海报记录失败`

**解决方案**:
```bash
# 检查活动是否存在
mysql -h127.0.0.1 -uroot -pAdmin168 dmh -e "SELECT * FROM campaigns LIMIT 1;"

# 检查海报目录权限
ls -la /tmp/dmh/posters/

# 如果目录不存在，创建它
mkdir -p /tmp/dmh/posters/
```

---

## 📈 性能优化建议

### 海报生成优化

1. **使用缓存机制**:
   - 相同参数的海报结果缓存 1 小时
   - 使用 Redis 存储缓存键

2. **异步生成**:
   - 将海报生成放入消息队列
   - 立即返回任务ID
   - 通过轮询或WebSocket通知完成

3. **CDN 加速**:
   - 生成的海报上传到 CDN
   - 减少图片传输延迟

### 二维码生成优化

1. **预生成二维码**:
   - 活动创建时预生成支付二维码
   - 避免每次请求都重新生成

2. **缓存策略**:
   - Redis 缓存二维码（TTL 2小时）
   - 使用更长的TTL减少生成次数

### 核销接口优化

1. **数据库索引优化**:
   - 确保核销码字段有索引
   - 优化查询语句

2. **批量操作**:
   - 支持批量核销接口
   - 减少网络往返次数

---

## 📝 测试检查清单

在完成测试后，请确保：

- [ ] 所有 4 个性能测试场景均已运行
- [ ] 测试结果已记录（截图或日志）
- [ ] 性能指标达到目标值
- [ ] 异常情况已记录
- [ ] 优化建议已文档化
- [ ] 测试报告已提交

---

## 📚 相关文档

- [OpenSpec 任务列表](../../openspec/changes/add-campaign-advanced-features/tasks.md)
- [海报生成 API 文档](../../docs/api/poster-api.md)
- [性能测试最佳实践](../../docs/performance-testing-guide.md)

---

**测试状态**: 待执行
**最后更新**: 2026-02-01
**负责人**: 待分配
