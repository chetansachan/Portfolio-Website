package service

import (
	"context"
	"fmt"
	"log"
	"portfolio-website/internals/reposiotry"
	"regexp"
)

type NewsletterService struct {
	BlogRepository *reposiotry.NewsletterRepo
}

func NewNewsletterService(blogRepo *reposiotry.NewsletterRepo) *NewsletterService {
	return &NewsletterService{BlogRepository: blogRepo}
}

func (s *NewsletterService) CreateSubscriber_Service(ctx context.Context, email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if re.MatchString(email) {
		fmt.Printf("%s is valid\n", email)
	} else {
		return fmt.Errorf("%s is invalid\n", email)
	}

	if err := s.BlogRepository.CreateSubscriber(ctx, email); err != nil {
		log.Printf("[ERROR] Failed to Subscribe: %v", err)
		return err
	}

	return nil
}
