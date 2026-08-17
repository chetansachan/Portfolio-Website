package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type ProjectController struct {
	ProjectService *service.ProjectService
}

func NewProjectController(projectservice *service.ProjectService) *ProjectController {
	return &ProjectController{ProjectService: projectservice}
}

func (c *ProjectController) PostProject_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid Method")
		return
	}

	var p model.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	project, err := c.ProjectService.PostProject_Service(r.Context(), p.Project_Name, p.Project_DES, p.Project_Banner, p.Project_Content, p.Is_Featured)
	if err != nil {
		log.Println(err)
		return
	}

	response := model.PostProjectResponse{
		ProjectResponse: model.ProjectResponse{
			Success: true,
			Message: "Project Posted Successfully",
		},
		Data: model.Project{
			Project_ID:      project.Project_ID,
			Project_Name:    project.Project_Name,
			Project_DES:     project.Project_DES,
			Project_Content: project.Project_Content,
			Project_Banner:  project.Project_Banner,
			Is_Featured:     project.Is_Featured,
			Created_At:      project.Created_At,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error Encoding JSON")
		return
	}
}
