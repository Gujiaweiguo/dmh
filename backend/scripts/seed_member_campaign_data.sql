-- Seed test data for campaigns and members
-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE dmh;

-- Ensure a test brand exists
INSERT INTO brands (name, logo, description, status)
SELECT '测试品牌C', 'https://via.placeholder.com/150', '会员系统测试品牌', 'active'
WHERE NOT EXISTS (SELECT 1 FROM brands WHERE name = '测试品牌C');

-- Ensure brand_manager is linked to the test brand
INSERT INTO user_brands (user_id, brand_id)
SELECT u.id, b.id
FROM users u
JOIN brands b ON b.name = '测试品牌C'
WHERE u.username = 'brand_manager'
  AND NOT EXISTS (
    SELECT 1 FROM user_brands ub
    WHERE ub.user_id = u.id AND ub.brand_id = b.id
  );

SET @brand_id := (SELECT id FROM brands WHERE name = '测试品牌C' ORDER BY id DESC LIMIT 1);
SET @admin_id := (SELECT id FROM users WHERE username = 'admin' LIMIT 1);
SET @brand_admin_id := (SELECT id FROM users WHERE username = 'brand_manager' LIMIT 1);

-- Create campaigns
INSERT INTO campaigns (brand_id, name, description, form_fields, reward_rule, start_time, end_time, status)
SELECT @brand_id, '会员招募活动-1', '会员系统测试活动A', '["姓名","手机号"]', 10.00,
       '2026-01-01 00:00:00', '2026-12-31 23:59:59', 'active'
WHERE NOT EXISTS (SELECT 1 FROM campaigns WHERE name = '会员招募活动-1');

INSERT INTO campaigns (brand_id, name, description, form_fields, reward_rule, start_time, end_time, status)
SELECT @brand_id, '会员招募活动-2', '会员系统测试活动B', '["姓名","手机号","地址"]', 5.00,
       '2026-02-01 00:00:00', '2026-08-31 23:59:59', 'active'
WHERE NOT EXISTS (SELECT 1 FROM campaigns WHERE name = '会员招募活动-2');

SET @campaign_a := (SELECT id FROM campaigns WHERE name = '会员招募活动-1' ORDER BY id DESC LIMIT 1);
SET @campaign_b := (SELECT id FROM campaigns WHERE name = '会员招募活动-2' ORDER BY id DESC LIMIT 1);

-- Member tags
INSERT INTO member_tags (name, category, color, description)
SELECT '高价值', '价值', '#f59e0b', '累计支付 >= 1000'
WHERE NOT EXISTS (SELECT 1 FROM member_tags WHERE name = '高价值');

INSERT INTO member_tags (name, category, color, description)
SELECT '活跃', '活跃度', '#10b981', '近30天有参与'
WHERE NOT EXISTS (SELECT 1 FROM member_tags WHERE name = '活跃');

SET @tag_high_value := (SELECT id FROM member_tags WHERE name = '高价值' LIMIT 1);
SET @tag_active := (SELECT id FROM member_tags WHERE name = '活跃' LIMIT 1);

-- Members
INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0001', '测试会员01', '', '13900000001', 1, 'campaign', 'active'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0001');

INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0002', '测试会员02', '', '13900000002', 2, 'campaign', 'active'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0002');

INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0003', '测试会员03', '', '13900000003', 1, 'campaign', 'active'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0003');

INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0004', '测试会员04', '', '13900000004', 0, 'campaign', 'active'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0004');

INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0005', '测试会员05', '', '13900000005', 2, 'campaign', 'disabled'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0005');

INSERT INTO members (unionid, nickname, avatar, phone, gender, source, status)
SELECT 'test_unionid_0006', '测试会员06', '', '13900000006', 1, 'campaign', 'active'
WHERE NOT EXISTS (SELECT 1 FROM members WHERE unionid = 'test_unionid_0006');

-- Member profiles
INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 2, 199.00, 10.00, '2026-01-05 10:00:00', '2026-02-15 09:00:00', 1
FROM members m
WHERE m.unionid = 'test_unionid_0001'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 1, 59.00, 5.00, '2026-02-12 14:00:00', '2026-02-12 14:00:00', 1
FROM members m
WHERE m.unionid = 'test_unionid_0002'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 3, 399.00, 20.00, '2026-01-10 11:00:00', '2026-03-01 12:00:00', 2
FROM members m
WHERE m.unionid = 'test_unionid_0003'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 0, 0.00, 0.00, NULL, NULL, 0
FROM members m
WHERE m.unionid = 'test_unionid_0004'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 1, 29.00, 0.00, '2026-02-20 16:00:00', '2026-02-20 16:00:00', 1
FROM members m
WHERE m.unionid = 'test_unionid_0005'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

