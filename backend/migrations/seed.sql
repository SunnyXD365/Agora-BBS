-- ==========================================
-- Agora-BBS 本地开发环境 Seed 数据脚本
-- 说明：包含基础版块、测试用户、示例帖子及评论
-- ==========================================

BEGIN;

-- 插入基础版块
INSERT INTO categories (id, name, slug, description, created_at, updated_at)
VALUES 
    (1, '综合讨论', 'general', '自由交流各类技术与生活话题', NOW(), NOW()),
    (2, '技术干货', 'tech', 'Go、Next.js、云原生与架构实践', NOW(), NOW()),
    (3, '反馈与建议', 'feedback', '社区功能建议与 Bug 反馈', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 插入测试用户
-- 注意：实际项目中密码使用 bcrypt 散列（以下哈希值对应的明文统一为：password123）
INSERT INTO users (id, username, email, password_hash, avatar, role, status, created_at, updated_at)
VALUES 
    (1, 'admin', 'admin@agora.com', '$2a$10$7EqJtq986P22m4kL1f8eu.E2.Q2xJgT6Dk4Z0P2q4e/K3Z1mB4C.i', 'https://api.dicebear.com/7.x/bottts/svg?seed=admin', 'admin', 'active', NOW(), NOW()),
    (2, 'alice', 'alice@agora.com', '$2a$10$7EqJtq986P22m4kL1f8eu.E2.Q2xJgT6Dk4Z0P2q4e/K3Z1mB4C.i', 'https://api.dicebear.com/7.x/bottts/svg?seed=alice', 'user', 'active', NOW(), NOW()),
    (3, 'bob', 'bob@agora.com', '$2a$10$7EqJtq986P22m4kL1f8eu.E2.Q2xJgT6Dk4Z0P2q4e/K3Z1mB4C.i', 'https://api.dicebear.com/7.x/bottts/svg?seed=bob', 'user', 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 插入示例帖子
INSERT INTO posts (id, user_id, category_id, title, content, status, view_count, created_at, updated_at)
VALUES 
    (1, 1, 1, '欢迎来到 Agora-BBS 社区！', '这是一个基于 Go + Next.js + Temporal 架构构建的现代论坛系统。欢迎在此畅所欲言！', 'published', 102, NOW(), NOW()),
    (2, 2, 2, 'Go 1.23 特性解析与最佳实践', 'Go 1.23 带来了不少迭代器相关的增强，本文来聊聊如何在实际项目中使用 iterator...', 'published', 45, NOW(), NOW()),
    (3, 3, 3, '建议增加暗黑模式支持', '希望前端能支持 Dark Mode 切换，夜间看社区有点烫眼。', 'published', 12, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 插入示例评论
INSERT INTO comments (id, post_id, user_id, content, created_at, updated_at)
VALUES 
    (1, 1, 2, '支持！期待后续更多 AI 功能落地。', NOW(), NOW()),
    (2, 1, 3, '环境一键拉起体验很好，赞！', NOW(), NOW()),
    (3, 3, 1, '收到建议，暗黑模式已列入后续路线图。', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 重置 Serial 自增序列（防止手动插入固定 ID 后，后续插入新记录出现 ID 冲突）
SELECT setval(pg_get_serial_sequence('categories', 'id'), COALESCE(MAX(id), 1)) FROM categories;
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1)) FROM users;
SELECT setval(pg_get_serial_sequence('posts', 'id'), COALESCE(MAX(id), 1)) FROM posts;
SELECT setval(pg_get_serial_sequence('comments', 'id'), COALESCE(MAX(id), 1)) FROM comments;

COMMIT;
