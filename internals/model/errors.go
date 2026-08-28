package model

import "errors"

var (
	ErrValidation          = errors.New("All fields are required")
	ErrSlugValidation      = errors.New("Slug field ig required")
	ErrIdValidation        = errors.New("Id field is required")
	ErrBlogNotFound        = errors.New("Blog not found")
	ErrBlogAlreadyExist    = errors.New("Blog already exist")
	ErrProjectNotFound     = errors.New("Project not found")
	ErrProjectAlreadyExist = errors.New("Project already exist")
	ErrCaseStudyNotFound     = errors.New("Case Study not found")
	ErrCaseStudyAlreadyExist = errors.New("Case Study already exist")
	ErrFileNotFound          = errors.New("File not found")
)
