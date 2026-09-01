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
		return err
	}

	return nil
}

func (r *AuthRepo) Login(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT * FROM users WHERE email = $1
	`

	var user model.User
	err := r.Pool.QueryRow(ctx, query, email).Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Is_Admin, &user.Login_At)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthRepo) SaveRefreshTkn(ctx context.Context, token string, userID uuid.UUID, userName string, expiredAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (token, user_id, user_name, expires_at)
        VALUES ($1, $2, $3, $4);
	`

	if _, err := r.Pool.Exec(ctx, query, token, userID, userName, expiredAt); err != nil {
		log.Println("Failed to save token")
		return err
	}

	return nil
}

func (r *AuthRepo) GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	query := `
		SELECT id, token, user_id, user_name, is_revoked, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	var rt model.RefreshToken
	err := r.Pool.QueryRow(ctx, query, token).Scan(
		&rt.ID, &rt.Token, &rt.UserID, &rt.UserName, &rt.IsRevoked, &rt.ExpiresAt, &rt.CreatedAt,
	)
	if err != nil {
		return model.RefreshToken{}, err
	}
	return rt, nil
}

func (r *AuthRepo) DeleteRefreshTkn(ctx context.Context, token string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`
	_, err := r.Pool.Exec(ctx, query, token)
	return err
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `
		SELECT id, name, email, password, is_admin, login_at
		FROM users
		WHERE id = $1
	`
	var user model.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.UserId, &user.Name, &user.Email, &user.Password, &user.Is_Admin, &user.Login_At,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *AuthRepo) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}

	return nil
}

func (r *AuthRepo) Logout(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE token = $1
	`
	_, err := r.Pool.Exec(ctx, query, token)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	return nil
}

func (r *AuthRepo) MakeAdmin(ctx context.Context, email string) error {
	query := `
		UPDATE users SET is_admin = true
		WHERE email = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, email)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}

	return nil
}
