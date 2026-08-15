package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reposiotry struct {
	pool *pgxpool.Pool
}

func (r *Reposiotry) PostBlog (b *model.Blog) error {
	query := `
		INSERT INTO blogs (blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5, $6) 
	`

	ctx := context.Background()
	_, err := r.pool.Exec(ctx, query, &b.Blog_ID, &b.Blog_Title, &b.Blog_DES, &b.Blog_Slug, &b.Blog_Content, &b.Blog_Banner)
	if err != nil {
		log.Println("Error Executing Query")
		return err
	}

	return nil
}
