-- podman exec -i agora-postgres psql -U agora_user -d agora_db < ./seed.sql
-- 1. 初始化分类 (Categories)
INSERT INTO categories (id, name, slug, description, sort_order)
VALUES
  (1, '综合讨论', 'general', '默认交流与综合话题讨论区', 0),
  (2, 'Go 语言技术', 'golang', 'Go 高并发架构、标准库与开源框架讨论', 1),
  (3, '前端开发', 'frontend', 'Vue / React 组件与前端工程化讨论', 2)
ON CONFLICT (id) DO UPDATE 
SET name = EXCLUDED.name, slug = EXCLUDED.slug, description = EXCLUDED.description;

-- 重置 categories_id_seq 序列，防止后续通过 API 发帖/创分类时主键冲突
SELECT setval(pg_get_serial_sequence('categories', 'id'), COALESCE(MAX(id), 1)) FROM categories;

-- 2. 初始化官方/测试账号 (Users)
-- 密码明文均为: 123456 (BCrypt 哈希: $2a$10$e883mAn9p2rIqYl/8jM88u8J7XwU22k95fB6iHk7X1SgQ8G3K3X6W)
INSERT INTO users (id, username, email, password_hash)
VALUES
  (1, 'agora_admin', 'admin@agora.com', '$2a$10$e883mAn9p2rIqYl/8jM88u8J7XwU22k95fB6iHk7X1SgQ8G3K3X6W'),
  (2, 'gopher_test', 'gopher@example.com', '$2a$10$e883mAn9p2rIqYl/8jM88u8J7XwU22k95fB6iHk7X1SgQ8G3K3X6W')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1)) FROM users;

-- 3. 初始化官方贴 (Topics)
INSERT INTO topics (id, category_id, user_id, title, content)
VALUES
  (1, 1, 1, '欢迎使用 Agora BBS 论坛系统', 'Agora-Backend 已打通基础核心架构，欢迎在各板块交流讨论！')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('topics', 'id'), COALESCE(MAX(id), 1)) FROM topics;