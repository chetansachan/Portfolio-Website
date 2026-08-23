package model

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	AdminId  uuid.UUID `json:"adminid"`
	Email    string    `validate:"required,email"`
	Password string    `validate:"required,password"`
	Login_At time.Time `json:"login_at"`
}
