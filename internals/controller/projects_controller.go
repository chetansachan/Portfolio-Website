package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"
)

type ProjectServiceInterface interface {
	PostProject_Service(ctx context.Context, name string, des string, banner string, content string, slug string, featured bool) (model.Project, error)
}

type ProjectController struct {
	ProjectService ProjectServiceInterface
}

func NewProjectController(projectservice ProjectServiceInterface) *ProjectController {
	return &ProjectController{ProjectService: projectservice}
}

func (c *ProjectController) PostProject_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var p model.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Invalid request body"})
		return
	}

	project, err := c.ProjectService.PostProject_Service(r.Context(), p.Project_Name, p.Project_DES, p.Project_Banner, p.Project_Content, p.Project_Slug, p.Is_Featured)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "All fields are required",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.PostProjectResponse{
		ProjectResponse: model.ProjectResponse{
			Success: true,
			Message: "Project created successfully"},
		Data: project,
	})
}
