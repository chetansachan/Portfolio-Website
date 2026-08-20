package service

import (
	"context"
	"fmt"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"

	"github.com/google/uuid"
)

type ProjectService struct {
	ProjectRepository *reposiotry.ProjectRepository
}

func NewProjectService(projectrepository *reposiotry.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectrepository}
}

func (s *ProjectService) CreateProject_Service(ctx context.Context, name string, des string, banner string, content string, slug string, featured bool) (model.Project, error) {
	if name == "" || des == "" || banner == "" || content == "" || slug == "" {
		log.Println("[WARN] validation failed")
		return model.Project{}, model.ErrValidation
	}

	project, err := s.ProjectRepository.CreateProject(ctx, name, des, banner, content, slug, featured)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Project{}, fmt.Errorf("failed to save project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) GetProjectById_Service(ctx context.Context, id uuid.UUID) (model.Project, error) {
	if id == uuid.Nil {
		log.Println("[WARN] ID is Required")
		return model.Project{}, model.ErrIdValidation
	}

	project, err := s.ProjectRepository.GetProjectById(ctx, id)
	if err != nil {
		log.Println("[Error] Error Geting project from database")
		return model.Project{}, model.ErrProjectNotFound
	}

	return project, nil
}
