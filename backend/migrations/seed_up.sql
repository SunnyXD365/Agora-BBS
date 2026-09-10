-- 密码明文均为: password123 (Bcrypt Hash)
INSERT INTO users (id, username, password_hash, email, avatar, role, status, trust_score, unlock_level)
VALUES 
(1, 'admin', '$2a$10$e83pP40W1H/1NfUv0YyLq.aJvS.tG1zY.z3rL4M4e444444444444', 'admin@agora.com', 'https://api.dicebear.com/7.x/bottts/svg?seed=admin', 'admin', 'active', 999, 10),
(2, 'alice', '$2a$10$e83pP40W1H/1NfUv0YyLq.aJvS.tG1zY.z3rL4M4e444444444444', 'alice@agora.com', 'https://api.dicebear.com/7.x/bottts/svg?seed=alice', 'user', 'active', 100, 1),
(3, 'bob',   '$2a$10$e83pP40W1H/1NfUv0YyLq.aJvS.tG1zY.z3rL4M4e444444444444', 'bob@agora.com',   'https://api.dicebear.com/7.x/bottts/svg?seed=bob',   'user', 'active', 100, 1)
ON CONFLICT (id) DO NOTHING;

-- 初始板块分类
INSERT INTO categories (id, name, slug, description, parent_id, sort_order)
VALUES 
(1, '官方公告', 'announcement', '社区重要通知与规则更新', NULL, 1),
(2, '技术交流', 'tech', 'Go、React、系统架构与算法讨论', NULL, 2),
(3, '综合闲聊', 'chat', '日常生活与灌水分享', NULL, 3)
ON CONFLICT (id) DO NOTHING;

-- 初始测试主题帖
INSERT INTO topics (id, category_id, user_id, title, content, view_count, post_count, like_count, is_sticky, is_essence, status)
VALUES 
(1, 1, 1, '欢迎来到 Agora BBS 论坛！', '这里是 Agora 论坛的官方欢迎帖，请遵守社区规范。', 105, 1, 12, true, true, 'normal'),
(2, 2, 2, 'Golang 1.24 并发优化实践', '分享一下最近在项目中使用的 Gin + Postgres 高并发架构心得。', 42, 0, 5, false, false, 'normal')
ON CONFLICT (id) DO NOTHING;

-- 初始测试评论
INSERT INTO posts (id, topic_id, user_id, parent_id, content, like_count, status)
VALUES 
(1, 1, 3, NULL, '支持！前排打卡！', 3, 'normal')
ON CONFLICT (id) DO NOTHING;

-- 重置 Serial 序列锁
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));
SELECT setval('topics_id_seq', (SELECT MAX(id) FROM topics));
SELECT setval('posts_id_seq', (SELECT MAX(id) FROM posts));
