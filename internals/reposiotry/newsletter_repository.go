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
	var blogID *uuid.UUID
	var projectID *uuid.UUID
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
		RETURNING id, blog_id, project_id, subscribers_group
	`).Scan(&jobID, &blogID, &projectID, &subscriberIDs)

	if err == pgx.ErrNoRows {
		return
	}
	if err != nil {
		fmt.Println("DB error fetching job:", err)
		return
	}

	// Fetch subscriber emails
	rows, err := r.Pool.Query(ctx, `
		SELECT email 
		FROM subscribers 
		WHERE is_active = true 
		AND id = ANY($1)
	`, subscriberIDs)
	if err != nil {
		fmt.Println("DB error fetching subscribers:", err)
		r.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var emailAddr string
		if err := rows.Scan(&emailAddr); err == nil {
			emails = append(emails, emailAddr)
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("DB rows error fetching subscribers:", err)
		r.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
		return
	}

	if blogID != nil {
		var blogTitle, blogSlug, blogSummary string
		err = r.Pool.QueryRow(ctx, `
			SELECT blog_title, blog_slug, blog_des 
			FROM blogs 
			WHERE blog_id = $1
		`, *blogID).Scan(&blogTitle, &blogSlug, &blogSummary)
		if err != nil {
			fmt.Println("Failed to fetch blog:", err)
			r.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
			return
		}

		for _, emailAddr := range emails {
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
	} else if projectID != nil {
		var projectName, projectSlug, projectSummary string
		err = r.Pool.QueryRow(ctx, `
			SELECT project_name, project_slug, project_des 
			FROM projects 
			WHERE project_id = $1
		`, *projectID).Scan(&projectName, &projectSlug, &projectSummary)
		if err != nil {
			fmt.Println("Failed to fetch project:", err)
			r.Pool.Exec(ctx, `UPDATE newsletter_jobs SET status = 'failed' WHERE id = $1`, jobID)
			return
		}

		for _, emailAddr := range emails {
			select {
			case <-ctx.Done():
				return
			case <-time.After(100 * time.Millisecond):
			}

			err := r.Email.SendProjectNotification(ctx, emailAddr, model.Project{
				Project_Name: projectName,
				Project_Slug: projectSlug,
				Project_DES:  projectSummary,
			})
			if err != nil {
				fmt.Printf("Failed to send to %s: %v\n", emailAddr, err)
			}
		}
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

func (r *NewsletterRepo) CreateProjectNewsletterJobs(projectId uuid.UUID) {
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
			INSERT INTO newsletter_jobs (project_id, subscribers_group, status)
			VALUES ($1, $2, 'pending')
		`, projectId, batch)
		if err != nil {
			fmt.Printf("failed to insert project job batch %d: %v\n", i/batchSize, err)
		}
	}

	fmt.Printf("created %d newsletter jobs for project %s\n",
		(len(allSubscriberIDs)+batchSize-1)/batchSize, projectId)
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
