CREATE TABLE IF NOT EXISTS admin_login_challenges (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(64) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_login_challenges_user_created
    ON admin_login_challenges(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_login_challenges_expires
    ON admin_login_challenges(expires_at);
