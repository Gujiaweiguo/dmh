-- ============================================
-- 分销商系统测试数据脚本
-- 包含分销商相关的完整测试数据
-- ============================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE dmh;

-- ============================================
-- 1. 插入分销商角色（如果不存在）
-- ============================================
INSERT IGNORE INTO roles (name, code, description) VALUES
('分销商', 'distributor', '拥有推广和分销权限的高级用户');

-- 获取分销商角色ID
SET @distributor_role_id := (SELECT id FROM roles WHERE code = 'distributor' LIMIT 1);

-- ============================================
-- 2. 创建测试用户
-- 密码都是 123456 的bcrypt加密
-- ============================================
INSERT IGNORE INTO users (id, username, password, phone, email, real_name, role, status) VALUES
(4, 'distributor1', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000001', 'distributor1@test.com', '一级分销商', 'distributor', 'active'),
(5, 'distributor2', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000002', 'distributor2@test.com', '二级分销商', 'distributor', 'active'),
(6, 'distributor3', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000003', 'distributor3@test.com', '三级分销商', 'distributor', 'active'),
(7, 'participant1', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000004', 'participant1@test.com', '普通用户1', 'participant', 'active'),
(8, 'participant2', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000005', 'participant2@test.com', '普通用户2', 'participant', 'active'),
(9, 'distributor4', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000006', 'distributor4@test.com', '品牌B一级分销商', 'distributor', 'active'),
(10, 'distributor5', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000007', 'distributor5@test.com', '品牌B二级分销商', 'distributor', 'active'),
(11, 'distributor6', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000008', 'distributor6@test.com', '品牌B三级分销商', 'distributor', 'active'),
(12, 'participant3', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000009', 'participant3@test.com', '普通用户3', 'participant', 'active'),
(13, 'participant4', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000010', 'participant4@test.com', '普通用户4', 'participant', 'active'),
(14, 'distributor_pending', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000011', 'pending@test.com', '待审核分销商', 'distributor', 'active'),
(15, 'distributor_suspended', '$2a$10$iL5hmpD0wGKSkRDCY92TL.y8wGarBWmnqVoFYlRxLM7xr0eSCzPEm', '13900000012', 'suspended@test.com', '已暂停分销商', 'distributor', 'active');

-- ============================================
-- 3. 分配角色给用户
-- ============================================
INSERT IGNORE INTO user_roles (user_id, role_id) VALUES
(4, @distributor_role_id),
(5, @distributor_role_id),
(6, @distributor_role_id),
(7, 3), -- participant1 -> participant
(8, 3), -- participant2 -> participant
(9, @distributor_role_id),
(10, @distributor_role_id),
(11, @distributor_role_id),
(12, 3),
(13, 3),
(14, @distributor_role_id),
(15, @distributor_role_id);

-- ============================================
-- 4. 将用户关联到品牌
-- ============================================
INSERT IGNORE INTO user_brands (user_id, brand_id) VALUES
(4, 1), -- distributor1 -> 品牌A
(5, 1), -- distributor2 -> 品牌A
(6, 1), -- distributor3 -> 品牌A
(7, 1), -- participant1 -> 品牌A
(8, 1), -- participant2 -> 品牌A
(9, 2), -- distributor4 -> 品牌B
(10, 2), -- distributor5 -> 品牌B
(11, 2), -- distributor6 -> 品牌B
(12, 2), -- participant3 -> 品牌B
(13, 2), -- participant4 -> 品牌B
(14, 1), -- distributor_pending -> 品牌A
(15, 1); -- distributor_suspended -> 品牌A

-- ============================================
-- 5. 创建分销商记录（三级分销链）
-- ============================================
INSERT IGNORE INTO distributors (id, user_id, brand_id, level, parent_id, status, total_earnings, subordinates_count) VALUES
(1, 4, 1, 1, NULL, 'active', 500.00, 2), -- 一级分销商（distributor1），有2个下级
(2, 5, 1, 2, 1, 'active', 200.00, 1), -- 二级分销商（distributor2），下级是distributor1
(3, 6, 1, 3, 2, 'active', 100.00, 0), -- 三级分销商（distributor3），下级是distributor2
(101, 9, 2, 1, NULL, 'active', 300.00, 2), -- 品牌B一级分销商
(102, 10, 2, 2, 101, 'pending', 80.00, 1), -- 品牌B待审核分销商
(103, 11, 2, 3, 102, 'suspended', 20.00, 0), -- 品牌B暂停分销商
(104, 14, 1, 1, NULL, 'pending', 0.00, 0), -- 品牌A待审核分销商
(105, 15, 1, 2, 1, 'suspended', 0.00, 0); -- 品牌A暂停分销商

-- ============================================
-- 6. 创建用户余额记录
-- ============================================
INSERT IGNORE INTO user_balances (user_id, balance, total_reward, version) VALUES
(4, 500.00, 500.00, 0), -- distributor1 累计收益500
(5, 200.00, 200.00, 0), -- distributor2 累计收益200
(6, 100.00, 100.00, 0), -- distributor3 累计收益100
(7, 0.00, 0.00, 0), -- participant1 普通用户
(8, 0.00, 0.00, 0), -- participant2 普通用户
(9, 300.00, 300.00, 0), -- distributor4 累计收益300
(10, 80.00, 80.00, 0), -- distributor5 累计收益80
(11, 20.00, 20.00, 0), -- distributor6 累计收益20
(12, 0.00, 0.00, 0), -- participant3 普通用户
(13, 0.00, 0.00, 0), -- participant4 普通用户
(14, 0.00, 0.00, 0), -- distributor_pending
(15, 0.00, 0.00, 0); -- distributor_suspended

-- ============================================
-- 7. 更新活动，启用分销并设置分销规则
-- ============================================
-- 修复历史导入造成的活动名称乱码
UPDATE campaigns
SET name = '分销测试活动',
    description = '用于测试分销系统的活动'
WHERE name = 'åˆ†é”€æµ‹è¯•æ´»åŠ¨';
-- 查找活动ID
SET @campaign_id := (SELECT id FROM campaigns WHERE name = '新年促销活动' LIMIT 1);

-- 更新活动，启用分销
UPDATE campaigns
SET enable_distribution = TRUE,
    distribution_level = 3,
    distribution_rewards = '{"1": 10.0, "2": 5.0, "3": 3.0}'
WHERE id = @campaign_id;

-- 如果没有找到活动，创建一个带分销规则的测试活动
INSERT IGNORE INTO campaigns (
    brand_id, name, description, form_fields,
    reward_rule, start_time, end_time,
    enable_distribution, distribution_level, distribution_rewards, status
) VALUES (
    1, '分销测试活动', '用于测试分销系统的活动',
    '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true}]',
    0.00, '2026-01-01 00:00:00', '2026-12-31 23:59:59',
    TRUE, 3, '{"1": 10.0, "2": 5.0, "3": 3.0}', 'active'
), (
    2, '分销测试活动B', '品牌B分销测试活动',
    '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true}]',
    0.00, '2026-01-01 00:00:00', '2026-12-31 23:59:59',
    TRUE, 2, '{"1": 8.0, "2": 4.0}', 'active'
);

-- 查找品牌B活动ID
SET @campaign_id_b := (SELECT id FROM campaigns WHERE name = '分销测试活动B' LIMIT 1);

-- 品牌B的订单数据
INSERT IGNORE INTO orders (
    id, campaign_id, phone, form_data, referrer_id,
    distributor_path, status, pay_status, amount, created_at, paid_at
) VALUES
-- 订单5：品牌B一级分销（distributor4）
(20, @campaign_id_b, '13900000020', '{"name":"测试用户20","phone":"13900000020"}', 9,
IF(@campaign_id_b IS NOT NULL, '101', ''), 'paid', 'paid', 120.00, NOW(), NOW()),

-- 订单6：品牌B二级分销（distributor4 -> distributor5）
(21, @campaign_id_b, '13900000021', '{"name":"测试用户21","phone":"13900000021"}', 10,
IF(@campaign_id_b IS NOT NULL, '101,102', ''), 'paid', 'paid', 220.00, NOW(), NOW()),

-- 订单7：品牌B三级分销（distributor4 -> distributor5 -> distributor6）
(22, @campaign_id_b, '13900000022', '{"name":"测试用户22","phone":"13900000022"}', 11,
IF(@campaign_id_b IS NOT NULL, '101,102,103', ''), 'paid', 'paid', 180.00, NOW(), NOW());

-- ============================================
-- 8. 创建订单数据（带分销路径）
-- ============================================
INSERT IGNORE INTO orders (
    id, campaign_id, phone, form_data, referrer_id,
    distributor_path, status, pay_status, amount, created_at, paid_at
) VALUES
-- 订单1：通过distributor1推荐，路径为 "1"（一级）
(10, @campaign_id, '13900000010', '{"name":"测试用户10","phone":"13900000010"}', 4, '1', 'paid', 'paid', 100.00, NOW(), NOW()),

-- 订单2：通过distributor3推荐，路径为 "1,2,3"（三级）
(11, @campaign_id, '13900000011', '{"name":"测试用户11","phone":"13900000011"}', 6, '1,2,3', 'paid', 'paid', 200.00, NOW(), NOW()),

-- 订单3：通过distributor2推荐，路径为 "1,2"（二级）
(12, @campaign_id, '13900000012', '{"name":"测试用户12","phone":"13900000012"}', 5, '1,2', 'paid', 'paid', 150.00, NOW(), NOW()),

-- 订单4：普通用户自己下单，无推荐
(13, @campaign_id, '13900000013', '{"name":"测试用户13","phone":"13900000013"}', 0, '', 'paid', 'paid', 80.00, NOW(), NOW());

-- ============================================
-- 9. 创建分销商奖励记录
-- ============================================
-- 订单10（100元，一级分销）
-- distributor1获得：100 * 10% = 10元
INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
) VALUES
(1, 4, 10, @campaign_id, 10.00, 1, 10.0, 'settled', NOW());

-- 订单11（200元，三级分销）
-- distributor1获得：200 * 10% = 20元
-- distributor2获得：200 * 5% = 10元
-- distributor3获得：200 * 3% = 6元
INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
) VALUES
(1, 4, 11, @campaign_id, 20.00, 1, 10.0, 'settled', NOW()),
(2, 5, 11, @campaign_id, 10.00, 2, 5.0, 'settled', NOW()),
(3, 6, 11, @campaign_id, 6.00, 3, 3.0, 'settled', NOW());

-- 订单12（150元，二级分销）
-- distributor1获得：150 * 10% = 15元
-- distributor2获得：150 * 5% = 7.5元
INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
) VALUES
(1, 4, 12, @campaign_id, 15.00, 1, 10.0, 'settled', NOW()),
(2, 5, 12, @campaign_id, 7.50, 2, 5.0, 'settled', NOW());

-- 品牌B奖励记录（待结算/已取消）
INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
)
SELECT 101, 9, 20, @campaign_id_b, 9.60, 1, 8.0, 'pending', NULL
WHERE @campaign_id_b IS NOT NULL;

INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
)
SELECT 101, 9, 21, @campaign_id_b, 17.60, 1, 8.0, 'settled', NOW()
WHERE @campaign_id_b IS NOT NULL;

INSERT IGNORE INTO distributor_rewards (
    distributor_id, user_id, order_id, campaign_id,
    amount, distributor_level, reward_percentage, status, settled_at
)
SELECT 102, 10, 21, @campaign_id_b, 8.80, 2, 4.0, 'cancelled', NULL
WHERE @campaign_id_b IS NOT NULL;

-- ============================================
-- 10. 创建提现申请数据
-- ============================================
INSERT IGNORE INTO withdrawals (
    id, user_id, brand_id, distributor_id, amount, status,
    pay_type, pay_account, pay_real_name,
    created_at, updated_at
) VALUES
-- 待审核提现
(1, 4, 1, 1, 50.00, 'pending',
 'wechat', 'wx_13900000001', '一级分销商',
 NOW(), NOW()),

-- 已通过提现
(2, 4, 1, 1, 100.00, 'approved',
 'wechat', 'wx_13900000001', '一级分销商',
 NOW(), NOW()),

-- 已拒绝提现
(3, 5, 1, 2, 300.00, 'rejected',
 'alipay', 'ali_13900000002', '二级分销商',
 NOW(), NOW()),

-- distributor2的待审核提现
(4, 5, 1, 2, 50.00, 'pending',
 'alipay', 'ali_13900000002', '二级分销商',
 NOW(), NOW());

-- 更新已通过和已拒绝提现的详细信息
UPDATE withdrawals
SET approved_by = 1,
    approved_at = NOW()
WHERE id = 2;

UPDATE withdrawals
SET approved_by = 1,
    approved_at = NOW(),
    rejected_reason = '金额超过可提现余额'
WHERE id = 3;

-- 处理中提现
INSERT IGNORE INTO withdrawals (
    id, user_id, brand_id, distributor_id, amount, status,
    pay_type, pay_account, pay_real_name,
    created_at, updated_at
) VALUES
(5, 9, 2, 101, 60.00, 'processing',
 'wechat', 'wx_13900000006', '品牌B一级分销商',
 NOW(), NOW()),

-- 已完成提现
(6, 10, 2, 102, 40.00, 'completed',
 'alipay', 'ali_13900000007', '品牌B二级分销商',
 NOW(), NOW());

-- 更新已完成提现的详细信息
UPDATE withdrawals
SET approved_by = 1,
    approved_at = NOW(),
    paid_at = NOW(),
    trade_no = 'TXN_20260120001',
    approved_notes = '系统测试打款'
WHERE id = 6;

-- ============================================
-- 11. 更新分销商累计收益（计算实际奖励总和）
-- ============================================
UPDATE distributors d
SET total_earnings = (
    SELECT COALESCE(SUM(amount), 0)
    FROM distributor_rewards dr
    WHERE dr.distributor_id = d.id
)
WHERE d.id IN (1, 2, 3, 101, 102, 103, 104, 105);

-- ============================================
-- 12. 创建推广链接数据
-- ============================================
INSERT IGNORE INTO distributor_links (
    distributor_id, campaign_id, link_code,
    click_count, order_count, status,
    created_at, updated_at
) VALUES
(1, @campaign_id, 'LINK_D1_C1', 100, 10, 'active', NOW(), NOW()),
(2, @campaign_id, 'LINK_D2_C1', 50, 5, 'active', NOW(), NOW()),
(3, @campaign_id, 'LINK_D3_C1', 30, 3, 'active', NOW(), NOW()),
(101, @campaign_id_b, 'LINK_D4_C2', 80, 8, 'active', NOW(), NOW()),
(102, @campaign_id_b, 'LINK_D5_C2', 20, 2, 'inactive', NOW(), NOW()),
(103, @campaign_id_b, 'LINK_D6_C2', 10, 1, 'expired', NOW(), NOW());

-- ============================================
-- 13. 创建分销商申请数据（注意：该表在最新迁移中已被移除）
-- ============================================
-- 注意：distributor_applications 表已在 20250120_update_distributor_system.sql 中删除
-- 如需测试申请审批功能，请先恢复该表或使用其他测试方式

SELECT '✓ 分销商系统测试数据导入完成' AS result;

-- ============================================
-- 数据汇总
-- ============================================
SELECT '========================================' AS '';
SELECT '测试账号信息汇总' AS '';
SELECT '========================================' AS '';

SELECT
    '平台管理员' AS role,
    username,
    phone,
    real_name,
    '123456' AS password
FROM users WHERE id = 1

UNION ALL

SELECT
    '品牌管理员',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 2

UNION ALL

SELECT
    '一级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 4

UNION ALL

SELECT
    '二级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 5

UNION ALL

SELECT
    '三级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 6

UNION ALL

SELECT
    '普通用户',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 7

UNION ALL

SELECT
    '品牌B一级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 9

UNION ALL

SELECT
    '品牌B二级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 10

UNION ALL

SELECT
    '品牌B三级分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 11

UNION ALL

SELECT
    '待审核分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 14

UNION ALL

SELECT
    '已暂停分销商',
    username,
    phone,
    real_name,
    '123456'
FROM users WHERE id = 15;
