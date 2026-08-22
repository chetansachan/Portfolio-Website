package newsletter

import (
	"context"
	"fmt"
	"portfolio-website/internals/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsletterWorker struct {
	Pool  *pgxpool.Pool
	Email *Client
}

func NewNewsletterWorker(pool *pgxpool.Pool, email *Client) *NewsletterWorker {
	return &NewsletterWorker{Pool: pool, Email: email}
}

func (w *NewsletterWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processPendingJobs(ctx)
		}
	}
}

func (w *NewsletterWorker) processPendingJobs(ctx context.Context) {
	var jobID uuid.UUID
	var blogID uuid.UUID
	var subscriberIDs []uuid.UUID

	err := w.Pool.QueryRow(ctx, `
		UPDATE newsletter_jobs 
		SET status = 'processing' 
		WHERE id = (
			SELECT id FROM newsletter_jobs 
			WHERE status = 'pending' 
			ORDER BY created_at 
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, blog_id, subscriber_ids
	`).Scan(&jobID, &blogID, &subscriberIDs)

	if err == pgx.ErrNoRows {
		return
	}
	if err != nil {
		fmt.Println("DB error fetching job:", err)
		return
	}

	var blogTitle, blogSlug, blogSummary string
	err = w.Pool.QueryRow(ctx, `
		SELECT blog_title, blog_slug, blog_des 
		FROM blogs 
		WHERE id = $1
	`, blogID).Scan(&blogTitle, &blogSlug, &blogSummary)
	if err != nil {
		fmt.Println("Failed to fetch blog:", err)
		w.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
		return
	}

	rows, err := w.Pool.Query(ctx, `
		SELECT email 
		FROM subscriber 
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

		err := w.Email.SendBlogNotification(ctx, emailAddr, model.Blog{
			Blog_Title: blogTitle,
			Blog_Slug: blogSlug,
			Blog_DES: blogSummary,
		})
		if err != nil {
			fmt.Printf("Failed to send to %s: %v\n", emailAddr, err)
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		return
	}

	_, err = w.Pool.Exec(ctx, `
		UPDATE newsletter_jobs 
		SET status = 'done', processed_at = NOW() 
		WHERE id = $1
	`, jobID)
	if err != nil {
		fmt.Println("Failed to mark job done:", err)
	}
}
