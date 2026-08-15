package reposiotry

import (
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reposiotry struct {
	pool *pgxpool.Pool
}

func (r *reposiotry) Blog_Repository(b *model.Blog) error {
	query := `
		INSERT INTO blogs (blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5, $6) 
	`

}
