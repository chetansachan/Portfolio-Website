package service

import (
	"context"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
)

type ProjectService struct {
	ProjectRepository *reposiotry.ProjectRepository
}

func NewProjectService(projectrepository *reposiotry.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectrepository}
}

func (s *ProjectService) PostProject_Service(ctx context.Context, name string, des string, banner string, content string, featured bool) (model.Project, error) {
	if name == "" || des == "" || banner == "" || content == "" {
		log.Println("Every field is required")
		return model.Project{}, nil
	}

	project, err := s.ProjectRepository.PostProject(ctx, name, des, banner, content, featured)
	if err != nil {
		log.Println("Error Executing Query")
		return model.Project{}, nil
	}

	return project, nil
}
