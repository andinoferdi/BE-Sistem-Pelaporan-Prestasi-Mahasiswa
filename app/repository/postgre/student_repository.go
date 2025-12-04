package repository

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func GetStudentIDByUserID(db *sql.DB, userID string) (string, error) {
	query := `
		SELECT s.id
		FROM students s
		WHERE s.user_id = $1
	`

	var studentID string
	err := db.QueryRow(query, userID).Scan(&studentID)
	if err != nil {
		return "", err
	}

	return studentID, nil
}

func GetStudentByUserID(db *sql.DB, userID string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, s.created_at
		FROM students s
		WHERE s.user_id = $1
	`

	student := new(model.Student)
	err := db.QueryRow(query, userID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

