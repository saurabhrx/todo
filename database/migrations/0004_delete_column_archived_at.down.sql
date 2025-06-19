BEGIN ;
ALTER TABLE user_session DROP COLUMN archived_at;
COMMIT;
