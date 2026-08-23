package service

import (
	"context"
	"fmt"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	AdminRepo *reposiotry.AdminRepo
}

func NewAdminService(adminRepo *reposiotry.AdminRepo) *AdminService {
	return &AdminService{AdminRepo: adminRepo}
}

func (s *AdminService) CreateAdmin_Service(ctx context.Context, email string, password string) (model.Admin, error) {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		fmt.Printf("%s is valid\n", email)
	} else {
		return model.Admin{}, fmt.Errorf("%s is invalid\n", email)
	}

	if password == "" {
		return model.Admin{}, fmt.Errorf("Password is requred Cannot be empty")
	}
	if len(password) <= 8 {
		return model.Admin{}, fmt.Errorf("Password length is aytleast 8 charaters")
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Println("Error hashing password")
	}
	password = string(hashpassword)

	admin, err := s.AdminRepo.CreateAdmin(ctx, email, password)
	if err != nil {
		return model.Admin{}, fmt.Errorf("Error Creating Admin")
	}

	return admin, nil
}
