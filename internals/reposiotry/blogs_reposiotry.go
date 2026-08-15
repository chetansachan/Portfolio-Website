package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PostBlog(pool *pgxpool.Pool,ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) error {
	query := `
	`

	if err != nil {
		log.Println("Error Executing Query")
		return err
	}

	return nil
}
