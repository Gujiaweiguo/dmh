-- 20250210_add_performance_indexes.sql
-- 性能优化：添加复合索引

-- user_feedback 表复合索引
ALTER TABLE user_feedback ADD INDEX idx_user_created (user_id, created_at);
ALTER TABLE user_feedback ADD INDEX idx_status_created (status, created_at);
ALTER TABLE user_feedback ADD INDEX idx_category_created (category, created_at);
ALTER TABLE user_feedback ADD INDEX idx_priority_created (priority, created_at);

-- 反馈统计查询优化
ALTER TABLE user_feedback ADD INDEX idx_status_category (status, category, created_at);

-- FAQ 查询优化
ALTER TABLE faqs ADD INDEX idx_category_sort (category, sort_order, created_at);
ALTER TABLE faqs ADD INDEX idx_published_sort (is_published, sort_order, created_at);

-- 功能使用记录查询优化
ALTER TABLE feature_usage_logs ADD INDEX idx_user_feature (user_id, feature, created_at);
ALTER TABLE feature_usage_logs ADD INDEX idx_feature_action (feature, action, created_at);
ALTER TABLE feature_usage_logs ADD INDEX idx_campaign_feature (campaign_id, feature, created_at);

-- 满意度调查查询优化
ALTER TABLE feature_satisfaction_surveys ADD INDEX idx_user_feature (user_id, feature, created_at);
ALTER TABLE feature_satisfaction_surveys ADD INDEX idx_feature_rating (feature, overall_satisfaction, created_at);
