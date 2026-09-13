DROP TABLE IF EXISTS bookmarks;

DROP INDEX IF EXISTS idx_posts_topic_status_created;
DROP INDEX IF EXISTS idx_topics_status_created;

ALTER TABLE posts
    DROP COLUMN IF EXISTS hidden_at,
    DROP COLUMN IF EXISTS cooling_ends_at,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS post_type;

ALTER TABLE topics
    DROP COLUMN IF EXISTS hidden_at,
    DROP COLUMN IF EXISTS cooling_ends_at,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS structured_content;

ALTER TABLE categories
    DROP COLUMN IF EXISTS requires_review,
    DROP COLUMN IF EXISTS is_active;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_status_check,
    DROP CONSTRAINT IF EXISTS users_role_check,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS role;
