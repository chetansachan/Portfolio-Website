package reposiotry

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	Pool *pgxpool.Pool
}

func NewBlogRepository(pool *pgxpool.Pool) *BlogRepository {
	return &BlogRepository{Pool: pool}
}

func (r *BlogRepository) PostBlog(ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) error {
	query := `
	INSERT INTO blogs (blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.Pool.Exec(ctx, query, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner)
	if err != nil {
		log.Println("Error Executing Query")
		return err
	}

	return nil
}
