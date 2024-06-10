-- Down Migration
BEGIN;

-- 1. Update all user records and set username to an empty string where username is NULL
DELETE FROM public."users" WHERE username IS NULL;

-- 2. Reapply NOT NULL constraint on username field
ALTER TABLE users ALTER COLUMN username SET NOT NULL;

COMMIT;
