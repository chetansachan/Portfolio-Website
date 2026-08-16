package model

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Success bool
	Message string
}

type GetBlogsResponse struct {
	Response
	Data Blog
}

type PostBlogResponse struct {
	Response
}

type Blog struct {
	Blog_ID      uuid.UUID `json:"blog_id"`
	Blog_Title   string    `json:"blog_title"`
	Blog_DES     string    `json:"blog_des"`
	Blog_Content string    `json:"blog_content"`
	Blog_Slug    string    `json:"blog_slug"`
	Blog_Banner  string    `json:"blog_banner"`
	Created_At   time.Time `json:"created_at"`
}
