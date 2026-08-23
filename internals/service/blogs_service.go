package service

import (
	"context"
	"fmt"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"

	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepository *reposiotry.BlogRepository
	NewsletterRepo *reposiotry.NewsletterRepo
}

func NewBlogService(blogRepo *reposiotry.BlogRepository, newsletterRepo *reposiotry.NewsletterRepo) *BlogService {
	return &BlogService{BlogRepository: blogRepo, NewsletterRepo: newsletterRepo}
}

func (s *BlogService) CreateBlog_Service(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
	if title == "" || des == "" || slug == "" || content == "" || banner == "" {
		log.Println("[WARN] validation failed")
		return model.Blog{}, model.ErrValidation
	}

	b, err := s.BlogRepository.GetBlogBySlug(ctx, slug)
	if err == nil {
		log.Println("Blog with this slug is already exist")
		return b, model.ErrBlogAlreadyExist
	}

	blog, err := s.BlogRepository.CreateBlog(ctx, title, des, slug, content, banner)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Blog{}, fmt.Errorf("failed to save blog: %w", err)
	}

	go s.NewsletterRepo.CreateNewsletterJobs(blog.Blog_ID)

	return blog, nil
}

func (s *BlogService) GetAllBlogs_Service(ctx context.Context, limit, offset int) ([]model.Blog, int, error) {
	blogs, totalblogs, err := s.BlogRepository.GetAllBlogs(ctx, limit, offset)
	if err != nil {
		log.Printf("[ERROR] failed to get blogs from database: %v", err)
		return []model.Blog{}, 0, model.ErrBlogNotFound
	}

	return blogs, totalblogs, nil
}

func (s *BlogService) GetBlogBySlug_Service(ctx context.Context, slug string) (model.Blog, error) {
	if slug == "" {
		log.Println("[WARN] validation failed")
		return model.Blog{}, model.ErrSlugValidation
	}

	blog, err := s.BlogRepository.GetBlogBySlug(ctx, slug)
	if err != nil {
		log.Printf("[ERROR] failed to get blog from database: %v", err)
		return model.Blog{}, model.ErrBlogNotFound
	}

	return blog, nil
}

func (s *BlogService) GetBlogById_Service(ctx context.Context, id uuid.UUID) (model.Blog, error) {
	if id == uuid.Nil {
		log.Println("[WARN] validation failed")
		return model.Blog{}, model.ErrSlugValidation
	}

	blog, err := s.BlogRepository.GetBlogById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] failed to get blog from database: %v", err)
		return model.Blog{}, model.ErrBlogNotFound
	}

	return blog, nil
}

func (s *BlogService) DeleteBlogById_Service(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return model.ErrIdValidation
	}

	if err := s.BlogRepository.DeleteBlogById(ctx, id); err != nil {
		log.Printf("[ERROR] failed to delete blog: %v", err)
		return err
	}

	return nil
}

func (s *BlogService) UpdateBlogById_Service(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string) (model.Blog, error) {
	if id == uuid.Nil {
		log.Println("[WARN] validation failed")
		return model.Blog{}, model.ErrIdValidation
	}

	b, err := s.BlogRepository.GetBlogById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] blog not found: %v", err)
		return model.Blog{}, model.ErrBlogNotFound
	}

	if title == "" {
		title = b.Blog_Title
	}
	if des == "" {
		des = b.Blog_DES
	}
	if slug == "" {
		slug = b.Blog_Slug
	}
	if content == "" {
		content = b.Blog_Content
	}
	if banner == "" {
		banner = b.Blog_Banner
	}

	blog, err := s.BlogRepository.UpdateBlogById(ctx, id, title, des, slug, content, banner)
	if err != nil {
		log.Printf("[ERROR] failed to update blog from database: %v", err)
		return model.Blog{}, model.ErrBlogNotFound
	}

	return blog, nil
}
