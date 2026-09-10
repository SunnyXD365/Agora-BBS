-- 1. 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(128) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',          -- 头像 URL
    role VARCHAR(32) DEFAULT 'user',          -- 角色：admin/user
    status VARCHAR(32) DEFAULT 'active',      -- 状态：active/banned
    trust_score INT DEFAULT 100,              -- 隐藏信任积分钩子
    unlock_level INT DEFAULT 1,               -- 渐进解锁等级钩子
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. 板块/分类表
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. 主题帖表
CREATE TABLE topics (
    id BIGSERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id),
    author_id BIGINT NOT NULL REFERENCES users(id),
    author_name VARCHAR(64) DEFAULT '',      -- 冗余作者名，方便前端直接展示
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,                   -- 普通富文本/Markdown 贴文
    structured_content JSONB DEFAULT '{}',   -- 观点/论据结构化数据钩子
    status VARCHAR(32) DEFAULT 'published',  -- published / cooling / reviewing / recalled
    cooling_ends_at TIMESTAMP WITH TIME ZONE DEFAULT NULL, -- 主楼冷静期倒计时钩子
    view_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. 帖子回复表
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id),
    author_name VARCHAR(64) DEFAULT '',      -- 冗余作者名
    parent_id BIGINT DEFAULT NULL REFERENCES posts(id) ON DELETE SET NULL, -- 楼中楼嵌套，父楼层删除时保留子楼层
    content TEXT NOT NULL,
    post_type VARCHAR(32) DEFAULT 'reply',   -- debate/evidence/experience/thanks
    status VARCHAR(32) DEFAULT 'published',  -- published / cooling / recalled
    cooling_ends_at TIMESTAMP WITH TIME ZONE DEFAULT NULL, -- Temporal 异步冷静期倒计时
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. 通用互动表（赞踩 / 语境标签 / 盲审评估）
CREATE TABLE interactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    target_type VARCHAR(32) NOT NULL,        -- 'topic' 或 'post'
    target_id BIGINT NOT NULL,
    interaction_type VARCHAR(32) NOT NULL,   -- 'like' / 'dislike' / 'tag_logic' / 'tag_perspective'
    reason TEXT DEFAULT '',                  -- 推荐理由 / 附带评价
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_type, target_id, interaction_type)
);

-- =================================================================
-- 性能优化索引（保障 1000+ QPS 压测与高频查询）
-- =================================================================

-- 话题表高频查询索引
CREATE INDEX idx_topics_category_id ON topics(category_id);
CREATE INDEX idx_topics_author_id ON topics(author_id);
CREATE INDEX idx_topics_created_at ON topics(created_at DESC);

-- 回复表高频查询索引
CREATE INDEX idx_posts_topic_id ON posts(topic_id);
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_parent_id ON posts(parent_id);

-- 互动表复合索引
CREATE INDEX idx_interactions_target ON interactions(target_type, target_id);
CREATE INDEX idx_interactions_user_id ON interactions(user_id);
