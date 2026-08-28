package model

import (
	"time"

	"github.com/google/uuid"
)

type Case_StudyPaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalCaseStudies  int  `json:"total_case_studies"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

type Case_StudyPaginatedResponse struct {
	CaseStudyResponse
	Data []Case_Study             `json:"data"`
	Meta Case_StudyPaginationMeta `json:"meta"`
}

type CaseStudyResponse struct {
	Success bool
	Message string
}

type GetCaseStudyResponse struct {
	CaseStudyResponse
	Data []Case_Study
}

type PostCaseStudyResponse struct {
	CaseStudyResponse
	Data Case_Study `json:"data"`
}

type Case_Study struct {
	Case_Study_ID      uuid.UUID `json:"case_study_id"`
	Case_Study_Title   string    `json:"case_study_title"`
	Case_Study_DES     string    `json:"case_study_des"`
	Case_Study_Content string    `json:"case_study_content"`
	Case_Study_Slug    string    `json:"case_study_slug"`
	Case_Study_Banner  string    `json:"case_study_banner"`
	Is_Featured        bool      `json:"is_featured"`
	Created_At         time.Time `json:"created_at"`
}
