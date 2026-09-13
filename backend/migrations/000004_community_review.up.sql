ALTER TABLE topics ADD COLUMN feedback_score INT NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN feedback_score INT NOT NULL DEFAULT 0;

CREATE TABLE contextual_feedbacks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(16) NOT NULL,
    target_id BIGINT NOT NULL,
    stance VARCHAR(16) NOT NULL,
    tag VARCHAR(32) NOT NULL,
    reason VARCHAR(120) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'published',
    llm_audit_status VARCHAR(24) NOT NULL DEFAULT 'not_selected',
    llm_audit_result JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT contextual_feedback_target_check CHECK (target_type IN ('topic', 'post')),
    CONSTRAINT contextual_feedback_stance_check CHECK (stance IN ('support', 'challenge')),
    CONSTRAINT contextual_feedback_status_check CHECK (status IN ('published', 'rejected', 'withdrawn')),
    CONSTRAINT contextual_feedback_reason_length CHECK (char_length(reason) BETWEEN 5 AND 120),
    CONSTRAINT uk_contextual_feedback_user_target UNIQUE(user_id, target_type, target_id)
);
CREATE INDEX idx_feedback_target ON contextual_feedbacks(target_type, target_id, status);

INSERT INTO contextual_feedbacks(user_id, target_type, target_id, stance, tag, reason, status, llm_audit_status, created_at, updated_at)
SELECT user_id, target_type, target_id, 'support', 'legacy_support', '历史点赞迁移记录', 'published', 'excluded', created_at, created_at
FROM likes ON CONFLICT(user_id, target_type, target_id) DO NOTHING;

CREATE TABLE blind_review_batches (
    id BIGSERIAL PRIMARY KEY,
    subject_type VARCHAR(16) NOT NULL,
    subject_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(24) NOT NULL DEFAULT 'pending',
    final_result VARCHAR(16),
    deadline TIMESTAMPTZ NOT NULL,
    llm_result JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    CONSTRAINT blind_review_subject_check CHECK (subject_type IN ('user', 'topic')),
    CONSTRAINT blind_review_batch_status_check CHECK (status IN ('pending', 'completed', 'expired', 'failed')),
    CONSTRAINT uk_blind_review_subject UNIQUE(subject_type, subject_id)
);

CREATE TABLE blind_review_tasks (
    id BIGSERIAL PRIMARY KEY,
    batch_id BIGINT NOT NULL REFERENCES blind_review_batches(id) ON DELETE CASCADE,
    reviewer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_status VARCHAR(24) NOT NULL DEFAULT 'pending',
    appropriateness BOOLEAN,
    sincerity BOOLEAN,
    review_result VARCHAR(16),
    reason VARCHAR(500) NOT NULL DEFAULT '',
    llm_check_status VARCHAR(24) NOT NULL DEFAULT 'pending',
    llm_check_result JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    CONSTRAINT blind_review_task_status_check CHECK (task_status IN ('pending', 'completed', 'expired')),
    CONSTRAINT uk_blind_review_batch_reviewer UNIQUE(batch_id, reviewer_id)
);
CREATE INDEX idx_review_tasks_reviewer_status ON blind_review_tasks(reviewer_id, task_status, created_at);

CREATE TABLE llm_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_key VARCHAR(160) UNIQUE NOT NULL,
    job_type VARCHAR(32) NOT NULL,
    aggregate_type VARCHAR(32) NOT NULL,
    aggregate_id BIGINT NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    provider VARCHAR(64) NOT NULL DEFAULT 'openai-compatible',
    model VARCHAR(128) NOT NULL DEFAULT '',
    result JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message TEXT NOT NULL DEFAULT '',
    prompt_tokens INT NOT NULL DEFAULT 0,
    completion_tokens INT NOT NULL DEFAULT 0,
    latency_ms INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ
);
CREATE INDEX idx_llm_jobs_status_created ON llm_jobs(status, created_at);

CREATE TABLE comment_clusters (
    id BIGSERIAL PRIMARY KEY,
    topic_id BIGINT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    tag VARCHAR(64) NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    weight DOUBLE PRECISION NOT NULL DEFAULT 0,
    model VARCHAR(128) NOT NULL DEFAULT '',
    generation INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_comment_clusters_topic_weight ON comment_clusters(topic_id, weight DESC);

CREATE TABLE post_cluster_assignments (
    cluster_id BIGINT NOT NULL REFERENCES comment_clusters(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    relevance DOUBLE PRECISION NOT NULL DEFAULT 0,
    PRIMARY KEY(cluster_id, post_id)
);
