package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	Success bool
	Message string
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
	Login_At time.Time `json:"login_at"`
}
