package service

import (
	"context"
	"log"
	"portfolio-website/internals/reposiotry"
)

func PostBlog_Service(ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) error {
	if Blog_Title == "" || Blog_DES == "" || Blog_Slug == "" || Blog_Content == "" || Blog_Banner == "" {
		log.Println("Every field is required")
		return nil
	}

	if err := reposiotry.PostBlog(ctx, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner); err != nil {
		log.Println("Error Executing Query")
		return nil
	}

	return nil
}
