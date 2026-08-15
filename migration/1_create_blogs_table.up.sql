CREATE TABLE IF NOT EXISTS blogs (
    blog_id SERIAL PRIMARY KEY,
    blog_title VARCHAR(150) NOT NULL,
    blog_des VARCHAR(350) NOT NULL,
    blog_slug VARCHAR(150) NOT NULL UNIQUE,
    blog_content TEXT NOT NULL,
    blog_banner VARCHAR(2048) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
