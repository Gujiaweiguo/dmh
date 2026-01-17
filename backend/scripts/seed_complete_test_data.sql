-- 完整测试数据脚本
-- 包含所有场景的测试数据
-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE dmh;

-- ============================================
-- 1. 品牌数据（补充）
-- ============================================
INSERT INTO brands (name, logo, description, status) VALUES
('华为科技', 'https://via.placeholder.com/150', '全球领先的ICT解决方案提供商', 'active'),
('小米科技', 'https://via.placeholder.com/150', '智能手机和智能硬件公司', 'active'),
('禁用品牌', 'https://via.placeholder.com/150', '已禁用的测试品牌', 'disabled')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 2. 品牌素材数据
-- ============================================
SET @brand_id_1 := (SELECT id FROM brands WHERE name = '星巴克咖啡' LIMIT 1);
SET @brand_id_2 := (SELECT id FROM brands WHERE name = '麦当劳' LIMIT 1);

INSERT INTO brand_assets (brand_id, name, type, category, tags, file_url, file_size, description, status) VALUES
(@brand_id_1, '品牌Logo', 'image', '品牌标识', 'logo,品牌', 'https://via.placeholder.com/200', 10240, '星巴克官方Logo', 'active'),
(@brand_id_1, '新年海报', 'image', '活动海报', '新年,促销', 'https://via.placeholder.com/800x600', 102400, '2026新年促销海报', 'active'),
(@brand_id_1, '产品视频', 'video', '产品展示', '咖啡,产品', 'https://example.com/video.mp4', 5242880, '咖啡制作过程视频', 'active'),
(@brand_id_2, '品牌Logo', 'image', '品牌标识', 'logo,品牌', 'https://via.placeholder.com/200', 10240, '麦当劳官方Logo', 'active'),
(@brand_id_2, '菜单图片', 'image', '产品展示', '菜单,产品', 'https://via.placeholder.com/600x400', 51200, '最新菜单图片', 'active')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 3. 活动数据（补充不同状态）
-- ============================================
INSERT INTO campaigns (brand_id, name, description, form_fields, reward_rule, start_time, end_time, status) VALUES
-- 已结束的活动
(1, '双十一促销', '双十一限时优惠', '[{"name":"name","label":"姓名","type":"text","required":true},{"name":"phone","label":"手机号","type":"phone","required":true}]', 15.00, '2025-11-01 00:00:00', '2025-11-11 23:59:59', 'ended'),
-- 暂停的活动
(1, '暂停测试活动', '暂停状态的活动', '[{"name":"name","label":"姓名","type":"text","required":true}]', 5.00, '2026-01-01 00:00:00', '2026-12-31 23:59:59', 'paused'),
-- 草稿活动
(2, '待发布活动', '草稿状态的活动', '[{"name":"name","label":"姓名","type":"text","required":true}]', 8.00, '2026-06-01 00:00:00', '2026-08-31 23:59:59', 'draft')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 4. 订单数据（补充不同状态）
-- ============================================
INSERT INTO orders (campaign_id, phone, form_data, referrer_id, status, amount, pay_status, trade_no, sync_status) VALUES
-- 待支付订单
(1, '13800000010', '{"name":"测试用户A","phone":"13800000010"}', 0, 'pending', 99.00, 'unpaid', '', 'pending'),
(2, '13800000011', '{"name":"测试用户B","phone":"13800000011"}', 2, 'pending', 199.00, 'unpaid', '', 'pending'),
-- 已取消订单
(1, '13800000012', '{"name":"测试用户C","phone":"13800000012"}', 0, 'cancelled', 99.00, 'refunded', 'TN_CANCEL_001', 'synced'),
-- 同步失败订单
(2, '13800000013', '{"name":"测试用户D","phone":"13800000013"}', 3, 'completed', 199.00, 'paid', 'TN_SYNC_FAIL_001', 'failed'),
-- 待同步订单
(3, '13800000014', '{"name":"测试用户E","phone":"13800000014"}', 2, 'completed', 299.00, 'paid', 'TN_PENDING_001', 'pending')
ON DUPLICATE KEY UPDATE phone=VALUES(phone);

-- ============================================
-- 5. 奖励数据（补充不同状态）
-- ============================================
SET @order_pending := (SELECT id FROM orders WHERE trade_no = 'TN_PENDING_001' LIMIT 1);

