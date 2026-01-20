-- 分销商系统数据库迁移脚本（新版本 - 自动成为分销商）
-- 创建日期: 2025-01-20
-- 基于 add-distributor-role 提案的新需求

-- ============================================================
-- 第一部分：删除不需要的旧表（申请审批流程相关）
-- ============================================================

-- 删除分销商申请表（不再需要申请审批）
DROP TABLE IF EXISTS `distributor_applications`;

-- 删除分销商级别奖励配置表（改为活动级别配置）
DROP TABLE IF EXISTS `distributor_level_rewards`;

-- 删除分销商推广链接表（改用二维码海报）
DROP TABLE IF EXISTS `distributor_links`;

-- ============================================================
-- 第二部分：修改现有表
-- ============================================================

-- 1. 修改 distributors 表（删除审批相关字段）
ALTER TABLE `distributors`
  DROP COLUMN `approved_by`,
  DROP COLUMN `approved_at`,
  MODIFY COLUMN `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/suspended)';

-- 2. 修改 distributor_rewards 表（简化字段）
ALTER TABLE `distributor_rewards`
  DROP COLUMN `user_id`,
  DROP COLUMN `reward_rate`,
  DROP COLUMN `from_user_id`,
  ADD COLUMN `distributor_level` INT NOT NULL DEFAULT 1 COMMENT '分销商级别(1/2/3)',
  ADD COLUMN `reward_percentage` DECIMAL(5,2) NOT NULL COMMENT '奖励比例';

-- ============================================================
-- 第三部分：扩展现有表
-- ============================================================

-- 3. 扩展 campaigns 表（增加分销相关字段）
ALTER TABLE `campaigns`
  ADD COLUMN `enable_distribution` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否启用分销',
  ADD COLUMN `distribution_level` INT NOT NULL DEFAULT 1 COMMENT '分销层级(1/2/3)',
  ADD COLUMN `distribution_rewards` JSON COMMENT '各级奖励比例 {"level1": 10, "level2": 8, "level3": 5}',
  ADD INDEX `idx_enable_distribution` (`enable_distribution`);

-- 4. 扩展 orders 表（增加分销链路径）
ALTER TABLE `orders`
  ADD COLUMN `distributor_path` VARCHAR(100) DEFAULT '' COMMENT '分销链路径 "一级ID,二级ID,三级ID"',
  ADD INDEX `idx_distributor_path` (`distributor_path`(50));

-- ============================================================
-- 第四部分：创建新表
-- ============================================================

-- 5. 创建 withdrawals 表（提现申请表）
CREATE TABLE IF NOT EXISTS `withdrawals` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '提现ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID',
  `distributor_id` BIGINT NOT NULL COMMENT '分销商ID',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '提现金额',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态(pending/approved/rejected/processing/completed/failed)',
  `pay_type` VARCHAR(20) NOT NULL COMMENT '提现方式(wechat/alipay/bank)',
  `pay_account` VARCHAR(100) NOT NULL COMMENT '提现账号',
  `pay_real_name` VARCHAR(50) COMMENT '真实姓名',
  `approved_by` BIGINT DEFAULT NULL COMMENT '审批人ID',
  `approved_at` DATETIME DEFAULT NULL COMMENT '审批时间',
  `approved_notes` TEXT COMMENT '审批备注',
  `rejected_reason` TEXT COMMENT '拒绝原因',
  `paid_at` DATETIME DEFAULT NULL COMMENT '打款时间',
  `trade_no` VARCHAR(100) DEFAULT NULL COMMENT '交易流水号',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_distributor_id` (`distributor_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提现申请表';

-- 6. 创建 poster_templates 表（海报模板表）
CREATE TABLE IF NOT EXISTS `poster_templates` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `type` VARCHAR(20) NOT NULL COMMENT '类型(campaign/distributor)',
  `campaign_id` BIGINT DEFAULT NULL COMMENT '活动ID（活动海报才有）',
  `distributor_id` BIGINT DEFAULT NULL COMMENT '分销商ID（分销商海报才有）',
  `template_url` VARCHAR(500) NOT NULL COMMENT '海报URL',
  `poster_data` JSON COMMENT '海报数据（包含活动信息、二维码等）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_campaign_id` (`campaign_id`),
  KEY `idx_distributor_id` (`distributor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报模板表';

-- ============================================================
-- 第五部分：数据迁移和初始化
-- ============================================================

-- 7. 为现有分销商记录更新状态（pending -> active）
UPDATE `distributors` SET `status` = 'active' WHERE `status` = 'pending';

-- 8. 为现有活动添加默认分销配置（可选）
-- 注意：这里只为示例活动添加默认配置，实际使用时品牌管理员需要手动配置
UPDATE `campaigns`
SET
  `enable_distribution` = FALSE,
  `distribution_level` = 1,
  `distribution_rewards` = JSON_OBJECT('level1', 10)
WHERE `enable_distribution` IS NULL;

-- 9. 确保 distributor 角色存在
INSERT INTO `roles` (`name`, `code`, `description`, `created_at`, `updated_at`)
VALUES ('分销商', 'distributor', '具备推广资格的高级顾客角色', NOW(), NOW())
ON DUPLICATE KEY UPDATE `description` = '具备推广资格的高级顾客角色';

-- 10. 添加提现权限（如果权限表存在）
-- 注意：这里需要根据实际的权限表结构调整
-- INSERT INTO `permissions` ...

-- ============================================================
-- 第六部分：清理和验证
-- ============================================================

-- 11. 验证表结构是否正确
-- SELECT COUNT(*) FROM information_schema.tables
-- WHERE table_schema = 'dmh'
-- AND table_name IN ('distributors', 'withdrawals', 'poster_templates', 'campaigns', 'orders');

-- 12. 显示修改结果
SHOW COLUMNS FROM `distributors`;
SHOW COLUMNS FROM `withdrawals`;
SHOW COLUMNS FROM `poster_templates`;
SHOW COLUMNS FROM `campaigns`;
SHOW COLUMNS FROM `orders`;
