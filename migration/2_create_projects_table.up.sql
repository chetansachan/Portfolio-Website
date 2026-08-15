CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    project_name VARCHAR(150) NOT NULL,
    project_des VARCHAR(350) NOT NULL,
    project_url VARCHAR(2048) NOT NULL,
    project_banner VARCHAR(2048) NOT NULL,
    is_featured BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);