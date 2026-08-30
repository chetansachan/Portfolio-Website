package model

import (
	"time"

	"github.com/google/uuid"
)

type PaginationProjectMeta struct {
	CurrentPage   int  `json:"current_page"`
	PageSize      int  `json:"page_size"`
	TotalProjects int  `json:"total_projects"`
	TotalPages    int  `json:"total_pages"`
	HasNext       bool `json:"has_next"`
	HasPrev       bool `json:"has_prev"`
}

type PaginatedProjectResponse struct {
	ProjectResponse
	Data []Project             `json:"data"`
	Meta PaginationProjectMeta `json:"meta"`
}

type ProjectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetProjectResponse struct {
	ProjectResponse
	Data []Project `json:"data"`
}

type PostProjectResponse struct {
	ProjectResponse
	Data Project `json:"data"`
}

type Project struct {
	Project_ID      uuid.UUID `json:"project_id"`
	Project_Name    string    `json:"project_name"`
	Project_DES     string    `json:"project_des"`
	Project_Content string    `json:"project_content"`
	Project_Banner  string    `json:"project_banner"`
	Project_Slug    string    `json:"project_slug"`
	Is_Featured     bool      `json:"is_featured"`
	Created_At      time.Time `json:"created_at"`
}
