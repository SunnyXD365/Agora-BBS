DROP TABLE IF EXISTS governance_events;
DROP TABLE IF EXISTS trust_logs;
DROP TABLE IF EXISTS reading_sessions;
DROP TRIGGER IF EXISTS users_create_governance_profiles ON users;
DROP FUNCTION IF EXISTS create_user_governance_profiles;
DROP TABLE IF EXISTS user_trust_profiles;
DROP TABLE IF EXISTS user_profiles;
