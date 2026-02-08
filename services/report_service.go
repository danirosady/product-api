package services

import (
	"product-api/models"
	"product-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReport(startInclusive time.Time, endExclusive time.Time) (*models.ReportResponse, error) {
	return s.repo.GetReport(startInclusive, endExclusive)
}
