package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Project_ID     uuid.UUID `json:"project_id"`
	Project_Name   string    `json:"project_name"`
	Project_DES    string    `json:"project_des"`
	Project_URL    string    `json:"project_url"`
	Project_Banner string    `json:"project_vid_url"`
	Is_Featured    bool      `json:"is_featured"`
	Created_At     time.Time `json:"created_at"`
}
