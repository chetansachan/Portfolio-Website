package reposiotry

import (
	"context"
	"database/sql"
	"log"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	Pool *pgxpool.Pool
}

func NewBlogRepository(pool *pgxpool.Pool) *BlogRepository {
	return &BlogRepository{Pool: pool}
}

func (r *BlogRepository) PostBlog(ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) (model.Blog, error) {
	query := `
	INSERT INTO blogs (blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5) RETURNING blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner, created_at
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return blog, nil
}

func (r *BlogRepository) GetBlogBySlug(ctx context.Context, slug string) (model.Blog, error) {
	query := `
		SELECT * FROM blogs WHERE blog_slug = $1
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, slug).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err == sql.ErrNoRows {
		log.Println("No user found with that slug")
		return model.Blog{}, nil
	} else if err != nil {
		log.Printf("Database error: %v", err)
	}

	return blog, nil
}

func (r *BlogRepository) DeleteBlogById(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM blogs WHERE blog_id = $1
	`

	_, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Database error: %v", err)
	}

	return nil
}

func (r *BlogRepository) UpdateBlogById(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string) (model.Blog, error) {
	query := `
		UPDATE blogs SET blog_title = $2, blog_des = $3, blog_slug = $4, blog_content = $5, blog_banner = $6
		WHERE blog_id = $1
		RETURNING blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner, created_at
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, id, title, des, slug, content, banner).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return blog, nil
}
