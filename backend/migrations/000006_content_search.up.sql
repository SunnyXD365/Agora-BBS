CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_topics_published_search_trgm
    ON topics USING GIN ((title || ' ' || content || ' ' || structured_content::text) gin_trgm_ops)
    WHERE status = 'published';

CREATE INDEX IF NOT EXISTS idx_posts_published_content_trgm
    ON posts USING GIN (content gin_trgm_ops)
    WHERE status = 'published';
