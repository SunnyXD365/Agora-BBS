CREATE TABLE user_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    onboarding_statement TEXT NOT NULL DEFAULT '',
    background_tag VARCHAR(64) NOT NULL DEFAULT '',
    onboarding_status VARCHAR(24) NOT NULL DEFAULT 'not_submitted',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_profiles_onboarding_status_check
        CHECK (onboarding_status IN ('not_submitted', 'pending_review', 'approved', 'rejected'))
);

CREATE TABLE user_trust_profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    trust_score INT NOT NULL DEFAULT 0,
    unlock_level INT NOT NULL DEFAULT 0,
    audit_probability DOUBLE PRECISION NOT NULL DEFAULT 0.10,
    verified_read_seconds BIGINT NOT NULL DEFAULT 0,
    compliant_interactions INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_trust_level_check CHECK (unlock_level BETWEEN 0 AND 3),
    CONSTRAINT user_audit_probability_check CHECK (audit_probability BETWEEN 0.05 AND 0.80)
);

-- Accounts that predate governance keep their existing posting abilities.
INSERT INTO user_profiles(user_id)
SELECT id FROM users ON CONFLICT (user_id) DO NOTHING;
INSERT INTO user_trust_profiles(user_id, trust_score, unlock_level, audit_probability)
SELECT id, 20, 3, 0.05 FROM users ON CONFLICT (user_id) DO NOTHING;

CREATE OR REPLACE FUNCTION create_user_governance_profiles() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_profiles(user_id) VALUES (NEW.id) ON CONFLICT DO NOTHING;
    INSERT INTO user_trust_profiles(user_id) VALUES (NEW.id) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_create_governance_profiles
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION create_user_governance_profiles();

CREATE TABLE reading_sessions (
    id BIGSERIAL PRIMARY KEY,
    public_id VARCHAR(64) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    last_progress SMALLINT NOT NULL DEFAULT 0,
    reading_seconds INT NOT NULL DEFAULT 0,
    reply_dwell_seconds INT NOT NULL DEFAULT 0,
    bottom_reached BOOLEAN NOT NULL DEFAULT FALSE,
    eligible BOOLEAN NOT NULL DEFAULT FALSE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    last_heartbeat_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    CONSTRAINT reading_progress_check CHECK (last_progress BETWEEN 0 AND 100)
);
CREATE INDEX idx_reading_sessions_user_topic ON reading_sessions(user_id, topic_id, created_at DESC);

CREATE TABLE trust_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(64) NOT NULL,
    score_delta INT NOT NULL,
    reason TEXT NOT NULL,
    reference_type VARCHAR(32),
    reference_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_trust_logs_user_created ON trust_logs(user_id, created_at DESC);

CREATE TABLE governance_events (
    id BIGSERIAL PRIMARY KEY,
    event_key VARCHAR(128) UNIQUE NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    aggregate_type VARCHAR(32) NOT NULL,
    aggregate_id BIGINT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(24) NOT NULL DEFAULT 'pending',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ
);
CREATE INDEX idx_governance_events_status_created ON governance_events(status, created_at);
