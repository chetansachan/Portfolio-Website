CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS projects (
    project_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_name VARCHAR(150) NOT NULL,
    project_des VARCHAR(350) NOT NULL,
    project_banner TEXT NOT NULL,
    project_content TEXT NOT NULL,
    project_slug TEXT NOT NULL,
    is_featured BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);