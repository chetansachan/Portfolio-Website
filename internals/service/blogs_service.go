package service

import (
	"context"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
)

type BlogService struct {
	BlogRepository *reposiotry.BlogRepository
}

func NewBlogService(blogRepo *reposiotry.BlogRepository) *BlogService {
	return &BlogService{BlogRepository: blogRepo}
}

func (s *BlogService) PostBlog_Service(ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) (model.Blog, error) {
	if Blog_Title == "" || Blog_DES == "" || Blog_Slug == "" || Blog_Content == "" || Blog_Banner == "" {
		log.Println("Every field is required")
		return model.Blog{}, nil
	}

	blog, err := s.BlogRepository.PostBlog(ctx, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner)
	if err != nil {
		log.Println("Error Executing Query")
		return model.Blog{}, nil
	}

	return blog, nil
}

func (s *BlogService) GetBlogBySlug_Service(ctx context.Context, slug string) (model.Blog, error) {
	
}
