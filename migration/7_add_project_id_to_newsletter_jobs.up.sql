-- Make blog_id nullable as a job can now be either for a blog or a project
ALTER TABLE newsletter_jobs ALTER COLUMN blog_id DROP NOT NULL;

-- Add project_id column referencing projects
ALTER TABLE newsletter_jobs ADD COLUMN project_id UUID REFERENCES projects(project_id) ON DELETE CASCADE;
