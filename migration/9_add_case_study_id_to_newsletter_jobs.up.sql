ALTER TABLE newsletter_jobs ADD COLUMN case_study_id UUID REFERENCES case_studies(case_study_id) ON DELETE CASCADE;
