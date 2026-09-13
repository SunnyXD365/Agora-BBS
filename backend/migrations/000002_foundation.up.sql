-- Additive compatibility migration for databases that already applied 000001.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(16) NOT NULL DEFAULT 'user',
    ADD COLUMN IF NOT EXISTS status VARCHAR(16) NOT NULL DEFAULT 'active';

ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS requires_review BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE topics
    ADD COLUMN IF NOT EXISTS structured_content JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS status VARCHAR(24) NOT NULL DEFAULT 'published',
    ADD COLUMN IF NOT EXISTS cooling_ends_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS hidden_at TIMESTAMPTZ;

UPDATE topics
SET structured_content = jsonb_build_object(
    'claim', content,
    'evidence', '',
    'uncertainty', ''
)
WHERE structured_content = '{}'::jsonb;

ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS post_type VARCHAR(24) NOT NULL DEFAULT 'experience',
    ADD COLUMN IF NOT EXISTS status VARCHAR(24) NOT NULL DEFAULT 'published',
    ADD COLUMN IF NOT EXISTS cooling_ends_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS hidden_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS bookmarks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_bookmarks_user_topic UNIQUE (user_id, topic_id)
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_created
    ON bookmarks(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_topics_status_created
    ON topics(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_topic_status_created
    ON posts(topic_id, status, created_at ASC);

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('user', 'admin'));
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users ADD CONSTRAINT users_status_check CHECK (status IN ('active', 'suspended'));
