-- Remove project_id column
ALTER TABLE newsletter_jobs DROP COLUMN IF EXISTS project_id;

-- Make blog_id column NOT NULL again
ALTER TABLE newsletter_jobs ALTER COLUMN blog_id SET NOT NULL;
