package reposiotry

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *CaseStudyRepository) GetAllCaseStudies(ctx context.Context, limit int, offset int) ([]model.Case_Study, int, error) {
	var totalCaseStudies int
	countQuery := `SELECT COUNT(*) FROM case_studies`
	if err := r.Pool.QueryRow(ctx, countQuery).Scan(&totalCaseStudies); err != nil {
		return nil, 0, fmt.Errorf("failed to count case studies: %w", err)
	}

	rows, err := r.Pool.Query(ctx, "SELECT * FROM case_studies ORDER BY created_at LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return []model.Case_Study{}, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var case_studies []model.Case_Study
	for rows.Next() {
		var case_study model.Case_Study
		if err := rows.Scan(&case_study.Case_Study_ID, &case_study.Case_Study_Title, &case_study.Case_Study_DES, &case_study.Case_Study_Slug, &case_study.Case_Study_Content, &case_study.Case_Study_Banner, &case_study.Created_At); err == sql.ErrNoRows {
			log.Printf("Database error: %v", err)
			return []model.Case_Study{}, 0, err
		}
		case_studies = append(case_studies, case_study)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row loop error: %w", err)
	}

	if len(case_studies) == 0 {
		return []model.Case_Study{}, 0, model.ErrCaseStudyNotFound
	}

	return case_studies, totalCaseStudies, nil
}

func (r *CaseStudyRepository) GetCaseStudyBySlug(ctx context.Context, slug string) (model.Case_Study, error) {
	query := `
		SELECT * FROM case_studies WHERE case_study_slug = $1
	`

	var case_study model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, slug).Scan(&case_study.Case_Study_ID, &case_study.Case_Study_Title, &case_study.Case_Study_DES, &case_study.Case_Study_Slug, &case_study.Case_Study_Content, &case_study.Case_Study_Banner, &case_study.Created_At); err == sql.ErrNoRows {
		log.Println("No Case Study found with that slug")
		return model.Case_Study{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Case_Study{}, err
	}

	return case_study, nil
}
