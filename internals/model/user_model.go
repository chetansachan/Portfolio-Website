package model

import (
	"time"

	"github.com/google/uuid"
)

type TokenDetail struct {
	AccessToken  string
	RefreshToken string
}

type LoginResponse struct {
	Success bool
	Message string
}

type LoginResponseWithTkn struct {
	Success bool
	Message string
	Token   TokenDetail
}

type GetLoginResponse struct {
	Response LoginResponse
	Data     User `json:"data"`
}

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Email    string    `validate:"required,email"`
	Password string    `validate:"required,password"`
	Is_Admin bool      `json:"is_admin"`
	Login_At time.Time `json:"login_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	IsRevoked bool      `json:"is_revoked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

