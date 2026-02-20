SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

USE dmh;

-- ============================================
-- 页面配置表
-- ============================================
CREATE TABLE IF NOT EXISTS page_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    components JSON COMMENT '组件配置（JSON格式）',
    theme JSON COMMENT '主题配置（JSON格式）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',
    INDEX idx_campaign_id (campaign_id),
    UNIQUE KEY uk_campaign_id (campaign_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='页面配置表';

-- ============================================
-- 核销记录表
-- ============================================
CREATE TABLE IF NOT EXISTS verification_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '核销记录ID',
    order_id BIGINT NOT NULL COMMENT '订单ID',
    verification_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '核销状态: pending/verified/cancelled',
    verified_at DATETIME DEFAULT NULL COMMENT '核销时间',
    verified_by BIGINT DEFAULT NULL COMMENT '核销人ID',
    verification_code VARCHAR(50) DEFAULT '' COMMENT '核销码',
    verification_method VARCHAR(20) NOT NULL DEFAULT 'manual' COMMENT '核销方式: manual/auto/qrcode',
    remark VARCHAR(500) DEFAULT '' COMMENT '备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_order_id (order_id),
    INDEX idx_verification_status (verification_status),
    INDEX idx_verified_by (verified_by),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='核销记录表';

-- ============================================
-- 海报模板表
-- ============================================
CREATE TABLE IF NOT EXISTS poster_templates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '模板ID',
    type VARCHAR(20) NOT NULL COMMENT '模板类型: campaign/distributor',
    campaign_id BIGINT DEFAULT NULL COMMENT '活动ID',
    distributor_id BIGINT DEFAULT NULL COMMENT '分销商ID',
    template_url VARCHAR(500) NOT NULL COMMENT '模板URL',
    poster_data JSON COMMENT '海报数据',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_type (type),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_distributor_id (distributor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报模板表';

-- ============================================
-- 海报生成记录表
-- ============================================
CREATE TABLE IF NOT EXISTS poster_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    record_type VARCHAR(20) NOT NULL DEFAULT 'personal' COMMENT '记录类型: personal/brand',
    campaign_id BIGINT NOT NULL COMMENT '活动ID',
    distributor_id BIGINT DEFAULT 0 COMMENT '分销商ID',
    template_name VARCHAR(100) NOT NULL COMMENT '模板名称',
    poster_url VARCHAR(500) NOT NULL COMMENT '海报URL',
    thumbnail_url VARCHAR(500) DEFAULT '' COMMENT '缩略图URL',
    file_size VARCHAR(50) DEFAULT '' COMMENT '文件大小',
    generation_time INT DEFAULT 0 COMMENT '生成耗时（毫秒）',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    share_count INT DEFAULT 0 COMMENT '分享次数',
    generated_by BIGINT DEFAULT NULL COMMENT '生成人ID',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/deleted',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_record_type (record_type),
    INDEX idx_campaign_id (campaign_id),
    INDEX idx_distributor_id (distributor_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报生成记录表';

-- ============================================
-- 海报模板配置表
-- ============================================
CREATE TABLE IF NOT EXISTS poster_template_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    template_id BIGINT NOT NULL COMMENT '模板ID',
    config JSON COMMENT '配置数据',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_template_id (template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海报模板配置表';

-- ============================================
-- 会员表
-- ============================================
CREATE TABLE IF NOT EXISTS members (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '会员ID',
    unionid VARCHAR(100) NOT NULL COMMENT '微信unionid',
    nickname VARCHAR(100) DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(500) DEFAULT '' COMMENT '头像',
    phone VARCHAR(20) DEFAULT '' COMMENT '手机号',
    gender INT DEFAULT 0 COMMENT '性别: 0未知 1男 2女',
    source VARCHAR(50) DEFAULT '' COMMENT '首次来源渠道',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态: active/disabled',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY uk_unionid (unionid),
    INDEX idx_phone (phone),
    INDEX idx_source (source),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员表';

-- ============================================
-- 会员画像扩展表
-- ============================================
CREATE TABLE IF NOT EXISTS member_profiles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '画像ID',
    member_id BIGINT NOT NULL COMMENT '会员ID',
    total_orders INT DEFAULT 0 COMMENT '累计订单数',
    total_payment DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计支付金额',
    total_reward DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计奖励金额',
    first_order_at DATETIME DEFAULT NULL COMMENT '首次下单时间',
    last_order_at DATETIME DEFAULT NULL COMMENT '最近下单时间',
    first_payment_at DATETIME DEFAULT NULL COMMENT '首次支付时间',
    last_payment_at DATETIME DEFAULT NULL COMMENT '最近支付时间',
    participated_campaigns INT DEFAULT 0 COMMENT '参与活动数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_member_id (member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员画像扩展表';

-- ============================================
-- 会员标签表
-- ============================================
CREATE TABLE IF NOT EXISTS member_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    name VARCHAR(50) NOT NULL COMMENT '标签名称',
    category VARCHAR(50) DEFAULT '' COMMENT '标签分类',
    color VARCHAR(20) DEFAULT '' COMMENT '标签颜色',
    description VARCHAR(200) DEFAULT '' COMMENT '标签描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_name (name),
    INDEX idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员标签表';

-- ============================================
-- 会员标签关联表
-- ============================================
CREATE TABLE IF NOT EXISTS member_tag_links (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    member_id BIGINT NOT NULL COMMENT '会员ID',
    tag_id BIGINT NOT NULL COMMENT '标签ID',
    created_by BIGINT NOT NULL COMMENT '操作人ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_member_tag (member_id, tag_id),
    INDEX idx_member_id (member_id),
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员标签关联表';

-- ============================================
-- 会员品牌关联表
-- ============================================
CREATE TABLE IF NOT EXISTS member_brand_links (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    member_id BIGINT NOT NULL COMMENT '会员ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    first_campaign_id BIGINT NOT NULL COMMENT '首次参与活动ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_member_brand (member_id, brand_id),
    INDEX idx_member_id (member_id),
    INDEX idx_brand_id (brand_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员品牌关联表';

-- ============================================
-- 会员合并请求表
-- ============================================
CREATE TABLE IF NOT EXISTS member_merge_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '请求ID',
    source_member_id BIGINT NOT NULL COMMENT '被合并会员ID',
    target_member_id BIGINT NOT NULL COMMENT '目标会员ID',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/completed/failed',
    reason TEXT COMMENT '合并原因',
    conflict_info JSON COMMENT '冲突信息',
    created_by BIGINT NOT NULL COMMENT '操作人ID',
    executed_at DATETIME DEFAULT NULL COMMENT '执行时间',
    error_msg TEXT COMMENT '错误信息',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_source_member_id (source_member_id),
    INDEX idx_target_member_id (target_member_id),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员合并请求表';

-- ============================================
-- 导出申请表
-- ============================================
CREATE TABLE IF NOT EXISTS export_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '申请ID',
    brand_id BIGINT NOT NULL COMMENT '品牌ID',
    requested_by BIGINT NOT NULL COMMENT '申请人ID',
    reason TEXT NOT NULL COMMENT '导出原因',
    filters JSON COMMENT '筛选条件',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/approved/rejected/completed',
    approved_by BIGINT DEFAULT NULL COMMENT '审批人ID',
    approved_at DATETIME DEFAULT NULL COMMENT '审批时间',
    reject_reason TEXT COMMENT '拒绝原因',
    file_url VARCHAR(500) DEFAULT '' COMMENT '导出文件URL',
    record_count INT DEFAULT 0 COMMENT '导出记录数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_brand_id (brand_id),
    INDEX idx_requested_by (requested_by),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导出申请表';

-- ============================================
-- 插入默认海报模板
-- ============================================
INSERT INTO poster_templates (type, template_url, poster_data) VALUES
('campaign', '/templates/campaign_default.png', '{"background": "#ffffff", "width": 750, "height": 1334}'),
('distributor', '/templates/distributor_default.png', '{"background": "#ffffff", "width": 750, "height": 1334}')
ON DUPLICATE KEY UPDATE template_url = VALUES(template_url);

-- ============================================
-- 插入默认会员标签
-- ============================================
INSERT INTO member_tags (name, category, color, description) VALUES
('高价值客户', '价值', '#ff4d4f', '累计消费超过1000元'),
('活跃用户', '活跃度', '#52c41a', '近30天有活动参与'),
('新用户', '生命周期', '#1890ff', '注册7天内')
ON DUPLICATE KEY UPDATE name = VALUES(name);
