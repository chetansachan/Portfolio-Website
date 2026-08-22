CREATE TABEL IF NOT EXIST newsletter_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blog_id UUID NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    subscribers_group INT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
)