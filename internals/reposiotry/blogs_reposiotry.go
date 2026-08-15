package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PostBlog(pool *pgxpool.Pool, ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) error {
	query := `
	INSERT INTO blogs (blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5, $6)
	`

	var b model.Blog
	_, err := pool.Exec(ctx, query, &b.Blog_Title, &b.Blog_DES, &b.Blog_Slug, &b.Blog_Content, &b.Blog_Banner)
	if err != nil {
		log.Println("Error Executing Query")
		return err
	}

	return nil
}
