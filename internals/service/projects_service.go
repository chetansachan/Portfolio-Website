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
	NewsletterRepo    *reposiotry.NewsletterRepo
}

func NewProjectService(projectrepository *reposiotry.ProjectRepository, newsletterRepo *reposiotry.NewsletterRepo) *ProjectService {
	return &ProjectService{
		ProjectRepository: projectrepository,
		NewsletterRepo:    newsletterRepo,
	}
}

func (s *ProjectService) CreateProject_Service(ctx context.Context, name string, des string, banner string, content string, slug string, featured bool) (model.Project, error) {
	if name == "" || des == "" || banner == "" || content == "" || slug == "" {
		log.Println("[WARN] validation failed")
		return model.Project{}, model.ErrValidation
	}

	p, err := s.ProjectRepository.GetProjectBySlug(ctx, slug)
	if err == nil {
		log.Println("Project with this slug is already exist")
		return p, model.ErrProjectAlreadyExist
	}

	project, err := s.ProjectRepository.CreateProject(ctx, name, des, banner, content, slug, featured)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Project{}, fmt.Errorf("failed to save project: %w", err)
	}

	go s.NewsletterRepo.CreateProjectNewsletterJobs(project.Project_ID)

	return project, nil
}

func (s *ProjectService) GetAllProjects_Service(ctx context.Context, page int, offset int) ([]model.Project, int, error) {
	project, totalprojects, err := s.ProjectRepository.GetAllProjects(ctx, page, offset)
	if err != nil {
		log.Printf("[ERROR] failed to get projects from database: %v", err)
		return []model.Project{}, 0, model.ErrProjectNotFound
	}

	return project, totalprojects, nil

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

func (s *ProjectService) GetProjectBySlug_Service(ctx context.Context, slug string) (model.Project, error) {
	if slug == "" {
		log.Println("[WARN] Slug is Required")
		return model.Project{}, model.ErrSlugValidation
	}

	project, err := s.ProjectRepository.GetProjectBySlug(ctx, slug)
	if err != nil {
		log.Println("[Error] Error Geting project from database")
		return model.Project{}, model.ErrProjectNotFound
	}

	return project, nil
}

func (s *ProjectService) DeleteProjectById_Service(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return model.ErrIdValidation
	}

	if err := s.ProjectRepository.DeleteProjectById(ctx, id); err != nil {
		log.Printf("[ERROR] failed to delete project: %v", err)
		return err
	}

	return nil
}

func (s *ProjectService) UpdateProjectById_Service(ctx context.Context, id uuid.UUID, name string, des string, content string, banner string, slug string, featured bool) (model.Project, error) {
	if id == uuid.Nil {
		log.Println("[WARN] validation failed")
		return model.Project{}, model.ErrIdValidation
	}

	p, err := s.ProjectRepository.GetProjectById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] project not found: %v", err)
		return model.Project{}, model.ErrProjectNotFound
	}

	if name == "" {
		name = p.Project_Name
	}
	if des == "" {
		des = p.Project_DES
	}
	if content == "" {
		content = p.Project_Content
	}
	if banner == "" {
		banner = p.Project_Banner
	}
	if slug == "" {
		slug = p.Project_Slug
	}

	project, err := s.ProjectRepository.UpdateProjectById(ctx, id, name, des, content, banner, slug, featured)
	if err != nil {
		log.Printf("[ERROR] failed to update project from database: %v", err)
		return model.Project{}, model.ErrProjectNotFound
	}

	return project, nil
}
