package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"

	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepository *reposiotry.BlogRepository
}

func NewBlogService(blogRepo *reposiotry.BlogRepository) *BlogService {
	return &BlogService{BlogRepository: blogRepo}
}

func (s *BlogService) PostBlog_Service(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
	if title == "" || des == "" || slug == "" || content == "" || banner == "" {
		log.Println("[WARN] validation failed")
		return model.Blog{}, model.ErrValidation
	}

	b, err := s.BlogRepository.GetBlogBySlug(ctx, slug)
	if err == nil {
		log.Println("Blog with this slug is already exist")
		return b, nil
	}

	blog, err := s.BlogRepository.PostBlog(ctx, title, des, slug, content, banner)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Blog{}, fmt.Errorf("failed to save blog: %w", err)
	}

	return blog, nil
}

func (s *BlogService) GetBlogBySlug_Service(ctx context.Context, slug string) (model.Blog, error) {
	if slug == "" {
		log.Println("[WARN] validation failed")
		return model.Blog{}, fmt.Errorf("Slug is Required: %s", slug)
	}

	blog, err := s.BlogRepository.GetBlogBySlug(ctx, slug)
	if err != nil {
		log.Printf("[ERROR] failed to gec blog from database: %v", err)
		return model.Blog{}, fmt.Errorf("Blog Not Found: %w", err)
	}

	return blog, nil
}

func (s *BlogService) DeleteBlogById_Service(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		errors.New("id is required to delete the blog")
	}

	if err := s.BlogRepository.DeleteBlogById(ctx, id); err != nil {
		log.Println("Blog Not Found")
		return err
	}

	return nil
}

func (s *BlogService) UpdateBlogById_Service(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string) (model.Blog, error) {
	blog, err := s.BlogRepository.UpdateBlogById(ctx, id, title, des, slug, content, banner)
	if err != nil {
		log.Println("Blog Not Found")
		return model.Blog{}, err
	}

	return blog, nil
}