INSERT INTO rewards (user_id, order_id, campaign_id, amount, status, settled_at) VALUES
-- 待结算奖励
(2, COALESCE(@order_pending, 1), 3, 20.00, 'pending', NULL),
-- 已取消奖励
(3, 1, 1, 10.00, 'cancelled', NULL)
ON DUPLICATE KEY UPDATE amount=VALUES(amount);

-- ============================================
-- 6. 提现申请数据
-- ============================================
INSERT INTO withdrawals (user_id, amount, bank_name, bank_account, account_name, status, remark, approved_by, approved_at) VALUES
-- 待审核提现
(2, 50.00, '中国工商银行', '6222021234567890123', '张三', 'pending', '首次提现申请', NULL, NULL),
(3, 30.00, '中国建设银行', '6227001234567890123', '李四', 'pending', '奖励提现', NULL, NULL),
-- 已通过提现
(2, 100.00, '中国工商银行', '6222021234567890123', '张三', 'approved', '审核通过', 1, NOW()),
-- 已拒绝提现
(3, 500.00, '中国建设银行', '6227001234567890123', '李四', 'rejected', '金额超过可提现余额', 1, NOW())
ON DUPLICATE KEY UPDATE amount=VALUES(amount);

-- ============================================
-- 7. 同步日志数据
-- ============================================
INSERT INTO sync_logs (order_id, sync_type, sync_status, attempts, error_msg, synced_at) VALUES
(1, 'order', 'synced', 1, NULL, NOW()),
(2, 'order', 'synced', 1, NULL, NOW()),
(3, 'order', 'failed', 3, '连接超时：无法连接到外部数据库', NULL),
(4, 'order', 'pending', 0, NULL, NULL),
(1, 'reward', 'synced', 1, NULL, NOW()),
(2, 'reward', 'synced', 1, NULL, NOW())
ON DUPLICATE KEY UPDATE sync_status=VALUES(sync_status);

-- ============================================
-- 8. 审计日志数据
-- ============================================
INSERT INTO audit_logs (user_id, username, action, resource, resource_id, details, client_ip, user_agent, status) VALUES
(1, 'admin', 'login', 'auth', '', '管理员登录系统', '192.168.1.100', 'Mozilla/5.0', 'success'),
(1, 'admin', 'create', 'user', '3', '创建用户: user001', '192.168.1.100', 'Mozilla/5.0', 'success'),
(1, 'admin', 'update', 'brand', '1', '更新品牌信息', '192.168.1.100', 'Mozilla/5.0', 'success'),
(2, 'brand_manager', 'login', 'auth', '', '品牌管理员登录', '192.168.1.101', 'Mozilla/5.0', 'success'),
(2, 'brand_manager', 'create', 'campaign', '1', '创建活动: 新年促销活动', '192.168.1.101', 'Mozilla/5.0', 'success'),
(2, 'brand_manager', 'update', 'campaign', '1', '更新活动状态', '192.168.1.101', 'Mozilla/5.0', 'success'),
(NULL, 'unknown', 'login', 'auth', '', '登录失败: 用户名或密码错误', '192.168.1.200', 'Mozilla/5.0', 'failed')
ON DUPLICATE KEY UPDATE action=VALUES(action);

-- ============================================
-- 9. 登录尝试记录
-- ============================================
INSERT INTO login_attempts (user_id, username, client_ip, user_agent, success, fail_reason) VALUES
(1, 'admin', '192.168.1.100', 'Mozilla/5.0', TRUE, ''),
(2, 'brand_manager', '192.168.1.101', 'Mozilla/5.0', TRUE, ''),
(NULL, 'hacker', '192.168.1.200', 'curl/7.68.0', FALSE, '用户不存在'),
(NULL, 'admin', '192.168.1.200', 'curl/7.68.0', FALSE, '密码错误'),
(NULL, 'admin', '192.168.1.200', 'curl/7.68.0', FALSE, '密码错误'),
(NULL, 'admin', '192.168.1.200', 'curl/7.68.0', FALSE, '密码错误')
ON DUPLICATE KEY UPDATE success=VALUES(success);

