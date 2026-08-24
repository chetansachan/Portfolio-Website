package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	Pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{Pool: pool}
}

func (r *AuthRepo) RegisterUser(ctx context.Context, name string, email string, password string) error {
	query := `
	INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
	`

	if _, err := r.Pool.Exec(ctx, query, name, email, password); err != nil {
		log.Printf("Database error: %v", err)
	}

	return nil
}

func (r *AuthRepo) Login(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`

	var user model.User
	if err := r.Pool.QueryRow(ctx, query, email).Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Is_Admin, &user.Login_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return user, nil
}

func (r *AuthRepo) SaveRefreshTkn(ctx context.Context, token string, userID uuid.UUID, userName string, ipAddress string, expiredAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (token, user_id, user_name, ip_address, expires_at)
        VALUES ($1, $2, $3, $4, $5);
	`

	if _, err := r.Pool.Exec(ctx, query, token, userID, userName, ipAddress, expiredAt); err != nil {
		log.Println("Failed to save token")
		return err
	}

	return nil
}
