-- 清理品牌管理员相关数据的SQL脚本

USE dmh;

-- 1. 清理用户表中的brand_admin角色
UPDATE users SET role = 'participant' WHERE role = 'brand_admin';

-- 2. 清理角色表中的brand_admin角色记录
DELETE FROM roles WHERE code = 'brand_admin';

-- 3. 清理用户角色关联表中的brand_admin关联
DELETE FROM user_roles WHERE role_id = (SELECT id FROM roles WHERE code = 'brand_admin');

-- 4. 清理角色权限关联表中的brand_admin权限
DELETE FROM role_permissions WHERE role_id = (SELECT id FROM roles WHERE code = 'brand_admin');

-- 5. 清理角色菜单关联表中的brand_admin菜单
DELETE FROM role_menus WHERE role_id = (SELECT id FROM roles WHERE code = 'brand_admin');

-- 6. 清理品牌关系管理相关菜单
DELETE FROM menus WHERE code = 'brand-relation';

-- 7. 清理用户品牌关联表（如果不再需要）
-- 注意：如果还需要保留品牌关联功能，请注释掉下面这行
-- DROP TABLE IF EXISTS user_brands;

-- 8. 更新测试用户brand_manager的角色为participant
UPDATE users SET role = 'participant' WHERE username = 'brand_manager';

-- 9. 重新分配brand_manager用户的角色关联
DELETE FROM user_roles WHERE user_id = (SELECT id FROM users WHERE username = 'brand_manager');
INSERT INTO user_roles (user_id, role_id) 
SELECT u.id, r.id 
FROM users u, roles r 
WHERE u.username = 'brand_manager' AND r.code = 'participant';

COMMIT;