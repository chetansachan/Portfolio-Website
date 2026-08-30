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
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo *reposiotry.AuthRepo
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email   string
	IsAdmin bool
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

	if !re.MatchString(email) {
		return fmt.Errorf("%s is invalid", email)
	}
	if password == "" {
		return fmt.Errorf("password is required and cannot be empty")
	}
	if len(password) < 8 {
		return fmt.Errorf("password length must be at least 8 characters")
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if err := s.AuthRepo.RegisterUser(ctx, name, email, string(hashpassword)); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (s *AuthService) Login_Service(ctx context.Context, email string, password string) (model.TokenDetail, error) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return model.TokenDetail{}, fmt.Errorf("%s is invalid", email)
	}

	user, err := s.AuthRepo.Login(ctx, email)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("error finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.TokenDetail{}, fmt.Errorf("password incorrect")
	}

	accessToken, err := s.GenerateToken(ctx, user, false)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := s.GenerateToken(ctx, user, true)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return model.TokenDetail{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, user model.User, isRefreshToken bool) (string, error) {
	secretStr := os.Getenv("MY_JWT_SECRETKEY")
	if secretStr == "" {
		return "", fmt.Errorf("MY_JWT_SECRETKEY is not set")
	}
	mySecretKey := []byte(secretStr)

	var claims jwt.Claims
	if !isRefreshToken {
		claims = &AccessTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Email:   user.Email,
			IsAdmin: user.Is_Admin,
		}
	} else {
		claims = &RefreshTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			UserId: user.UserId.String(),
		}
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(mySecretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	if isRefreshToken {
		hashRefreshTkn := HashToken(token)
		expiredAt := time.Now().Add(7 * 24 * time.Hour)
		if err := s.SaveRefreshToken(ctx, user, hashRefreshTkn, expiredAt); err != nil {
			return "", fmt.Errorf("error saving refresh token: %w", err)
		}
	}

	return token, nil
}

func (s *AuthService) SaveRefreshToken(ctx context.Context, user model.User, hashRefreshTkn string, expiredAt time.Time) error {
	if err := s.AuthRepo.SaveRefreshTkn(ctx, hashRefreshTkn, user.UserId, user.Name, expiredAt); err != nil {
		log.Printf("Failed to save refresh token for user %s: %v", user.Email, err)
		return err
	}
	return nil
}

func (s *AuthService) RotateTokens_Service(ctx context.Context, refreshToken string) (model.TokenDetail, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		secretStr := os.Getenv("MY_JWT_SECRETKEY")
		if secretStr == "" {
			return nil, fmt.Errorf("MY_JWT_SECRETKEY is not set")
		}
		return []byte(secretStr), nil
	})

	if err != nil || !token.Valid {
		return model.TokenDetail{}, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return model.TokenDetail{}, fmt.Errorf("invalid refresh token claims")
	}

	hashedTkn := HashToken(refreshToken)
	storedTkn, err := s.AuthRepo.GetRefreshToken(ctx, hashedTkn)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("refresh token not found or already used: %w", err)
	}

	if storedTkn.IsRevoked {
		return model.TokenDetail{}, fmt.Errorf("refresh token has been revoked")
	}

	if storedTkn.ExpiresAt.Before(time.Now()) {
		return model.TokenDetail{}, fmt.Errorf("refresh token has expired")
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("invalid user id in token: %w", err)
	}

	user, err := s.AuthRepo.GetUserByID(ctx, userID)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("failed to retrieve user details: %w", err)
	}

	newAccessToken, err := s.GenerateToken(ctx, user, false)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("error generating new access token: %w", err)
	}

	newRefreshToken, err := s.GenerateToken(ctx, user, true)
	if err != nil {
		return model.TokenDetail{}, fmt.Errorf("error generating new refresh token: %w", err)
	}

	if err := s.AuthRepo.DeleteRefreshTkn(ctx, hashedTkn); err != nil {
		log.Printf("failed to delete old refresh token: %v", err)
	}

	return model.TokenDetail{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

