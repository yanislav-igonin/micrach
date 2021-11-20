ALTER TABLE posts
ALTER COLUMN updated_at DROP NOT NULL;

ALTER TABLE posts
ALTER COLUMN updated_at DROP DEFAULT;

UPDATE posts
SET updated_at = null
WHERE is_parent != true;
