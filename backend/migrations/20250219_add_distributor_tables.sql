SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

USE dmh;

-- ============================================
-- 分销员系统相关表
-- ============================================

-- 分销商表
CREATE TABLE IF NOT EXISTS distributors (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '分销商ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    level INT NOT NULL DEFAULT 1 COMMENT '分销级别: 1/2/3',
    parent_id BIGINT DEFAULT NULL COMMENT '上级分销商ID',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/active/suspended',
    approved_by BIGINT DEFAULT NULL COMMENT '审核人ID',
    approved_at DATETIME DEFAULT NULL COMMENT '审核时间',
    total_earnings DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计收益',
    subordinates_count INT NOT NULL DEFAULT 0 COMMENT '下级数量',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',
    INDEX idx_distributor_user (user_id),
    INDEX idx_distributor_brand (brand_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status),
    INDEX idx_level (level),
    UNIQUE KEY uk_user_brand (user_id, brand_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商表';

-- 分销商申请表
CREATE TABLE IF NOT EXISTS distributor_applications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '申请ID',
    user_id BIGINT NOT NULL COMMENT '申请人用户ID',
    brand_id BIGINT NOT NULL COMMENT '申请品牌ID',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/approved/rejected',
    reason TEXT COMMENT '申请理由',
    reviewed_by BIGINT DEFAULT NULL COMMENT '审核人ID',
    reviewed_at DATETIME DEFAULT NULL COMMENT '审核时间',
    review_notes TEXT COMMENT '审核备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_brand_id (brand_id),
    INDEX idx_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商申请表';

-- 分销商级别奖励配置表
CREATE TABLE IF NOT EXISTS distributor_level_rewards (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    level INT NOT NULL COMMENT '分销级别: 1/2/3',
    reward_percentage DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '奖励百分比',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_level_reward_brand (brand_id),
    UNIQUE KEY uk_brand_level (brand_id, level),
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商级别奖励配置表';

-- 分销商奖励记录表
CREATE TABLE IF NOT EXISTS distributor_rewards (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '奖励ID',
    distributor_id BIGINT NOT NULL COMMENT '分销商ID',
    user_id BIGINT NOT NULL COMMENT '用户ID（分销商对应的用户）',
    order_id BIGINT NOT NULL COMMENT '订单ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '奖励金额',
    level INT NOT NULL COMMENT '奖励级别: 1/2/3',
    reward_rate DECIMAL(5,2) NOT NULL COMMENT '奖励比例',
    from_user_id BIGINT DEFAULT NULL COMMENT '购买用户ID',
    status VARCHAR(20) NOT NULL DEFAULT 'settled' COMMENT '状态: pending/settled/cancelled',
    settled_at DATETIME DEFAULT NULL COMMENT '结算时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_distributor_id (distributor_id),
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_status (status),
    FOREIGN KEY (distributor_id) REFERENCES distributors(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商奖励记录表';

-- 分销商推广链接表
CREATE TABLE IF NOT EXISTS distributor_links (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '链接ID',
    distributor_id BIGINT NOT NULL COMMENT '分销商ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    link_code VARCHAR(50) NOT NULL COMMENT '推广码',
    click_count INT NOT NULL DEFAULT 0 COMMENT '点击次数',
    order_count INT NOT NULL DEFAULT 0 COMMENT '订单数量',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
    expires_at DATETIME DEFAULT NULL COMMENT '过期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_distributor_id (distributor_id),
    INDEX idx_campaign_id (campaign_id),
    UNIQUE KEY uk_link_code (link_code),
    FOREIGN KEY (distributor_id) REFERENCES distributors(id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销商推广链接表';

-- 添加 distributor_id 到 orders 表（如果不存在）
SET @exist_col := (SELECT COUNT(*) FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = 'dmh' AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'distributor_id');
SET @sql := IF(@exist_col = 0, 
    'ALTER TABLE orders ADD COLUMN distributor_id BIGINT DEFAULT NULL COMMENT ''分销商ID'' AFTER referrer_id, ADD INDEX idx_distributor_id (distributor_id)', 
    'SELECT ''Column distributor_id already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ============================================
-- 添加分销员角色
-- ============================================

INSERT INTO roles (name, code, description) VALUES
('分销员', 'distributor', '品牌分销员，完成品牌活动订单支付后可获得奖励')
ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description);

-- ============================================
-- 添加分销员相关权限
-- ============================================

INSERT INTO permissions (name, code, resource, action, description) VALUES
('分销中心查看', 'distributor:read', 'distributor', 'read', '查看分销中心'),
('分销申请', 'distributor:apply', 'distributor', 'apply', '申请成为分销员'),
('推广链接管理', 'distributor:link', 'distributor', 'link', '管理推广链接'),
('收益查看', 'distributor:earnings', 'distributor', 'earnings', '查看分销收益')
ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description);

-- 分销员角色权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'distributor' 
AND p.code IN ('distributor:read', 'distributor:apply', 'distributor:link', 'distributor:earnings',
    'campaign:read', 'order:create', 'reward:read', 'withdrawal:apply')
ON DUPLICATE KEY UPDATE role_id = role_id;

-- ============================================
-- 创建分销员测试用户
-- ============================================

-- 密码都是 123456 的bcrypt加密
INSERT INTO users (username, password, phone, email, real_name, role, status) VALUES
('distributor001', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000010', 'distributor001@dmh.com', '王小明', 'participant', 'active'),
('distributor002', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13800000011', 'distributor002@dmh.com', '李小红', 'participant', 'active')
ON DUPLICATE KEY UPDATE username = VALUES(username);

-- 为分销员用户分配角色
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r 
WHERE u.username IN ('distributor001', 'distributor002') AND r.code = 'distributor'
ON DUPLICATE KEY UPDATE user_id = user_id;

-- 插入分销商记录
INSERT INTO distributors (user_id, brand_id, level, parent_id, status, approved_by, approved_at, total_earnings, subordinates_count) VALUES
((SELECT id FROM users WHERE username = 'distributor001'), 1, 1, NULL, 'active', 2, NOW(), 150.00, 1),
((SELECT id FROM users WHERE username = 'distributor002'), 1, 2, (SELECT id FROM (SELECT id FROM distributors WHERE level = 1 AND brand_id = 1 LIMIT 1) AS tmp), 'active', 2, NOW(), 50.00, 0)
ON DUPLICATE KEY UPDATE status = VALUES(status);

-- 更新上级的下级数量（使用子查询绕过 MySQL 限制）
UPDATE distributors d SET subordinates_count = (
    SELECT cnt FROM (
        SELECT id, COUNT(*) as cnt FROM distributors GROUP BY id
    ) AS tmp WHERE tmp.id = d.id
) WHERE level = 1 AND brand_id = 1;

-- 插入分销商级别奖励配置
INSERT INTO distributor_level_rewards (brand_id, level, reward_percentage) VALUES
(1, 1, 10.00),
(1, 2, 5.00),
(1, 3, 2.00)
ON DUPLICATE KEY UPDATE reward_percentage = VALUES(reward_percentage);

-- 插入用户余额
INSERT INTO user_balances (user_id, balance, total_reward)
SELECT u.id, 0.00, d.total_earnings 
FROM users u 
JOIN distributors d ON d.user_id = u.id 
WHERE u.username IN ('distributor001', 'distributor002')
ON DUPLICATE KEY UPDATE total_reward = VALUES(total_reward);

-- 插入分销推广链接示例
INSERT INTO distributor_links (distributor_id, campaign_id, link_code, click_count, order_count, status) VALUES
((SELECT d.id FROM distributors d JOIN users u ON d.user_id = u.id WHERE u.username = 'distributor001' AND d.brand_id = 1), 1, 'DIST001', 25, 3, 'active'),
((SELECT d.id FROM distributors d JOIN users u ON d.user_id = u.id WHERE u.username = 'distributor002' AND d.brand_id = 1), 1, 'DIST002', 10, 1, 'active')
ON DUPLICATE KEY UPDATE click_count = VALUES(click_count);
