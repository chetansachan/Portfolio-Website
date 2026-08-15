package service

import (
	"context"
	"log"
	"portfolio-website/internals/reposiotry"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PostBlog_Service(pool *pgxpool.Pool, ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) error {
	if Blog_Title == "" || Blog_DES == "" || Blog_Slug == "" || Blog_Content == "" || Blog_Banner == "" {
		log.Println("Every field is required")
		return nil
	}

	if err := reposiotry.PostBlog(pool, ctx, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner); err != nil {
		log.Println("Error Executing Query")
		return nil
	}

	return nil
}
