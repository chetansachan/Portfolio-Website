package reposiotry

import (
	"context"
	"github.com/jackc/pgx/v5"
	"fmt"
	"log"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CaseStudyRepository struct {
	Pool *pgxpool.Pool
}

func NewCaseStudyRepository(pool *pgxpool.Pool) *CaseStudyRepository {
	return &CaseStudyRepository{Pool: pool}
}

func (r *CaseStudyRepository) CreateCaseStudy(ctx context.Context, Title string, DES string, Slug string, Content string, Banner string, featured bool) (model.Case_Study, error) {
	query := `
	INSERT INTO case_studies (case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at
	`

	var case_study model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, Title, DES, Slug, Content, Banner, featured).Scan(
		&case_study.Case_Study_ID,
		&case_study.Case_Study_Title,
		&case_study.Case_Study_DES,
		&case_study.Case_Study_Slug,
		&case_study.Case_Study_Content,
		&case_study.Case_Study_Banner,
		&case_study.Is_Featured,
		&case_study.Created_At,
	); err != nil {
		log.Printf("Database error: %v", err)
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (r *CaseStudyRepository) GetAllCaseStudies(ctx context.Context, limit int, offset int, featured *bool) ([]model.Case_Study, int, error) {
	var totalCaseStudies int
	countQuery := `SELECT COUNT(*) FROM case_studies WHERE ($1::boolean IS NULL OR is_featured = $1)`
	if err := r.Pool.QueryRow(ctx, countQuery, featured).Scan(&totalCaseStudies); err != nil {
		return nil, 0, fmt.Errorf("failed to count case studies: %w", err)
	}

	rows, err := r.Pool.Query(ctx, "SELECT case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at FROM case_studies WHERE ($1::boolean IS NULL OR is_featured = $1) ORDER BY created_at LIMIT $2 OFFSET $3", featured, limit, offset)
	if err != nil {
		return []model.Case_Study{}, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	case_studies := []model.Case_Study{}
	for rows.Next() {
		var case_study model.Case_Study
		if err := rows.Scan(
			&case_study.Case_Study_ID,
			&case_study.Case_Study_Title,
			&case_study.Case_Study_DES,
			&case_study.Case_Study_Slug,
			&case_study.Case_Study_Content,
			&case_study.Case_Study_Banner,
			&case_study.Is_Featured,
			&case_study.Created_At,
		); err == pgx.ErrNoRows {
			log.Printf("Database error: %v", err)
			return []model.Case_Study{}, 0, err
		} else if err != nil {
			log.Printf("Database error: %v", err)
			return []model.Case_Study{}, 0, err
		}
		case_studies = append(case_studies, case_study)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row loop error: %w", err)
	}

	return case_studies, totalCaseStudies, nil
}

func (r *CaseStudyRepository) GetCaseStudyBySlug(ctx context.Context, slug string) (model.Case_Study, error) {
	query := `
		SELECT case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at FROM case_studies WHERE case_study_slug = $1
	`

	var case_study model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, slug).Scan(
		&case_study.Case_Study_ID,
		&case_study.Case_Study_Title,
		&case_study.Case_Study_DES,
		&case_study.Case_Study_Slug,
		&case_study.Case_Study_Content,
		&case_study.Case_Study_Banner,
		&case_study.Is_Featured,
		&case_study.Created_At,
	); err == pgx.ErrNoRows {
		log.Println("No Case Study found with that slug")
		return model.Case_Study{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (r *CaseStudyRepository) GetCaseStudyById(ctx context.Context, id uuid.UUID) (model.Case_Study, error) {
	query := `
		SELECT case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at FROM case_studies WHERE case_study_id = $1
	`

	var case_study model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, id).Scan(
		&case_study.Case_Study_ID,
		&case_study.Case_Study_Title,
		&case_study.Case_Study_DES,
		&case_study.Case_Study_Slug,
		&case_study.Case_Study_Content,
		&case_study.Case_Study_Banner,
		&case_study.Is_Featured,
		&case_study.Created_At,
	); err == pgx.ErrNoRows {
		log.Println("No Case Study found with that ID")
		return model.Case_Study{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (r *CaseStudyRepository) DeleteCaseStudyById(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM case_studies WHERE case_study_id = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrCaseStudyNotFound
	}

	return nil
}

func (r *CaseStudyRepository) UpdateCaseStudyById(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string, featured bool) (model.Case_Study, error) {
	query := `
		UPDATE case_studies 
		SET case_study_title = $2, case_study_des = $3, case_study_slug = $4, case_study_content = $5, case_study_banner = $6, is_featured = $7
		WHERE case_study_id = $1
		RETURNING case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at
	`

	var case_study model.Case_Study
	if err := r.Pool.QueryRow(ctx, query, id, title, des, slug, content, banner, featured).Scan(
		&case_study.Case_Study_ID,
		&case_study.Case_Study_Title,
		&case_study.Case_Study_DES,
		&case_study.Case_Study_Slug,
		&case_study.Case_Study_Content,
		&case_study.Case_Study_Banner,
		&case_study.Is_Featured,
		&case_study.Created_At,
	); err != nil {
		log.Printf("Database error: %v", err)
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (r *CaseStudyRepository) GetFeaturedCaseStudies(ctx context.Context) ([]model.Case_Study, error) {
	query := `
		SELECT case_study_id, case_study_title, case_study_des, case_study_slug, case_study_content, case_study_banner, is_featured, created_at 
		FROM case_studies 
		WHERE is_featured = true 
		ORDER BY created_at DESC
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return []model.Case_Study{}, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var case_studies []model.Case_Study
	for rows.Next() {
		var case_study model.Case_Study
		if err := rows.Scan(
			&case_study.Case_Study_ID,
			&case_study.Case_Study_Title,
			&case_study.Case_Study_DES,
			&case_study.Case_Study_Slug,
			&case_study.Case_Study_Content,
			&case_study.Case_Study_Banner,
			&case_study.Is_Featured,
			&case_study.Created_At,
		); err != nil {
			log.Printf("Database error: %v", err)
			return []model.Case_Study{}, err
		}
		case_studies = append(case_studies, case_study)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row loop error: %w", err)
	}

	return case_studies, nil
}
