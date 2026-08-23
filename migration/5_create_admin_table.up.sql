CREATE TABLE IF NOT EXISTS admin (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password VARCHAR(80) NOT NULL,
    login_at TIMESTAMP DEFAULT NOW()
);