package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	Pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{Pool: pool}
}

func (r *AuthRepo) RegisterUser(ctx context.Context, email string, password string) error {
	query := `
	INSERT INTO users (email, password) VALUES ($1, $2)
	`

	if _, err := r.Pool.Exec(ctx, query, email, password); err != nil {
		log.Printf("Database error: %v", err)
	}

	return nil
}

func (r *AuthRepo) Login(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`

	var admin model.User
	if err := r.Pool.QueryRow(ctx, query, email).Scan(&admin.UserId, &admin.Email, &admin.Password, &admin.Login_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return admin, nil
}