-- ============================================
-- 10. 安全事件记录
-- ============================================
INSERT INTO security_events (event_type, severity, user_id, username, client_ip, user_agent, description, details, handled, handled_by, handled_at) VALUES
('brute_force', 'high', NULL, 'admin', '192.168.1.200', 'curl/7.68.0', '检测到暴力破解尝试', '{"attempts": 5, "time_window": "5分钟"}', TRUE, 1, NOW()),
('suspicious_login', 'medium', 2, 'brand_manager', '10.0.0.100', 'Mozilla/5.0', '异常地点登录', '{"location": "未知地区", "usual_location": "北京"}', FALSE, NULL, NULL),
('password_expired', 'low', 3, 'user001', '192.168.1.102', 'Mozilla/5.0', '密码即将过期', '{"days_remaining": 7}', FALSE, NULL, NULL)
ON DUPLICATE KEY UPDATE event_type=VALUES(event_type);

-- ============================================
-- 11. 会员标签（补充）
-- ============================================
INSERT INTO member_tags (name, category, color, description) VALUES
('新会员', '生命周期', '#3b82f6', '注册30天内'),
('沉睡会员', '活跃度', '#6b7280', '超过90天未活跃'),
('VIP', '等级', '#eab308', 'VIP会员'),
('潜在流失', '风险', '#ef4444', '有流失风险的会员')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 12. 页面配置数据
-- ============================================
SET @campaign_1 := (SELECT id FROM campaigns WHERE name = '新年促销活动' LIMIT 1);
SET @campaign_2 := (SELECT id FROM campaigns WHERE name = '春季招生活动' LIMIT 1);

INSERT INTO page_configs (campaign_id, components, theme) VALUES
(@campaign_1, '[{"type":"poster","props":{"imageUrl":"https://via.placeholder.com/750x400","alt":"新年促销"}},{"type":"title","props":{"text":"2026新年大促","fontSize":24,"color":"#ff4d4f"}},{"type":"detail","props":{"content":"新年特惠，报名即送优惠券！"}},{"type":"button","props":{"text":"立即报名","color":"#ff4d4f"}}]', '{"primaryColor":"#ff4d4f","backgroundColor":"#fff5f5"}'),
(@campaign_2, '[{"type":"poster","props":{"imageUrl":"https://via.placeholder.com/750x400","alt":"春季招生"}},{"type":"title","props":{"text":"春季班火热招生","fontSize":24,"color":"#52c41a"}},{"type":"time","props":{"startTime":"2026-02-01","endTime":"2026-05-31"}},{"type":"button","props":{"text":"立即报名","color":"#52c41a"}}]', '{"primaryColor":"#52c41a","backgroundColor":"#f6ffed"}')
ON DUPLICATE KEY UPDATE components=VALUES(components);

-- ============================================
-- 13. 用户会话数据
-- ============================================
INSERT INTO user_sessions (id, user_id, client_ip, user_agent, login_at, last_active_at, expires_at, status) VALUES
('session_admin_001', 1, '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', NOW(), NOW(), DATE_ADD(NOW(), INTERVAL 8 HOUR), 'active'),
('session_brand_001', 2, '192.168.1.101', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', NOW(), NOW(), DATE_ADD(NOW(), INTERVAL 8 HOUR), 'active'),
('session_expired_001', 3, '192.168.1.102', 'Mozilla/5.0', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 16 HOUR), 'expired')
ON DUPLICATE KEY UPDATE status=VALUES(status);

-- ============================================
-- 14. 导出申请数据
-- ============================================
SET @brand_c := (SELECT id FROM brands WHERE name = '测试品牌C' LIMIT 1);

INSERT INTO export_requests (brand_id, requested_by, reason, filters, status, approved_by, approved_at, reject_reason, file_url, record_count) VALUES
(@brand_c, 2, '月度会员数据分析', '{"status":"active","date_range":"2026-01"}', 'approved', 1, NOW(), NULL, 'https://example.com/exports/members_202601.xlsx', 150),
(@brand_c, 2, '活跃会员导出', '{"tags":["活跃"]}', 'pending', NULL, NULL, NULL, NULL, 0),
(@brand_c, 2, '全量会员导出', '{}', 'rejected', 1, NOW(), '请说明具体用途', NULL, 0)
ON DUPLICATE KEY UPDATE status=VALUES(status);

SELECT '✓ 完整测试数据导入完成' AS result;
