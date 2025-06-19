BEGIN ;

CREATE UNIQUE INDEX IF NOT EXISTS active_user on users(email) WHERE archived_at IS NULL ;


COMMIT ;