-- Up Migration
BEGIN;

-- 1. Remove NOT NULL constraint on username field
ALTER TABLE public."users" ALTER COLUMN username DROP NOT NULL;

-- 2. Update all user records and set username to NULL where username is an empty string
UPDATE public."users" SET username = NULL WHERE username = '';

COMMIT;