INSERT INTO member_profiles (member_id, total_orders, total_payment, total_reward, first_order_at, last_order_at, participated_campaigns)
SELECT m.id, 2, 120.00, 8.00, '2026-01-25 18:00:00', '2026-02-28 09:30:00', 2
FROM members m
WHERE m.unionid = 'test_unionid_0006'
  AND NOT EXISTS (SELECT 1 FROM member_profiles mp WHERE mp.member_id = m.id);

-- Member brand links
INSERT INTO member_brand_links (member_id, brand_id, first_campaign_id)
SELECT m.id, @brand_id, @campaign_a
FROM members m
WHERE m.unionid IN ('test_unionid_0001','test_unionid_0002','test_unionid_0003','test_unionid_0004','test_unionid_0005','test_unionid_0006')
  AND NOT EXISTS (
    SELECT 1 FROM member_brand_links mbl
    WHERE mbl.member_id = m.id AND mbl.brand_id = @brand_id
  );

-- Tag links
INSERT INTO member_tag_links (member_id, tag_id, created_by)
SELECT m.id, @tag_high_value, @admin_id
FROM members m
WHERE m.unionid IN ('test_unionid_0001','test_unionid_0003')
  AND NOT EXISTS (
    SELECT 1 FROM member_tag_links mtl
    WHERE mtl.member_id = m.id AND mtl.tag_id = @tag_high_value
  );

INSERT INTO member_tag_links (member_id, tag_id, created_by)
SELECT m.id, @tag_active, @admin_id
FROM members m
WHERE m.unionid IN ('test_unionid_0002','test_unionid_0003','test_unionid_0006')
  AND NOT EXISTS (
    SELECT 1 FROM member_tag_links mtl
    WHERE mtl.member_id = m.id AND mtl.tag_id = @tag_active
  );

-- Orders
INSERT INTO orders (campaign_id, member_id, unionid, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status)
SELECT @campaign_a, m.id, m.unionid, m.phone,
       JSON_OBJECT('name', m.nickname, 'phone', m.phone),
       0, 'paid', 99.00, 'paid', 'TEST_TRADE_0001', 'synced'
FROM members m
WHERE m.unionid = 'test_unionid_0001'
  AND NOT EXISTS (SELECT 1 FROM orders o WHERE o.trade_no = 'TEST_TRADE_0001');

INSERT INTO orders (campaign_id, member_id, unionid, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status)
SELECT @campaign_a, m.id, m.unionid, m.phone,
       JSON_OBJECT('name', m.nickname, 'phone', m.phone),
       0, 'paid', 59.00, 'paid', 'TEST_TRADE_0002', 'synced'
FROM members m
WHERE m.unionid = 'test_unionid_0002'
  AND NOT EXISTS (SELECT 1 FROM orders o WHERE o.trade_no = 'TEST_TRADE_0002');

INSERT INTO orders (campaign_id, member_id, unionid, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status)
SELECT @campaign_b, m.id, m.unionid, m.phone,
       JSON_OBJECT('name', m.nickname, 'phone', m.phone),
       0, 'paid', 200.00, 'paid', 'TEST_TRADE_0003', 'synced'
FROM members m
WHERE m.unionid = 'test_unionid_0003'
  AND NOT EXISTS (SELECT 1 FROM orders o WHERE o.trade_no = 'TEST_TRADE_0003');

INSERT INTO orders (campaign_id, member_id, unionid, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status)
SELECT @campaign_b, m.id, m.unionid, m.phone,
       JSON_OBJECT('name', m.nickname, 'phone', m.phone),
       0, 'paid', 120.00, 'paid', 'TEST_TRADE_0004', 'synced'
FROM members m
WHERE m.unionid = 'test_unionid_0006'
  AND NOT EXISTS (SELECT 1 FROM orders o WHERE o.trade_no = 'TEST_TRADE_0004');

-- Rewards (linked to orders)
INSERT INTO rewards (user_id, member_id, order_id, campaign_id, amount, status, settled_at)
SELECT @brand_admin_id, o.member_id, o.id, o.campaign_id, 10.00, 'settled', NOW()
FROM orders o
WHERE o.trade_no = 'TEST_TRADE_0001'
  AND NOT EXISTS (SELECT 1 FROM rewards r WHERE r.order_id = o.id);

INSERT INTO rewards (user_id, member_id, order_id, campaign_id, amount, status, settled_at)
SELECT @brand_admin_id, o.member_id, o.id, o.campaign_id, 5.00, 'settled', NOW()
FROM orders o
WHERE o.trade_no = 'TEST_TRADE_0002'
  AND NOT EXISTS (SELECT 1 FROM rewards r WHERE r.order_id = o.id);
