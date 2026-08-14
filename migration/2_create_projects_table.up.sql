CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    project_name VARCHAR(50) NOT NULL,
    project_des VARCHAR(350) NOT NULL,
    project_url VARCHAR(350) NOT NULL,
    project_vid_url VARCHAR(350) NOT NULL,
    is_featured BOOLEAN,
    created_at VARCHAR(100) NOT NULL
);