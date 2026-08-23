package workers

import (
	"context"
	"portfolio-website/internals/config"
	"portfolio-website/internals/reposiotry"
	"time"
)

type NewsletterWorker struct {
	Email *config.Client
	BlogRepo *reposiotry.BlogRepository
}

func NewNewsletterWorker(email *config.Client, blogRepo *reposiotry.BlogRepository) *NewsletterWorker {
	return &NewsletterWorker{Email: email, BlogRepo: blogRepo}
}

func (w *NewsletterWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			w.BlogRepo.ProcessPendingJobs(ctx)
		}
	}
}
