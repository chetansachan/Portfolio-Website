package model

import (
	"portfolio-website/internals/service"
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	Success bool
	Message string
}

type LoginResponseWithTkn struct {
	Success bool
	Message string
	Token service.LoginResponse
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
