package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type NewsletterController struct {
	NewsletterService *service.NewsletterService
}

func NewNewsletterController(newsletterService *service.NewsletterService) *NewsletterController {
	return &NewsletterController{NewsletterService: newsletterService}
}

func (c *NewsletterController) CreateSubscriber_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Invalid request body"})
		return
	}

	if err := c.NewsletterService.CreateSubscriber_Service(r.Context(), req.Email); err != nil {
		if errors.Is(err, model.ErrInvalidEmail) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid email",
			})
			return
		}
		log.Printf("[ERROR] CreateSubscriber failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Subscribed",
	})
}
