DROP TRIGGER IF EXISTS set_timestamp_comments ON comments;
DROP TRIGGER IF EXISTS set_timestamp_posts ON posts;
DROP TRIGGER IF EXISTS set_timestamp_users ON users;

DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP TABLE IF EXISTS revoked_tokens;

DROP TABLE IF EXISTS comments;

DROP TABLE IF EXISTS posts;

DROP TABLE IF EXISTS users; 