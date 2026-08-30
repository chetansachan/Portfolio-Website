package service

import (
	"context"
	"fmt"
	"log"
	"errors"
	"portfolio-website/internals/model"
	"portfolio-website/internals/reposiotry"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CaseStudyService struct {
	CaseStudyRepository *reposiotry.CaseStudyRepository
	NewsletterRepo      *reposiotry.NewsletterRepo
}

func NewCaseStudyService(caseStudyRepo *reposiotry.CaseStudyRepository, newsletterRepo *reposiotry.NewsletterRepo) *CaseStudyService {
	return &CaseStudyService{CaseStudyRepository: caseStudyRepo, NewsletterRepo: newsletterRepo}
}

func (s *CaseStudyService) CreateCaseStudy_Service(ctx context.Context, title, des, slug, content, banner string, featured bool) (model.Case_Study, error) {
	if title == "" || des == "" || slug == "" || content == "" || banner == "" {
		log.Println("[WARN] validation failed")
		return model.Case_Study{}, model.ErrValidation
	}

	c, err := s.CaseStudyRepository.GetCaseStudyBySlug(ctx, slug)
	if err == nil {
		log.Println("Case study with this slug is already exist")
		return c, model.ErrCaseStudyAlreadyExist
	}

	casestudy, err := s.CaseStudyRepository.CreateCaseStudy(ctx, title, des, slug, content, banner, featured)
	if err != nil {
		log.Printf("[ERROR] database insert failed: %v", err)
		return model.Case_Study{}, fmt.Errorf("failed to save case study: %w", err)
	}

	go s.NewsletterRepo.CreateCaseStudyNewsletterJobs(casestudy.Case_Study_ID)

	return casestudy, nil
}

func (s *CaseStudyService) GetAllCaseStudies_Service(ctx context.Context, limit, offset int, featured *bool) ([]model.Case_Study, int, error) {
	case_studies, totalcasestudies, err := s.CaseStudyRepository.GetAllCaseStudies(ctx, limit, offset, featured)
	if err != nil {
		log.Printf("[ERROR] failed to get case studies from database: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Case_Study{}, 0, model.ErrCaseStudyNotFound
		}
		return []model.Case_Study{}, 0, err
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
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Case_Study{}, model.ErrCaseStudyNotFound
		}
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (s *CaseStudyService) GetCaseStudyById_Service(ctx context.Context, id uuid.UUID) (model.Case_Study, error) {
	if id == uuid.Nil {
		log.Println("[WARN] validation failed")
		return model.Case_Study{}, model.ErrIdValidation
	}

	case_study, err := s.CaseStudyRepository.GetCaseStudyById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] failed to get case study from database: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Case_Study{}, model.ErrCaseStudyNotFound
		}
		return model.Case_Study{}, err
	}

	return case_study, nil
}

func (s *CaseStudyService) DeleteCaseStudyById_Service(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return model.ErrIdValidation
	}

	if err := s.CaseStudyRepository.DeleteCaseStudyById(ctx, id); err != nil {
		log.Printf("[ERROR] failed to delete case study: %v", err)
		return err
	}

	return nil
}

func (s *CaseStudyService) UpdateCaseStudyById_Service(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string, featured bool) (model.Case_Study, error) {
	if id == uuid.Nil {
		log.Println("[WARN] validation failed")
		return model.Case_Study{}, model.ErrIdValidation
	}

	c, err := s.CaseStudyRepository.GetCaseStudyById(ctx, id)
	if err != nil {
		log.Printf("[ERROR] case study not found: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Case_Study{}, model.ErrCaseStudyNotFound
		}
		return model.Case_Study{}, err
	}

	if title == "" {
		title = c.Case_Study_Title
	}
	if des == "" {
		des = c.Case_Study_DES
	}
	if slug == "" {
		slug = c.Case_Study_Slug
	}
	if content == "" {
		content = c.Case_Study_Content
	}
	if banner == "" {
		banner = c.Case_Study_Banner
	}

	case_study, err := s.CaseStudyRepository.UpdateCaseStudyById(ctx, id, title, des, slug, content, banner, featured)
	if err != nil {
		log.Printf("[ERROR] failed to update case study: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Case_Study{}, model.ErrCaseStudyNotFound
		}
		return model.Case_Study{}, err
	}

	return case_study, nil
}
