package service

import (
	"context"
	"fmt"
	"log"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"
)

type CaseStudyService struct {
	CaseStudyRepository *reposiotry.CaseStudyRepository
	NewsletterRepo      *reposiotry.NewsletterRepo
}

func NewCaseStudyService(caseStudyRepo *reposiotry.CaseStudyRepository, newsletterRepo *reposiotry.NewsletterRepo) *CaseStudyService {
	return &CaseStudyService{CaseStudyRepository: caseStudyRepo, NewsletterRepo: newsletterRepo}
}

func (s *CaseStudyService) CreateCaseStudy_Service(ctx context.Context, title, des, slug, content, banner string) (model.Case_Study, error) {
	if title == "" || des == "" || slug == "" || content == "" || banner == "" {
		log.Println("[WARN] validation failed")
		return model.Case_Study{}, model.ErrValidation
	}

	c, err := s.CaseStudyRepository.GetCaseStudyBySlug(ctx, slug)
	if err == nil {
		log.Println("Case study with this slug is already exist")
		return c, model.ErrCaseStudyAlreadyExist
	}

	casestudy, err := s.CaseStudyRepository.CreateCaseStudy(ctx, title, des, slug, content, banner)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Case_Study{}, fmt.Errorf("failed to save case study: %w", err)
	}

	go s.NewsletterRepo.CreateNewsletterJobs(casestudy.Case_Study_ID)

	return casestudy, nil
}

func (s *CaseStudyService) GetAllCaseStudies_Service(ctx context.Context, limit, offset int) ([]model.Case_Study, int, error) {
	case_studies, totalcasestudies, err := s.CaseStudyRepository.GetAllCaseStudies(ctx, limit, offset)
	if err != nil {
		log.Printf("[ERROR] failed to get case studies from database: %v", err)
		return []model.Case_Study{}, 0, model.ErrCaseStudyNotFound
	}

	return case_studies, totalcasestudies, nil
}

func (s *CaseStudyService) GetCaseStudyBySlug_Service(ctx context.Context, slug string) (model.Case_Study, error) {
	if slug == "" {
		log.Println("[WARN] validation failed")
		return model.Case_Study{}, model.ErrSlugValidation
	}

	case_study, err := s.CaseStudyRepository.GetCaseStudyBySlug(ctx, slug)
	if err != nil {
		log.Printf("[ERROR] failed to get case study from database: %v", err)
		return model.Case_Study{}, model.ErrBlogNotFound
	}

	return case_study, nil
}
