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

type ProjectRepository struct {
	Pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{Pool: pool}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, name string, des string, banner string, content string, slug string, featured bool) (model.Project, error) {
	query := `
		INSERT INTO projects (project_name, project_des, project_banner, project_content, project_slug, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING project_id, project_name, project_des, project_banner, project_content, project_slug, is_featured, created_at
	`

	var project model.Project
	if err := r.Pool.QueryRow(ctx, query, name, des, banner, content, slug, featured).Scan(&project.Project_ID, &project.Project_Name, &project.Project_DES, &project.Project_Banner, &project.Project_Content, &project.Project_Slug, &project.Is_Featured, &project.Created_At); err != nil {
		log.Printf("Database error: %v", err)
		return model.Project{}, err
	}

	return project, nil
}
func (r *ProjectRepository) GetAllProjects(ctx context.Context, limit int, offset int, featured *bool) ([]model.Project, int, error) {
	var totalprojects int
	countQuery := `SELECT COUNT(*) FROM projects WHERE ($1::boolean IS NULL OR is_featured = $1)`
	if err := r.Pool.QueryRow(ctx, countQuery, featured).Scan(&totalprojects); err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	rows, err := r.Pool.Query(ctx, "SELECT project_id, project_name, project_des, project_banner, project_content, project_slug, is_featured, created_at FROM projects WHERE ($1::boolean IS NULL OR is_featured = $1) ORDER BY created_at LIMIT $2 OFFSET $3", featured, limit, offset)
	if err != nil {
		return []model.Project{}, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	projects := []model.Project{}
	for rows.Next() {
		var project model.Project
		err := rows.Scan(&project.Project_ID, &project.Project_Name, &project.Project_DES, &project.Project_Banner, &project.Project_Content, &project.Project_Slug, &project.Is_Featured, &project.Created_At)
		if err != nil {
			log.Printf("Database error: %v", err)
			return []model.Project{}, 0, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row loop error: %w", err)
	}

	return projects, totalprojects, nil
}

func (r *ProjectRepository) GetProjectById(ctx context.Context, id uuid.UUID) (model.Project, error) {
	query := `
		SELECT project_id, project_name, project_des, project_banner, project_content, project_slug, is_featured, created_at FROM projects WHERE project_id = $1
	`

	var project model.Project
	if err := r.Pool.QueryRow(ctx, query, id).Scan(&project.Project_ID, &project.Project_Name, &project.Project_DES, &project.Project_Banner, &project.Project_Content, &project.Project_Slug, &project.Is_Featured, &project.Created_At); err == pgx.ErrNoRows {
		log.Println("No Project Found with that id")
		return model.Project{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepository) GetProjectBySlug(ctx context.Context, slug string) (model.Project, error) {
	query := `
		SELECT project_id, project_name, project_des, project_banner, project_content, project_slug, is_featured, created_at FROM projects WHERE project_slug = $1
	`

	var project model.Project
	if err := r.Pool.QueryRow(ctx, query, slug).Scan(&project.Project_ID, &project.Project_Name, &project.Project_DES, &project.Project_Banner, &project.Project_Content, &project.Project_Slug, &project.Is_Featured, &project.Created_At); err == pgx.ErrNoRows {
		log.Println("No Project Found with that slug")
		return model.Project{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepository) DeleteProjectById(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM projects WHERE project_id = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) UpdateProjectById(ctx context.Context, id uuid.UUID, name string, des string, content string, banner string, slug string, featured bool) (model.Project, error) {
	query := `
		UPDATE projects SET project_name = $2, project_des = $3, project_content = $4, project_banner = $5, project_slug = $6, is_featured = $7
		WHERE project_id = $1
		RETURNING project_id, project_name, project_des, project_banner, project_content, project_slug, is_featured, created_at
	`

	var project model.Project
	if err := r.Pool.QueryRow(ctx, query, id, name, des, content, banner, slug, featured).Scan(&project.Project_ID, &project.Project_Name, &project.Project_DES, &project.Project_Banner, &project.Project_Content, &project.Project_Slug, &project.Is_Featured, &project.Created_At); err != nil {
		log.Printf("Database error: %v", err)
		return model.Project{}, err
	}

	return project, nil
}
