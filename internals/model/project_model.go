package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectResponse struct {
	Success bool
	Message string
}

type GetProjectResponse struct {
	ProjectResponse
	Data []Project
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
	Project_Slug 	string	  `json:"project_slug"`	 
	Is_Featured     bool      `json:"is_featured"`
	Created_At      time.Time `json:"created_at"`
}
