package repository

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

type IStudentRepository interface {
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetStudentByUserID(ctx context.Context, userID string) (*model.Student, error)
	GetStudentByID(ctx context.Context, id string) (*model.Student, error)
	GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]model.Student, error)
}

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) IStudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	query := `
		SELECT s.id
		FROM students s
		WHERE s.user_id = $1
	`

	var studentID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&studentID)
	if err != nil {
		return "", err
	}

	return studentID, nil
}

func (r *StudentRepository) GetStudentByUserID(ctx context.Context, userID string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, s.created_at
		FROM students s
		WHERE s.user_id = $1
	`

	student := new(model.Student)
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, s.created_at
		FROM students s
		WHERE s.id = $1
	`

	student := new(model.Student)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, s.created_at
		FROM students s
		WHERE s.advisor_id = $1
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(
			&student.ID, &student.UserID, &student.StudentID,
			&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
			&student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
