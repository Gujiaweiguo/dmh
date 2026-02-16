-- Migration: Fix distributors table columns to match model
-- Date: 2026-02-17
-- Issue: Model expects total_earnings, subordinates_count, approved_by, approved_at, deleted_at
--         but database has total_reward, total_orders instead

-- Add missing columns
ALTER TABLE `distributors`
ADD COLUMN `total_earnings` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计收益' AFTER `status`,
ADD COLUMN `subordinates_count` INT NOT NULL DEFAULT 0 COMMENT '下级人数' AFTER `total_earnings`,
ADD COLUMN `approved_by` BIGINT NULL COMMENT '审批人ID' AFTER `subordinates_count`,
ADD COLUMN `approved_at` DATETIME NULL COMMENT '审批时间' AFTER `approved_by`,
ADD COLUMN `deleted_at` DATETIME NULL COMMENT '软删除时间' AFTER `updated_at`;

-- Migrate data from total_reward to total_earnings
UPDATE `distributors` SET `total_earnings` = COALESCE(`total_reward`, 0.00) WHERE `total_reward` IS NOT NULL;

-- Add indexes for new columns
CREATE INDEX `idx_distributors_approved_at` ON `distributors` (`approved_at`);
CREATE INDEX `idx_distributors_deleted_at` ON `distributors` (`deleted_at`);
