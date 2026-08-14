CREATE TABLE IF NOT EXISTS blogs (
    blog_id SERIAL PRIMARY KEY,
    blog_name VARCHAR(50) NOT NULL,
    blog_des VARCHAR(350) NOT NULL,
    blog_content TEXT NOT NULL,
    created_at VARCHAR(100) NOT NULL
);