package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

type ProjectServiceInterface interface {
	CreateProject_Service(ctx context.Context, name string, des string, banner string, content string, slug string, featured bool) (model.Project, error)
	GetProjectById_Service(ctx context.Context, id uuid.UUID) (model.Project, error)
	GetProjectBySlug_Service(ctx context.Context, slug string) (model.Project, error)
	DeleteProjectById_Service(ctx context.Context, id uuid.UUID) error
	UpdateProjectById_Service(ctx context.Context, id uuid.UUID, name string, des string, content string, banner string, slug string, featured bool) (model.Project, error)
}

type ProjectController struct {
	ProjectService ProjectServiceInterface
}

func NewProjectController(projectservice ProjectServiceInterface) *ProjectController {
	return &ProjectController{ProjectService: projectservice}
}

func (c *ProjectController) CreateProject_Controller(w http.ResponseWriter, r *http.Request) {
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

	project, err := c.ProjectService.CreateProject_Service(r.Context(), p.Project_Name, p.Project_DES, p.Project_Banner, p.Project_Content, p.Project_Slug, p.Is_Featured)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "All fields are required",
			})

		case errors.Is(err, model.ErrProjectAlreadyExist):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Project Already Exist",
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

func (c *ProjectController) GetProjectById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	myId := r.PathValue("id")
	id, err := uuid.Parse(myId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Invalid Id Format",
		})
		return
	}

	project, err := c.ProjectService.GetProjectById_Service(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Id is Required",
			})

		case errors.Is(err, model.ErrProjectNotFound):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Project Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(model.PostProjectResponse{
		ProjectResponse: model.ProjectResponse{
			Success: true,
			Message: "Project Found Successfully",
		},
		Data: project,
	})
}

func (c *ProjectController) GetProjectBySlug_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	slug := r.PathValue("slug")
	project, err := c.ProjectService.GetProjectBySlug_Service(r.Context(), slug)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSlugValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Slug is Required",
			})

		case errors.Is(err, model.ErrProjectNotFound):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Project Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ProjectResponse{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(model.PostProjectResponse{
		ProjectResponse: model.ProjectResponse{
			Success: true,
			Message: "Project Found Successfully",
		},
		Data: project,
	})
}

func (c *ProjectController) DeleteProjectById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ProjectResponse{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := c.ProjectService.DeleteProjectById_Service(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, model.ErrProjectNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Project Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal server error",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Project Deleted Successfully",
	})
}

func (c *ProjectController) UpdateProjectById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ProjectResponse{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ProjectResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var p model.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	project, err := c.ProjectService.UpdateProjectById_Service(r.Context(), id, p.Project_Name, p.Project_DES, p.Project_Content, p.Project_Banner, p.Project_Slug, p.Is_Featured)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Id field is required",
			})

		case errors.Is(err, model.ErrProjectNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Project Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal server error",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.PostProjectResponse{
		ProjectResponse: model.ProjectResponse{
			Success: true,
			Message: "Project Updated Successfully",
		},
		Data: project,
	})
}
