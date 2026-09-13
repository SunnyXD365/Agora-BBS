DROP TABLE IF EXISTS post_cluster_assignments;
DROP TABLE IF EXISTS comment_clusters;
DROP TABLE IF EXISTS llm_jobs;
DROP TABLE IF EXISTS blind_review_tasks;
DROP TABLE IF EXISTS blind_review_batches;
DROP TABLE IF EXISTS contextual_feedbacks;
ALTER TABLE posts DROP COLUMN IF EXISTS feedback_score;
ALTER TABLE topics DROP COLUMN IF EXISTS feedback_score;
