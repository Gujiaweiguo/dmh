-- Migration: create_promoters_materials_tables
-- Date: 2026-02-24
-- Description: Create promoters and materials tables for integration tests

-- Promoters table
CREATE TABLE IF NOT EXISTS `promoters` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `brand_id` BIGINT NOT NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'active',
    `level` VARCHAR(50) DEFAULT '',
    `total_orders` BIGINT NOT NULL DEFAULT 0,
    `total_rewards` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    `conversion_rate` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    `campaign_count` BIGINT NOT NULL DEFAULT 0,
    `last_active_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME NULL,
    INDEX `idx_promoter_user` (`user_id`),
    INDEX `idx_promoter_brand` (`brand_id`),
    UNIQUE INDEX `idx_user_brand` (`user_id`, `brand_id`),
    INDEX `idx_status` (`status`),
    CONSTRAINT `fk_promoters_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_promoters_brand` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Promoter links table
CREATE TABLE IF NOT EXISTS `promoter_links` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `promoter_id` BIGINT NOT NULL,
    `campaign_id` BIGINT NOT NULL,
    `link_code` VARCHAR(50) NOT NULL,
    `click_count` BIGINT NOT NULL DEFAULT 0,
    `order_count` BIGINT NOT NULL DEFAULT 0,
    `status` VARCHAR(20) NOT NULL DEFAULT 'active',
    `expires_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_promoter_id` (`promoter_id`),
    INDEX `idx_campaign_id` (`campaign_id`),
    UNIQUE INDEX `idx_link_code` (`link_code`),
    CONSTRAINT `fk_promoter_links_promoter` FOREIGN KEY (`promoter_id`) REFERENCES `promoters` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_promoter_links_campaign` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Promoter rewards table
CREATE TABLE IF NOT EXISTS `promoter_rewards` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `promoter_id` BIGINT NOT NULL,
    `type` VARCHAR(20) NOT NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending',
    `amount` DECIMAL(10,2) NOT NULL,
    `description` TEXT,
    `campaign_id` BIGINT NULL,
    `order_id` BIGINT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_promoter_id` (`promoter_id`),
    INDEX `idx_type` (`type`),
    INDEX `idx_status` (`status`),
    INDEX `idx_campaign_id` (`campaign_id`),
    INDEX `idx_order_id` (`order_id`),
    CONSTRAINT `fk_promoter_rewards_promoter` FOREIGN KEY (`promoter_id`) REFERENCES `promoters` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_promoter_rewards_campaign` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_promoter_rewards_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Materials table
CREATE TABLE IF NOT EXISTS `materials` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `type` VARCHAR(20) NOT NULL DEFAULT 'image',
    `url` VARCHAR(500),
    `content` TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME NULL,
    INDEX `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert test data
INSERT INTO `promoters` (`user_id`, `brand_id`, `status`, `level`, `total_orders`, `total_rewards`) VALUES
(2, 1, 'active', 'VIP', 10, 1000.00);

INSERT INTO `materials` (`name`, `type`, `url`, `content`) VALUES
('test_image.png', 'image', '/uploads/test_image.png', NULL),
('test_text.txt', 'text', NULL, 'Sample text content for testing');
