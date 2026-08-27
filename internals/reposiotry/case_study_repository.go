package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CaseStudyRepository struct {
	Pool *pgxpool.Pool
}

func NewCaseStudyRepository(pool *pgxpool.Pool) *CaseStudyRepository {
	return &CaseStudyRepository{Pool: pool}
}

func (r *CaseStudyRepository) CreateCaseStudy(ctx context.Context, Title string, DES string, Slug string, Content string, Banner string) (model.Case_Study, error) {
	query := `
	INSERT INTO case_studies (title, des, slug, content, banner) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, des, slug, content, banner, created_at
	`

	var case_studies model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, Title, DES, Slug, Content, Banner).Scan(&case_studies.Case_Study_ID, &case_studies.Case_Study_Title, &case_studies.Case_Study_DES, &case_studies.Case_Study_Slug, &case_studies.Case_Study_Content, &case_studies.Case_Study_Banner, &case_studies.Created_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return case_studies, nil
}
