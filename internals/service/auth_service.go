package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo *reposiotry.AuthRepo
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email string
	Role  string
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId string
}

func NewAuthService(authRepo *reposiotry.AuthRepo) *AuthService {
	return &AuthService{AuthRepo: authRepo}
}

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) RegisterUser_Service(ctx context.Context, name string, email string, password string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		fmt.Printf("%s is valid\n", email)
	} else {
		return fmt.Errorf("%s is invalid\n", email)
	}
	if password == "" {
		return fmt.Errorf("Password is requred Cannot be empty")
	}
	if len(password) <= 8 {
		return fmt.Errorf("Password length is aytleast 8 charaters")
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Println("Error hashing password")
	}
	password = string(hashpassword)

	err2 := s.AuthRepo.RegisterUser(ctx, email, password)
	if err2 != nil {
		return fmt.Errorf("Error Creating Admin")
	}

	return nil
}

func (s *AuthService) Login_Service(ctx context.Context, name string, email string, password string, ip string) (model.User, error) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		fmt.Printf("%s is valid\n", email)
	} else {
		return model.User{}, fmt.Errorf("%s is invalid\n", email)
	}

	user, err := s.AuthRepo.Login(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("Error finding Admin")
	}

	if err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err2 != nil {
		return model.User{}, fmt.Errorf("Password Incorrect")
	}

	// Access token
	access_token, err := GernateToken(ctx, user, false)

	// Refresh token
	refresh_token, err := GernateToken(ctx, user, true)
	hashrefreshtkn := HashToken(refresh_token)

	expired_at := time.Now().Add(7 * 24 * time.Hour)
	if err := s.SaveRefreshToken(ctx, user, ip, hashrefreshtkn, expired_at); err != nil {
		return model.User{}, fmt.Errorf("Error Saving refresh token")
	}

	fmt.Println(access_token)
	fmt.Println(refresh_token)

	return user, nil
}

func GernateToken(ctx context.Context, user model.User, isRefeshToken bool) (string, error) {
	secretStr := os.Getenv("MY_JWT_SECRETKEY")
	mySecretKey := []byte(secretStr)

	var claims jwt.Claims
	if !isRefeshToken {
		claims = &AccessTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.Email,
			"Admin",
		}
	} else {
		claims = &RefreshTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.UserId.String(),
		}
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(t)
	token, err := t.SignedString(mySecretKey)
	if err != nil {
		return "", fmt.Errorf("Error Gernating token")
	}

	return token, nil
}

func (s *AuthService) SaveRefreshToken(ctx context.Context, user model.User, ipAddress string, hashrefreshtkn string, expired_at time.Time) error {
	if err := s.AuthRepo.SaveRefreshTkn(ctx, hashrefreshtkn, user.UserId, user.Name, ipAddress, expired_at); err != nil {
		log.Println("Failed to save token")
		return err
	}

	return nil
}
