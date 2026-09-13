-- A newcomer reaches the L3 trust threshold after an approved self-introduction.
ALTER TABLE user_trust_profiles ALTER COLUMN trust_score SET DEFAULT 15;

-- Only adjust untouched newcomer profiles; historical trust changes remain intact.
UPDATE user_trust_profiles tp
SET trust_score = 15, updated_at = CURRENT_TIMESTAMP
FROM user_profiles p
WHERE p.user_id = tp.user_id
  AND tp.trust_score = 0
  AND tp.compliant_interactions = 0
  AND p.onboarding_status = 'not_submitted'
  AND NOT EXISTS (SELECT 1 FROM trust_logs l WHERE l.user_id = tp.user_id);
