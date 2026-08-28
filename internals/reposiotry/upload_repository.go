package reposiotry

import (
	"context"
	"log"
	"portfolio-website/internals/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadRepo struct {
	Pool *pgxpool.Pool
}

func NewUploadRepo(pool *pgxpool.Pool) *UploadRepo {
	return &UploadRepo{Pool: pool}
}

func (r *UploadRepo) SaveFiles(ctx context.Context, name string, ext string, mimeType string) error {
	query := `
		INSERT INTO files (file_name, file_ext, file_mimetype)
		VALUES ($1, $2, $3)
	`

	if _, err := r.Pool.Exec(ctx, query, name, ext, mimeType); err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	return nil
}

func (r *UploadRepo) DeleteFileByName(ctx context.Context, name string) error {
	query := `
		DELETE FROM files WHERE file_name = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, name)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrFileNotFound
	}

	return nil
}
