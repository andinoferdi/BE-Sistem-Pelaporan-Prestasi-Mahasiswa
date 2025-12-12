package service

import (
	"context"
	"errors"
)

type IReportService interface {
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	GetStudentReport(ctx context.Context, studentID string) (map[string]interface{}, error)
}

type ReportService struct {
}

func NewReportService() IReportService {
	return &ReportService{}
}

func (s *ReportService) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	return nil, errors.New("fitur ini belum diimplementasikan")
}

func (s *ReportService) GetStudentReport(ctx context.Context, studentID string) (map[string]interface{}, error) {
	return nil, errors.New("fitur ini belum diimplementasikan")
}
