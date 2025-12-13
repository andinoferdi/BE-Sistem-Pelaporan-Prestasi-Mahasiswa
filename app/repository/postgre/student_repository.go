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
	GetAllStudents(ctx context.Context) ([]model.Student, error)
	CreateStudent(ctx context.Context, req model.CreateStudentRequest) (*model.Student, error)
	UpdateStudent(ctx context.Context, id string, req model.UpdateStudentRequest) (*model.Student, error)
	UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error
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
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1
	`

	student := new(model.Student)
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.FullName, &student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`

	student := new(model.Student)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.FullName, &student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
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
			&student.FullName, &student.CreatedAt,
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

func (r *StudentRepository) GetAllStudents(ctx context.Context) ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
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
			&student.FullName, &student.CreatedAt,
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

func (r *StudentRepository) CreateStudent(ctx context.Context, req model.CreateStudentRequest) (*model.Student, error) {
	query := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	student := new(model.Student)
	var advisorID sql.NullString
	err := r.db.QueryRowContext(ctx, query,
		req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID,
	).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &advisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

func (r *StudentRepository) CreateStudentWithTx(ctx context.Context, tx *sql.Tx, req model.CreateStudentRequest) (*model.Student, error) {
	query := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	student := new(model.Student)
	var advisorID sql.NullString
	err := tx.QueryRowContext(ctx, query,
		req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID,
	).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &advisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

func (r *StudentRepository) UpdateStudent(ctx context.Context, id string, req model.UpdateStudentRequest) (*model.Student, error) {
	query := `
		UPDATE students
		SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
		WHERE id = $5
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	student := new(model.Student)
	var advisorID sql.NullString
	err := r.db.QueryRowContext(ctx, query,
		req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID, id,
	).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &advisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

func (r *StudentRepository) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`

	var advisorIDValue interface{}
	if advisorID == "" {
		advisorIDValue = nil
	} else {
		advisorIDValue = advisorID
	}

	result, err := r.db.ExecContext(ctx, query, advisorIDValue, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
