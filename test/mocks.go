package test

import (
	"context"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

// =============================================================================
// Mock Blog Service
// =============================================================================

type MockBlogService struct {
	CreateBlogFunc    func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error)
	GetAllBlogsFunc   func(ctx context.Context, size int, offset int) ([]model.Blog, int, error)
	GetBlogBySlugFunc func(ctx context.Context, slug string) (model.Blog, error)
	GetBlogByIdFunc   func(ctx context.Context, id uuid.UUID) (model.Blog, error)
	DeleteBlogFunc    func(ctx context.Context, id uuid.UUID) error
	UpdateBlogFunc    func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error)
	PostFunc   func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error)
	DeleteFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *MockBlogService) CreateBlog_Service(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
	if m.CreateBlogFunc != nil {
		return m.CreateBlogFunc(ctx, title, des, slug, content, banner)
	}
	if m.PostFunc != nil {
		return m.PostFunc(ctx, title, des, slug, content, banner)
	}
	return model.Blog{}, nil
}

func (m *MockBlogService) GetAllBlogs_Service(ctx context.Context, size int, offset int) ([]model.Blog, int, error) {
	if m.GetAllBlogsFunc != nil {
		return m.GetAllBlogsFunc(ctx, size, offset)
	}
	return nil, 0, nil
}

func (m *MockBlogService) GetBlogBySlug_Service(ctx context.Context, slug string) (model.Blog, error) {
	if m.GetBlogBySlugFunc != nil {
		return m.GetBlogBySlugFunc(ctx, slug)
	}
	return model.Blog{}, nil
}

func (m *MockBlogService) GetBlogById_Service(ctx context.Context, id uuid.UUID) (model.Blog, error) {
	if m.GetBlogByIdFunc != nil {
		return m.GetBlogByIdFunc(ctx, id)
	}
	return model.Blog{}, nil
}

func (m *MockBlogService) DeleteBlogById_Service(ctx context.Context, id uuid.UUID) error {
	if m.DeleteBlogFunc != nil {
		return m.DeleteBlogFunc(ctx, id)
	}
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockBlogService) UpdateBlogById_Service(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
	if m.UpdateBlogFunc != nil {
		return m.UpdateBlogFunc(ctx, id, title, des, slug, content, banner)
	}
	return model.Blog{}, nil
}

// =============================================================================
// Mock Project Service
// =============================================================================

type MockProjectService struct {
	CreateProjectFunc    func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error)
	GetAllProjectsFunc   func(ctx context.Context, page int, offset int) ([]model.Project, int, error)
	GetProjectByIdFunc   func(ctx context.Context, id uuid.UUID) (model.Project, error)
	GetProjectBySlugFunc func(ctx context.Context, slug string) (model.Project, error)
	DeleteProjectFunc    func(ctx context.Context, id uuid.UUID) error
	UpdateProjectFunc    func(ctx context.Context, id uuid.UUID, name, des, content, banner, slug string, featured bool) (model.Project, error)
}

func (m *MockProjectService) CreateProject_Service(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, name, des, banner, content, slug, featured)
	}
	return model.Project{}, nil
}

func (m *MockProjectService) GetAllProjects_Service(ctx context.Context, page int, offset int) ([]model.Project, int, error) {
	if m.GetAllProjectsFunc != nil {
		return m.GetAllProjectsFunc(ctx, page, offset)
	}
	return nil, 0, nil
}

func (m *MockProjectService) GetProjectById_Service(ctx context.Context, id uuid.UUID) (model.Project, error) {
	if m.GetProjectByIdFunc != nil {
		return m.GetProjectByIdFunc(ctx, id)
	}
	return model.Project{}, nil
}

func (m *MockProjectService) GetProjectBySlug_Service(ctx context.Context, slug string) (model.Project, error) {
	if m.GetProjectBySlugFunc != nil {
		return m.GetProjectBySlugFunc(ctx, slug)
	}
	return model.Project{}, nil
}

func (m *MockProjectService) DeleteProjectById_Service(ctx context.Context, id uuid.UUID) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectService) UpdateProjectById_Service(ctx context.Context, id uuid.UUID, name, des, content, banner, slug string, featured bool) (model.Project, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, id, name, des, content, banner, slug, featured)
	}
	return model.Project{}, nil
}
