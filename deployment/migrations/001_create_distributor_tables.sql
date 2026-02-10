-- Migration Script: Create distributor_rewards table
-- 创建时间: 2026-02-09
-- 说明: 创建分销奖励记录表，支持多级分销奖励系统

-- 分销商表（如果不存在则创建）
CREATE TABLE IF NOT EXISTS `distributors` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '分销商ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID（关联users表）',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID（关联brands表）',
  `name` VARCHAR(100) NOT NULL COMMENT '分销商名称',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '分销商手机号',
  `status` VARCHAR(20) DEFAULT 'pending' COMMENT '状态：pending/active/suspended',
  `total_reward` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计奖励金额',
  `withdrawable_amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '可提现金额',
  `total_orders` INT DEFAULT 0 COMMENT '累计订单数',
  `level` INT DEFAULT 1 COMMENT '分销层级：1/2/3',
  `parent_id` BIGINT DEFAULT 0 COMMENT '上级分销商ID',
  `referral_code` VARCHAR(20) DEFAULT NULL COMMENT '推荐码',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_user_brand` (`user_id`, `brand_id`),
  INDEX `idx_parent` (`parent_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_referral_code` (`referral_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销商表';

-- 分销奖励表（如果不存在则创建）
CREATE TABLE IF NOT EXISTS `distributor_rewards` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '奖励记录ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID（关联users表）',
  `distributor_id` BIGINT NOT NULL COMMENT '分销商ID（关联distributors表）',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID（关联brands表）',
  `order_id` BIGINT NOT NULL COMMENT '订单ID（关联orders表）',
  `campaign_id` BIGINT NOT NULL COMMENT '活动ID（关联campaigns表）',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '奖励金额',
  `level` INT NOT NULL DEFAULT 1 COMMENT '分销层级：1/2/3（1=一级，2=二级，3=三级）',
  `percentage` DECIMAL(5,2) DEFAULT 0.00 COMMENT '奖励比例（%）',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/settled/failed/cancelled',
  `settled_at` DATETIME DEFAULT NULL COMMENT '结算时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_user` (`user_id`),
  INDEX `idx_distributor` (`distributor_id`),
  INDEX `idx_order` (`order_id`),
  INDEX `idx_campaign` (`campaign_id`),
  INDEX `idx_level` (`level`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销奖励表';

-- 分销商层级关系表（如果不存在则创建）
CREATE TABLE IF NOT EXISTS `distributor_relations` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关系ID',
  `parent_id` BIGINT NOT NULL COMMENT '上级分销商ID',
  `child_id` BIGINT NOT NULL COMMENT '下级分销商ID',
  `level` INT NOT NULL DEFAULT 1 COMMENT '层级关系：1=父子，2=祖孙',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  UNIQUE KEY `uk_parent_child` (`parent_id`, `child_id`),
  INDEX `idx_parent` (`parent_id`),
  INDEX `idx_child` (`child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销商层级关系表';

-- 分销申请表（如果不存在则创建）
CREATE TABLE IF NOT EXISTS `distributor_applications` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '申请ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID（关联users表）',
  `brand_id` BIGINT NOT NULL COMMENT '品牌ID（关联brands表）',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '申请提现金额',
  `bank_name` VARCHAR(100) DEFAULT NULL COMMENT '银行名称',
  `bank_account` VARCHAR(50) DEFAULT NULL COMMENT '银行账号',
  `account_name` VARCHAR(100) DEFAULT NULL COMMENT '账户名称',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/approved/rejected/processing/completed',
  `approved_by` BIGINT DEFAULT NULL COMMENT '审批人ID（关联users表）',
  `approved_at` DATETIME DEFAULT NULL COMMENT '审批时间',
  `paid_at` DATETIME DEFAULT NULL COMMENT '付款时间',
  `rejection_reason` TEXT DEFAULT NULL COMMENT '拒绝原因',
  `remark` TEXT DEFAULT NULL COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_user` (`user_id`),
  INDEX `idx_brand` (`brand_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销商提现申请表';

-- 插入测试数据（可选）
-- 插入一些测试分销商
INSERT INTO `distributors` (`user_id`, `brand_id`, `name`, `status`, `level`)
VALUES
  (2, 1, '品牌经理分销商', 'active', 1),
  (3, 1, '用户001分销商', 'active', 2)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- 注释：这些测试数据关联到已有的用户（brand_manager和user001）
-- 可以根据需要调整或删除

-- Migration完成提示
SELECT 'Migration completed: distributor tables created' AS status;
