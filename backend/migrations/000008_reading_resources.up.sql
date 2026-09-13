-- Generalize reading sessions so important system documents can contribute to
-- verified reading without pretending to be forum topics.
ALTER TABLE reading_sessions
    ALTER COLUMN topic_id DROP NOT NULL,
    ADD COLUMN resource_type VARCHAR(24) NOT NULL DEFAULT 'topic',
    ADD COLUMN resource_key VARCHAR(64) NOT NULL DEFAULT '';

UPDATE reading_sessions
SET resource_type = 'topic', resource_key = topic_id::text;

ALTER TABLE reading_sessions
    ADD CONSTRAINT reading_sessions_resource_check CHECK (
        (resource_type = 'topic' AND topic_id IS NOT NULL AND resource_key = topic_id::text)
        OR
        (resource_type = 'guide' AND topic_id IS NULL AND resource_key = 'forum-guide')
    );

CREATE INDEX idx_reading_sessions_user_resource
    ON reading_sessions(user_id, resource_type, resource_key, completed, created_at DESC);
