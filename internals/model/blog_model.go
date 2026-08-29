package model

import (
	"time"

	"github.com/google/uuid"
)

type PaginationMeta struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	TotalBlogs int  `json:"total_blogs"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

type PaginatedResponse struct {
	Response
	Data []Blog         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type Response struct {
	Success bool	`json:"success"`
	Message string `json:"message"`
}

type GetBlogsResponse struct {
	Response
	Data []Blog
}

type PostBlogResponse struct {
	Response
	Data Blog `json:"data"`
}

type Blog struct {
	Blog_ID      uuid.UUID `json:"blog_id"`
	Blog_Title   string    `json:"blog_title"`
	Blog_DES     string    `json:"blog_des"`
	Blog_Content string    `json:"blog_content"`
	Blog_Slug    string    `json:"blog_slug"`
	Blog_Banner  string    `json:"blog_banner"`
	Is_Featured  bool      `json:"is_featured"`
	Created_At   time.Time `json:"created_at"`
}
