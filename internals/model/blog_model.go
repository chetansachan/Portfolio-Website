package model

import (
	"github.com/google/uuid"
)

type Blog struct {
	Blog_id      uuid.UUID `json:"blog_id"`
	Blog_name    string    `json:"blog_name"`
	Blog_des     string    `json:"blog_des"`
	Blog_content string    `json:"blog_content"`
	Created_at   string    `json:"created_at"`
}
