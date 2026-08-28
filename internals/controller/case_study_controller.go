package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"portfolio-website/internals/model"
	"strconv"

	"github.com/google/uuid"
)

type CaseStudyServiceInterface interface {
	CreateCaseStudy_Service(ctx context.Context, title, des, slug, content, banner string, featured bool) (model.Case_Study, error)
	GetAllCaseStudies_Service(ctx context.Context, size int, offset int) ([]model.Case_Study, int, error)
	GetCaseStudyBySlug_Service(ctx context.Context, slug string) (model.Case_Study, error)
	GetCaseStudyById_Service(ctx context.Context, id uuid.UUID) (model.Case_Study, error)
	DeleteCaseStudyById_Service(ctx context.Context, id uuid.UUID) error
	UpdateCaseStudyById_Service(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string, featured bool) (model.Case_Study, error)
	GetFeaturedCaseStudies_Service(ctx context.Context) ([]model.Case_Study, error)
}

type CaseStudyController struct {
	CaseStudyService CaseStudyServiceInterface
}

func NewCaseStudyController(caseStudyService CaseStudyServiceInterface) *CaseStudyController {
	return &CaseStudyController{CaseStudyService: caseStudyService}
}

func (c *CaseStudyController) CreateCaseStudy_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var cs model.Case_Study
	if err := json.NewDecoder(r.Body).Decode(&cs); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Invalid request body"})
		return
	}

	case_study, err := c.CaseStudyService.CreateCaseStudy_Service(r.Context(), cs.Case_Study_Title, cs.Case_Study_DES, cs.Case_Study_Slug, cs.Case_Study_Content, cs.Case_Study_Banner, cs.Is_Featured)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "All fields are required",
			})

		case errors.Is(err, model.ErrCaseStudyAlreadyExist):
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Case Study Already exist",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.PostCaseStudyResponse{
		CaseStudyResponse: model.CaseStudyResponse{Success: true, Message: "Case Study created successfully"},
		Data:              case_study,
	})
}

func (c *CaseStudyController) GetAllCaseStudies_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	queryParams := r.URL.Query()
	pageStr := queryParams.Get("page")
	sizeStr := queryParams.Get("size")

	page := 1
	size := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Invalid 'page' parameter: must be a positive integer",
			})
			return
		}
		page = p
	}

	if sizeStr != "" {
		s, err := strconv.Atoi(sizeStr)
		if err != nil || s <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Invalid 'size' parameter: must be a positive integer",
			})
			return
		}
		if s > 100 {
			s = 100
		}
		size = s
	}

	offset := (page - 1) * size

	case_studies, totalCaseStudies, err := c.CaseStudyService.GetAllCaseStudies_Service(r.Context(), size, offset)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "No Case Studies Found",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	totalPages := 0
	if totalCaseStudies > 0 {
		totalPages = int(math.Ceil(float64(totalCaseStudies) / float64(size)))
	}

	response := model.Case_StudyPaginatedResponse{
		CaseStudyResponse: model.CaseStudyResponse{
			Success: true,
			Message: "Case Studies Found Successfully",
		},
		Data: case_studies,
		Meta: model.Case_StudyPaginationMeta{
			CurrentPage:      page,
			PageSize:         size,
			TotalCaseStudies: totalCaseStudies,
			TotalPages:       totalPages,
			HasNext:          page < totalPages,
			HasPrev:          page > 1 && totalCaseStudies > 0,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *CaseStudyController) GetCaseStudyBySlug_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	slug := r.PathValue("slug")
	case_study, err := c.CaseStudyService.GetCaseStudyBySlug_Service(r.Context(), slug)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSlugValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Slug field is required",
			})

		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Case study Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	response := model.GetCaseStudyResponse{
		CaseStudyResponse: model.CaseStudyResponse{
			Success: true,
			Message: "Case Study Found Successfully",
		},
		Data: []model.Case_Study{case_study},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(response)
}

func (c *CaseStudyController) GetCaseStudyById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	case_study, err := c.CaseStudyService.GetCaseStudyById_Service(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Valid ID is required",
			})
		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Case study not found",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(model.PostCaseStudyResponse{
		CaseStudyResponse: model.CaseStudyResponse{
			Success: true,
			Message: "Case Study Found Successfully",
		},
		Data: case_study,
	})
}

func (c *CaseStudyController) DeleteCaseStudyById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := c.CaseStudyService.DeleteCaseStudyById_Service(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Case Study Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.CaseStudyResponse{
		Success: true,
		Message: "Case Study Deleted Successfully",
	})
}

func (c *CaseStudyController) UpdateCaseStudyById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var cs model.Case_Study
	if err := json.NewDecoder(r.Body).Decode(&cs); err != nil {
		log.Printf("Error Decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{
			Success: false,
			Message: "Error Decoding JSON",
		})
		return
	}

	case_study, err := c.CaseStudyService.UpdateCaseStudyById_Service(r.Context(), id, cs.Case_Study_Title, cs.Case_Study_DES, cs.Case_Study_Slug, cs.Case_Study_Content, cs.Case_Study_Banner, cs.Is_Featured)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Id field is required",
			})

		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Case Study Not Found",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.PostCaseStudyResponse{
		CaseStudyResponse: model.CaseStudyResponse{
			Success: true,
			Message: "Case Study Updated Successfully",
		},
		Data: case_study,
	})
}

func (c *CaseStudyController) GetFeaturedCaseStudies_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.CaseStudyResponse{Success: false, Message: "Method not allowed"})
		return
	}

	caseStudies, err := c.CaseStudyService.GetFeaturedCaseStudies_Service(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrCaseStudyNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "No Featured Case Studies Found",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.CaseStudyResponse{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.GetCaseStudyResponse{
		CaseStudyResponse: model.CaseStudyResponse{
			Success: true,
			Message: "Featured case studies retrieved successfully",
		},
		Data: caseStudies,
	})
}
