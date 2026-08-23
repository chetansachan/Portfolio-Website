package workers

import (
	"context"
	"portfolio-website/internals/config"
	"portfolio-website/internals/reposiotry"
	"time"
)

type NewsletterWorker struct {
	Email          *config.Client
	NewsletterRepo *reposiotry.NewsletterRepo
}

func NewNewsletterWorker(email *config.Client, newsletterRepo *reposiotry.NewsletterRepo) *NewsletterWorker {
	return &NewsletterWorker{Email: email, NewsletterRepo: newsletterRepo}
}

func (w *NewsletterWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			w.NewsletterRepo.ProcessPendingJobs(ctx)
		}
	}
}
