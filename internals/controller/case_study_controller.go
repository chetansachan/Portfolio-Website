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
)

type CaseStudyServiceInterface interface {
	CreateCaseStudy_Service(ctx context.Context, title, des, slug, content, banner string) (model.Case_Study, error)
	GetAllCaseStudies_Service(ctx context.Context, size int, offset int) ([]model.Case_Study, int, error)
	GetCaseStudyBySlug_Service(ctx context.Context, slug string) (model.Case_Study, error)
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

	case_study, err := c.CaseStudyService.CreateCaseStudy_Service(r.Context(), cs.Case_Study_Title, cs.Case_Study_DES, cs.Case_Study_Slug, cs.Case_Study_Content, cs.Case_Study_Banner)
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

