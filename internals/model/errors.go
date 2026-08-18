package model

import "errors"

var (
	ErrValidation = errors.New("All fields are required")
	ErrSlugValidation = errors.New("Slug field ig required")
	ErrIdValidation = errors.New("Id field is required")
	ErrBlogNotFound   = errors.New("Blog not found")
)
