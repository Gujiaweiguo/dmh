-- 分销商系统数据库迁移脚本（最终版 - 兼容MySQL 8.0）
-- 创建日期: 2025-01-20
-- 基于 add-distributor-role 提案的新需求

-- ============================================================
-- 第一部分：创建缺失的表
-- ============================================================

-- 1. 创建 distributors 表（分销商信息表）
CREATE TABLE IF NOT EXISTS `distributors` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '分销商ID',
  `user_id` BIGINT NOT NULL COMMENT '关联用户ID',
  `brand_id` BIGINT NOT NULL COMMENT '关联品牌ID',
  `level` INT NOT NULL DEFAULT 1 COMMENT '分销级别(1/2/3)',
  `parent_id` BIGINT DEFAULT NULL COMMENT '上级分销商ID',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态(active/suspended)',
  `total_earnings` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计收益',
  `subordinates_count` INT NOT NULL DEFAULT 0 COMMENT '下级人数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_brand` (`user_id`, `brand_id`),
  KEY `idx_distributor_user` (`user_id`),
  KEY `idx_distributor_brand` (`brand_id`),
  KEY `idx_parent` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商信息表';

-- 2. 创建 distributor_rewards 表（分销商奖励记录表）
CREATE TABLE IF NOT EXISTS `distributor_rewards` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '奖励ID',
  `distributor_id` BIGINT NOT NULL COMMENT '分销商ID',
  `order_id` BIGINT NOT NULL COMMENT '订单ID',
  `campaign_id` BIGINT NOT NULL COMMENT '活动ID',
  `amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
  `distributor_level` INT NOT NULL DEFAULT 1 COMMENT '分销商级别(1/2/3)',
  `reward_percentage` DECIMAL(5,2) NOT NULL COMMENT '奖励比例',
  `status` VARCHAR(20) NOT NULL DEFAULT 'settled' COMMENT '状态',
  `settled_at` DATETIME NULL COMMENT '结算时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_distributor` (`distributor_id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_campaign` (`campaign_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商奖励记录表';

-- 3. 创建 poster_templates 表（海报模板表）
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
-- 第二部分：扩展现有表
-- ============================================================

-- 4. 扩展 campaigns 表（增加分销相关字段）
-- 注意：MySQL 8.0 不支持 IF NOT EXISTS，需要手动检查
-- 使用存储过程来安全添加字段
DELIMITER $$
DROP PROCEDURE IF EXISTS add_campaign_distribution_fields$$
CREATE PROCEDURE add_campaign_distribution_fields()
BEGIN
  -- 添加 enable_distribution 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'campaigns'
    AND column_name = 'enable_distribution'
  ) THEN
    ALTER TABLE `campaigns`
    ADD COLUMN `enable_distribution` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否启用分销';
  END IF;

  -- 添加 distribution_level 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'campaigns'
    AND column_name = 'distribution_level'
  ) THEN
    ALTER TABLE `campaigns`
    ADD COLUMN `distribution_level` INT NOT NULL DEFAULT 1 COMMENT '分销层级(1/2/3)';
  END IF;

  -- 添加 distribution_rewards 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'campaigns'
    AND column_name = 'distribution_rewards'
  ) THEN
    ALTER TABLE `campaigns`
    ADD COLUMN `distribution_rewards` JSON COMMENT '各级奖励比例 {"level1": 10, "level2": 8, "level3": 5}';
  END IF;

  -- 添加索引
  IF NOT EXISTS (
    SELECT * FROM information_schema.statistics
    WHERE table_schema = 'dmh'
    AND table_name = 'campaigns'
    AND index_name = 'idx_enable_distribution'
  ) THEN
    ALTER TABLE `campaigns`
    ADD INDEX `idx_enable_distribution` (`enable_distribution`);
  END IF;
END$$
DELIMITER ;

CALL add_campaign_distribution_fields();
DROP PROCEDURE IF EXISTS add_campaign_distribution_fields;

-- 5. 扩展 orders 表（增加分销链路径）
DELIMITER $$
DROP PROCEDURE IF EXISTS add_order_distributor_path$$
CREATE PROCEDURE add_order_distributor_path()
BEGIN
  -- 添加 distributor_path 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'orders'
    AND column_name = 'distributor_path'
  ) THEN
    ALTER TABLE `orders`
    ADD COLUMN `distributor_path` VARCHAR(100) DEFAULT '' COMMENT '分销链路径 "一级ID,二级ID,三级ID"';
  END IF;

  -- 添加索引
  IF NOT EXISTS (
    SELECT * FROM information_schema.statistics
    WHERE table_schema = 'dmh'
    AND table_name = 'orders'
    AND index_name = 'idx_distributor_path'
  ) THEN
    ALTER TABLE `orders`
    ADD INDEX `idx_distributor_path` (`distributor_path`(50));
  END IF;
END$$
DELIMITER ;

CALL add_order_distributor_path();
DROP PROCEDURE IF EXISTS add_order_distributor_path;

-- ============================================================
-- 第三部分：初始化数据
-- ============================================================

-- 6. 确保 distributor 角色存在
INSERT IGNORE INTO `roles` (`name`, `code`, `description`, `created_at`, `updated_at`)
VALUES ('分销商', 'distributor', '具备推广资格的高级顾客角色', NOW(), NOW());

-- 7. 为现有活动添加默认分销配置
UPDATE `campaigns`
SET
  `enable_distribution` = FALSE,
  `distribution_level` = 1,
  `distribution_rewards` = JSON_OBJECT('level1', 10)
WHERE `enable_distribution` IS NULL;

-- ============================================================
-- 第四部分：验证
-- ============================================================

-- 显示创建结果
SELECT 'Migration completed successfully!' AS status;
SHOW TABLES LIKE 'distributor%';
SHOW TABLES LIKE 'poster_templates';

-- ============================================================
-- 第五部分：扩展 withdrawals 表
-- ============================================================

-- 扩展 withdrawals 表（增加分销相关字段）
DELIMITER $$
DROP PROCEDURE IF EXISTS add_withdrawal_distributor_fields$$
CREATE PROCEDURE add_withdrawal_distributor_fields()
BEGIN
  -- 添加 brand_id 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'brand_id'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `brand_id` BIGINT NOT NULL DEFAULT 0 COMMENT '品牌ID' AFTER `user_id`,
    ADD INDEX `idx_brand_id` (`brand_id`);
  END IF;

  -- 添加 distributor_id 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'distributor_id'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `distributor_id` BIGINT NOT NULL DEFAULT 0 COMMENT '分销商ID' AFTER `brand_id`,
    ADD INDEX `idx_distributor_id` (`distributor_id`);
  END IF;

  -- 添加 pay_type 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'pay_type'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `pay_type` VARCHAR(20) NOT NULL DEFAULT 'wechat' COMMENT '提现方式' AFTER `status`;
  END IF;

  -- 添加 pay_account 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'pay_account'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `pay_account` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '提现账号' AFTER `pay_type`;
  END IF;

  -- 添加 pay_real_name 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'pay_real_name'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `pay_real_name` VARCHAR(50) DEFAULT '' COMMENT '真实姓名' AFTER `pay_account`;
  END IF;

  -- 添加 approved_notes 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'approved_notes'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `approved_notes` TEXT COMMENT '审批备注' AFTER `approved_at`;
  END IF;

  -- 添加 rejected_reason 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'rejected_reason'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `rejected_reason` TEXT COMMENT '拒绝原因' AFTER `approved_notes`;
  END IF;

  -- 添加 paid_at 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'paid_at'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `paid_at` DATETIME NULL COMMENT '打款时间' AFTER `rejected_reason`;
  END IF;

  -- 添加 trade_no 字段
  IF NOT EXISTS (
    SELECT * FROM information_schema.columns
    WHERE table_schema = 'dmh'
    AND table_name = 'withdrawals'
    AND column_name = 'trade_no'
  ) THEN
    ALTER TABLE `withdrawals`
    ADD COLUMN `trade_no` VARCHAR(100) DEFAULT '' COMMENT '交易流水号' AFTER `paid_at`;
  END IF;
END$$
DELIMITER ;

CALL add_withdrawal_distributor_fields();
DROP PROCEDURE IF EXISTS add_withdrawal_distributor_fields;

-- 显示更新结果
SELECT 'Withdrawals table extended successfully!' AS status;
