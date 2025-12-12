package service

import (
	"context"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
)

type IStudentService interface {
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error)
	GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error)
}

type StudentService struct {
	studentRepo repositorypostgre.IStudentRepository
}

func NewStudentService(studentRepo repositorypostgre.IStudentRepository) IStudentService {
	return &StudentService{studentRepo: studentRepo}
}

func (s *StudentService) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	return s.studentRepo.GetStudentIDByUserID(ctx, userID)
}

func (s *StudentService) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

func (s *StudentService) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if advisorID == "" {
		return nil, errors.New("advisor ID wajib diisi")
	}
	return s.studentRepo.GetStudentsByAdvisorID(ctx, advisorID)
}
