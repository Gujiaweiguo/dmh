-- 测试数据脚本
-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 插入测试品牌
INSERT INTO brands (id, name, logo, description, status, created_at, updated_at) VALUES
(1, '星巴克咖啡', 'https://example.com/starbucks.png', '全球知名咖啡连锁品牌', 'active', NOW(), NOW()),
(2, '麦当劳', 'https://example.com/mcdonalds.png', '全球快餐连锁品牌', 'active', NOW(), NOW()),
(3, '耐克运动', 'https://example.com/nike.png', '全球运动品牌领导者', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入测试活动
INSERT INTO campaigns (id, name, description, form_fields, reward_rule, start_time, end_time, status, created_at, updated_at) VALUES
(1, '新年促销活动', '2026年新年特惠活动，报名即送优惠券', '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true},{"name":"course","label":"意向课程","type":"select","options":["Java开发","Python开发","前端开发"],"required":true}]', 10.00, '2026-01-01 00:00:00', '2026-03-31 23:59:59', 'active', NOW(), NOW()),
(2, '春季招生活动', '春季班火热招生中，推荐有礼', '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true},{"name":"age","label":"年龄","type":"number","required":false}]', 15.00, '2026-02-01 00:00:00', '2026-05-31 23:59:59', 'active', NOW(), NOW()),
(3, '会员专享活动', '老会员专享福利活动', '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true},{"name":"member_id","label":"会员号","type":"text","required":true}]', 20.00, '2026-01-15 00:00:00', '2026-06-30 23:59:59', 'active', NOW(), NOW()),
(4, '暑期特训营', '暑期编程特训营，限时优惠', '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true}]', 25.00, '2026-06-01 00:00:00', '2026-08-31 23:59:59', 'draft', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入测试订单
INSERT INTO orders (id, campaign_id, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status, created_at, updated_at) VALUES
(1, 1, '13900000001', '{"name":"张三","phone":"13900000001","course":"Java开发"}', 2, 'completed', 99.00, 'paid', 'TN20260115001', 'synced', NOW(), NOW()),
(2, 1, '13900000002', '{"name":"李四","phone":"13900000002","course":"Python开发"}', 2, 'completed', 99.00, 'paid', 'TN20260115002', 'synced', NOW(), NOW()),
(3, 2, '13900000003', '{"name":"王五","phone":"13900000003","age":25}', 3, 'completed', 199.00, 'paid', 'TN20260115003', 'pending', NOW(), NOW()),
(4, 1, '13900000004', '{"name":"赵六","phone":"13900000004","course":"前端开发"}', 0, 'pending', 99.00, 'unpaid', '', 'pending', NOW(), NOW()),
(5, 3, '13900000005', '{"name":"钱七","phone":"13900000005","member_id":"VIP001"}', 2, 'completed', 299.00, 'paid', 'TN20260115005', 'synced', NOW(), NOW())
ON DUPLICATE KEY UPDATE phone=VALUES(phone);

-- 插入测试奖励记录
INSERT INTO rewards (id, user_id, order_id, campaign_id, amount, status, settled_at, created_at, updated_at) VALUES
(1, 2, 1, 1, 10.00, 'settled', NOW(), NOW(), NOW()),
(2, 2, 2, 1, 10.00, 'settled', NOW(), NOW(), NOW()),
(3, 3, 3, 2, 15.00, 'settled', NOW(), NOW(), NOW()),
(4, 2, 5, 3, 20.00, 'settled', NOW(), NOW(), NOW())
ON DUPLICATE KEY UPDATE amount=VALUES(amount);

-- 插入用户余额
INSERT INTO user_balances (id, user_id, balance, total_reward, version, created_at, updated_at) VALUES
(1, 2, 40.00, 40.00, 4, NOW(), NOW()),
(2, 3, 15.00, 15.00, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE balance=VALUES(balance);
