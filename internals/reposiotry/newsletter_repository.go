package reposiotry

import (
	"context"
	"fmt"
	"log"
	"portfolio-website/internals/config"
	"portfolio-website/internals/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsletterRepo struct {
	Pool  *pgxpool.Pool
	Email *config.Client
}

func NewNewsletterRepo(pool *pgxpool.Pool, email *config.Client) *NewsletterRepo {
	return &NewsletterRepo{Pool: pool, Email: email}
}

func (r *NewsletterRepo) ProcessPendingJobs(ctx context.Context) {
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

func (r *NewsletterRepo) CreateNewsletterJobs(blogId uuid.UUID) {
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

func (r *NewsletterRepo) CreateSubscriber(ctx context.Context, email string) error {
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
