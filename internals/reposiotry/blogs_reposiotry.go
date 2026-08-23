package reposiotry

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"portfolio-website/internals/config"
	"portfolio-website/internals/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository struct {
	Pool  *pgxpool.Pool
	Email *config.Client
}

func NewBlogRepository(pool *pgxpool.Pool, email *config.Client) *BlogRepository {
	return &BlogRepository{Pool: pool, Email: email}
}

func (r *BlogRepository) CreateBlog(ctx context.Context, Blog_Title string, Blog_DES string, Blog_Slug string, Blog_Content string, Blog_Banner string) (model.Blog, error) {
	query := `
	INSERT INTO blogs (blog_title, blog_des, blog_slug, blog_content, blog_banner) VALUES ($1, $2, $3, $4, $5) RETURNING blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner, created_at
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, Blog_Title, Blog_DES, Blog_Slug, Blog_Content, Blog_Banner).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err != nil {
		log.Printf("Database error: %v", err)
	}

	return blog, nil
}

func (r *BlogRepository) GetAllBlogs(ctx context.Context, limit int, offset int) ([]model.Blog, int, error) {
	var totalBlogs int
	countQuery := `SELECT COUNT(*) FROM blogs`
	if err := r.Pool.QueryRow(ctx, countQuery).Scan(&totalBlogs); err != nil {
		return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
	}

	rows, err := r.Pool.Query(ctx, "SELECT * FROM blogs ORDER BY created_at LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return []model.Blog{}, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var blogs []model.Blog
	for rows.Next() {
		var blog model.Blog
		if err := rows.Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err == sql.ErrNoRows {
			log.Printf("Database error: %v", err)
			return []model.Blog{}, 0, err
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row loop error: %w", err)
	}

	return blogs, totalBlogs, nil
}

func (r *BlogRepository) GetBlogBySlug(ctx context.Context, slug string) (model.Blog, error) {
	query := `
		SELECT * FROM blogs WHERE blog_slug = $1
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, slug).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err == sql.ErrNoRows {
		log.Println("No user found with that slug")
		return model.Blog{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Blog{}, err
	}

	return blog, nil
}

func (r *BlogRepository) GetBlogById(ctx context.Context, id uuid.UUID) (model.Blog, error) {
	query := `
		SELECT * FROM blogs WHERE blog_id = $1
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, id).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err == sql.ErrNoRows {
		log.Println("No user found with that id")
		return model.Blog{}, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return model.Blog{}, err
	}

	return blog, nil
}

func (r *BlogRepository) DeleteBlogById(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM blogs WHERE blog_id = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	// Check if any row was actually deleted
	if cmdTag.RowsAffected() == 0 {
		return model.ErrBlogNotFound
	}

	return nil
}

func (r *BlogRepository) UpdateBlogById(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string) (model.Blog, error) {
	query := `
		UPDATE blogs SET blog_title = $2, blog_des = $3, blog_slug = $4, blog_content = $5, blog_banner = $6
		WHERE blog_id = $1
		RETURNING blog_id, blog_title, blog_des, blog_slug, blog_content, blog_banner, created_at
	`

	var blog model.Blog
	if err := r.Pool.QueryRow(ctx, query, id, title, des, slug, content, banner).Scan(&blog.Blog_ID, &blog.Blog_Title, &blog.Blog_DES, &blog.Blog_Slug, &blog.Blog_Content, &blog.Blog_Banner, &blog.Created_At); err != nil {
		log.Printf("Database error: %v", err)
		return model.Blog{}, err
	}

	return blog, nil
}

func (r *BlogRepository) ProcessPendingJobs(ctx context.Context) {
	var jobID uuid.UUID
	var blogID uuid.UUID
	var subscriberIDs []uuid.UUID

	err := r.Pool.QueryRow(ctx, `
		UPDATE newsletter_jobs 
		SET status = 'processing' 
		WHERE id = (
			SELECT id FROM newsletter_jobs 
			WHERE status = 'pending' 
			ORDER BY created_at 
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, blog_id, subscribers_group
	`).Scan(&jobID, &blogID, &subscriberIDs)

	if err == pgx.ErrNoRows {
		return
	}
	if err != nil {
		fmt.Println("DB error fetching job:", err)
		return
	}

	var blogTitle, blogSlug, blogSummary string
	err = r.Pool.QueryRow(ctx, `
		SELECT blog_title, blog_slug, blog_des 
		FROM blogs 
		WHERE blog_id = $1
	`, blogID).Scan(&blogTitle, &blogSlug, &blogSummary)
	if err != nil {
		fmt.Println("Failed to fetch blog:", err)
		r.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
		return
	}

	rows, err := r.Pool.Query(ctx, `
		SELECT email 
		FROM subscribers 
		WHERE is_active = true 
		AND id = ANY($1)
	`, subscriberIDs)
	if err != nil {
		fmt.Println("DB error fetching subscribers:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var emailAddr string
		if err := rows.Scan(&emailAddr); err != nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}

		err := r.Email.SendBlogNotification(ctx, emailAddr, model.Blog{
			Blog_Title: blogTitle,
			Blog_Slug:  blogSlug,
			Blog_DES:   blogSummary,
		})
		if err != nil {
			fmt.Printf("Failed to send to %s: %v\n", emailAddr, err)
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		return
	}

	_, err = r.Pool.Exec(ctx, `
		UPDATE newsletter_jobs 
		SET status = 'done', processed_at = NOW() 
		WHERE id = $1
	`, jobID)
	if err != nil {
		fmt.Println("Failed to mark job done:", err)
	}
}

func (r *BlogRepository) CreateNewsletterJobs(blogId uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := r.Pool.Query(ctx, `
		SELECT id FROM subscribers WHERE is_active = true
	`)
	if err != nil {
		fmt.Println("failed to fetch subscribers:", err)
		return
	}

	var allSubscriberIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			continue
		}
		allSubscriberIDs = append(allSubscriberIDs, id)
	}
	rows.Close()

	if len(allSubscriberIDs) == 0 {
		return
	}

	batchSize := 100
	for i := 0; i < len(allSubscriberIDs); i += batchSize {
		end := i + batchSize
		if end > len(allSubscriberIDs) {
			end = len(allSubscriberIDs)
		}
		batch := allSubscriberIDs[i:end]

		_, err := r.Pool.Exec(ctx, `
			INSERT INTO newsletter_jobs (blog_id, subscribers_group, status)
			VALUES ($1, $2, 'pending')
		`, blogId, batch)
		if err != nil {
			fmt.Printf("failed to insert job batch %d: %v\n", i/batchSize, err)
		}
	}

	fmt.Printf("created %d newsletter jobs for blog %s\n",
		(len(allSubscriberIDs)+batchSize-1)/batchSize, blogId)
}

func (r *BlogRepository) CreateSubscriber(ctx context.Context, email string) error {
	query := `
		INSERT INTO subscribers (email, is_active) 
		VALUES ($1, true) 
		ON CONFLICT (email) 
		DO UPDATE SET is_active = true, unsubscribed_at = NULL
	`

	if _, err := r.Pool.Exec(ctx, query, email); err != nil {
		log.Printf("Database error: %v", err)
		return err
	}

	return nil
}
