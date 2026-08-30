package service

import (
	"context"
	"errors"
	"portfolio-website/internals/model"
	"testing"
)

func TestIsPasswordStrong(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"weakpassword", false},
		{"WEAKPASSWORD1", false},
		{"WeakPassword", false},
		{"WeakPassword1", false},
		{"WeakPassword!", false},
		{"WeakPassword1!", true},
		{"StrongPass123#", true},
		{"NOLOWER123!", false},
		{"noupper123!", false},
		{"NoDigitSpec", false},
	}

	for _, tc := range tests {
		result := isPasswordStrong(tc.password)
		if result != tc.expected {
			t.Errorf("isPasswordStrong(%q) = %v; want %v", tc.password, result, tc.expected)
		}
	}
}

func TestRegisterUser_Service_Validation(t *testing.T) {
	service := &AuthService{AuthRepo: nil}

	tests := []struct {
		name     string
		email    string
		password string
		expected error
	}{
		{"Name", "invalid-email", "StrongPass1!", model.ErrInvalidEmail},
		{"Name", "valid@email.com", "", model.ErrPasswordRequired},
		{"Name", "valid@email.com", "Short1!", model.ErrPasswordTooShort},
		{"Name", "valid@email.com", "weakpassword", model.ErrPasswordWeak},
		{"Name", "valid@email.com", "Weak12345678", model.ErrPasswordWeak},
	}

	for _, tc := range tests {
		err := service.RegisterUser_Service(context.Background(), tc.name, tc.email, tc.password)
		if !errors.Is(err, tc.expected) {
			t.Errorf("RegisterUser_Service(..., %q, %q) error = %v; want %v", tc.email, tc.password, err, tc.expected)
		}
	}
}
