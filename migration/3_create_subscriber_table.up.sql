CREATE TABLE IF NOT EXiST subscriber (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    subscribed_at TIMESTAMP DEFAULT NOW(),
    unsubscribed_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);