DROP INDEX IF EXISTS idx_reading_sessions_user_resource;
ALTER TABLE reading_sessions DROP CONSTRAINT IF EXISTS reading_sessions_resource_check;
DELETE FROM reading_sessions WHERE topic_id IS NULL;
ALTER TABLE reading_sessions
    DROP COLUMN IF EXISTS resource_key,
    DROP COLUMN IF EXISTS resource_type,
    ALTER COLUMN topic_id SET NOT NULL;
