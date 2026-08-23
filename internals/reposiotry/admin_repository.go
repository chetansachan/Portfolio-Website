package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	Pool *pgxpool.Pool
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{Pool: pool}
}

func (r *AdminRepo) CreateAdmin(ctx context.Context, email string, password string) (model.Admin, error) {
	query := `
	INSERT INTO admin (email, password) VALUES ($1, $2) RETURNING id, email, password, login_at
	`

	var admin model.Admin
	if err := r.Pool.QueryRow(ctx, query, email, password).Scan(&admin.AdminId, &admin.Email, &admin.Password, &admin.Login_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return admin, nil
}
