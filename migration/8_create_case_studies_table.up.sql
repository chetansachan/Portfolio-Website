CREATE TABLE IF NOT EXISTS case_studies (
    case_study_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    case_study_title VARCHAR(150) NOT NULL,
    case_study_des VARCHAR(350) NOT NULL,
    case_study_slug VARCHAR(150) NOT NULL UNIQUE,
    case_study_content TEXT NOT NULL,
    case_study_banner VARCHAR(2048) NOT NULL,
    is_featured BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
