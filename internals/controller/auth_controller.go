package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
	"time"

	"github.com/jackc/pgx/v5"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (c *AuthController) RegisterUser_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid JSON Payload",
		})
		return
	}

	err := c.AuthService.RegisterUser_Service(r.Context(), user.Name, user.Email, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Invalid Email",
			})
		case errors.Is(err,model.ErrUserAlreadyExist):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "User Already Exist",
			})
		case errors.Is(err, model.ErrPasswordRequired):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Password is required",
			})
		case errors.Is(err, model.ErrPasswordTooShort):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Password too short",
			})
		case errors.Is(err, model.ErrPasswordWeak):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Password too weak",
			})
		default:
			log.Printf("[ERROR] RegisterUser failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "User Registered Successfully",
	},
	)
}

func (c *AuthController) Login_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid JSON Payload",
		})
		return
	}

	tokens, err := c.AuthService.Login_Service(r.Context(), u.Email, u.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: err.Error(),
			})
		default:
			if errors.Is(err, pgx.ErrNoRows) || err.Error() == "password incorrect" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(model.Response{
					Success: false,
					Message: "Invalid email or password",
				})
				return
			}
			log.Printf("[ERROR] Login failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.LoginResponse{
		Success: true,
		Message: "Login Successful",
		Tokens: model.TokenDetail{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	},
	)
}

func (c *AuthController) RefreshToken_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.LoginResponse{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.LoginResponse{
			Success: false,
			Message: "Invalid JSON Payload",
		})
		return
	}

	if req.RefreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.LoginResponse{
			Success: false,
			Message: "Refresh token is required",
		})
		return
	}

	tokens, err := c.AuthService.RotateTokens_Service(r.Context(), req.RefreshToken)
	if err != nil {
		log.Printf("[WARN] Refresh token rotation failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.LoginResponse{
			Success: false,
			Message: "Invalid or expired refresh token",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.LoginResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Tokens: model.TokenDetail{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
