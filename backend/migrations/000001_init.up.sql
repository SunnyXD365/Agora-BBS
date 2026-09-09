-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(128) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',          -- 头像 URL
    role VARCHAR(32) DEFAULT 'user',         -- 角色：admin/user
    status VARCHAR(32) DEFAULT 'active',     -- 状态：active/banned
    trust_score INT DEFAULT 100,             -- [扩展钩子] 中期忽略，后期启用信任机制
    unlock_level INT DEFAULT 1,              -- [扩展钩子] 中期忽略，后期控权
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 板块/分类表
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) UNIQUE NOT NULL,
    description TEXT DEFAULT ''
);

-- 主题帖表
CREATE TABLE topics (
    id BIGSERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id),
    author_id BIGINT NOT NULL REFERENCES users(id),
    author_name VARCHAR(64) DEFAULT '',      -- 冗余作者名，方便前端展示
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,                   -- [传统论坛] 普通富文本/Markdown 贴文
    structured_content JSONB DEFAULT '{}',   -- [扩展钩子] 后期存观点/论据等结构化数据
    status VARCHAR(32) DEFAULT 'published',  -- [扩展钩子] 中期全为 published，后期支持 cooling/reviewing
    view_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 帖子回复表
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id),
    author_name VARCHAR(64) DEFAULT '',      -- 冗余作者名
    parent_id BIGINT DEFAULT NULL REFERENCES posts(id), -- 支持楼中楼嵌套
    content TEXT NOT NULL,
    post_type VARCHAR(32) DEFAULT 'reply',   -- [扩展钩子] 后期可填 debate/evidence
    status VARCHAR(32) DEFAULT 'published',  -- [扩展钩子] 后期可填 cooling/recalled
    cooling_ends_at TIMESTAMP WITH TIME ZONE DEFAULT NULL, -- [扩展钩子] 后期配合 Temporal 定时
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 通用互动表
CREATE TABLE interactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    target_type VARCHAR(32) NOT NULL,        -- 'topic' 或 'post'
    target_id BIGINT NOT NULL,
    interaction_type VARCHAR(32) NOT NULL,   -- 中期填 'like'/'dislike'，后期可扩展为 'tag_logic' 等
    reason TEXT DEFAULT '',                  -- [扩展钩子] 后期存评价推荐理由
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_type, target_id, interaction_type)
);
