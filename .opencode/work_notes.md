# 重要操作记录
## 数据库操作
- 本系统没有mysql工具，需要通过进入docker容器进行数据库操作
- 命令：docker exec -it <container_name> mysql -uroot -p123456 dmh

## 系统架构
- 后端框架：Go + go-zero
- 前端框架：Vue3 (H5) + React (Admin)
- 数据库：MySQL (在Docker容器中)


## 分销商系统实施进度（2025-01-20）

### 已完成
1. ✅ 数据库迁移
   - 创建 distributors 表
   - 创建 distributor_rewards 表
   - 创建 poster_templates 表
   - 扩展 campaigns 表（增加分销相关字段）
   - 扩展 orders 表（增加 distributor_path 字段）
   - 执行命令：docker exec -i mysql8 mysql -uroot -p'#Admin168' dmh < migration_file.sql

2. ✅ Model层更新
   - 更新 Campaign 模型（新增 EnableDistribution, DistributionLevel, DistributionRewards 字段）
   - 更新 Order 模型（新增 DistributorPath, PaidAt 字段）

3. ✅ 核心业务逻辑
   - 创建 auto_upgrade_logic.go（自动成为分销商逻辑）
   - 创建 multi_level_reward_logic.go（多级奖励计算逻辑）

### 进行中
- 后端API实现（2/10完成）

### 待完成
- 活动级别分销奖励规则配置
- 海报生成API
- 提现API
- 管理员相关API
- 前端H5实现
- 前端管理后台实现
- 测试
EOF
