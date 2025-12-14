package repository

// #1 proses: import library yang diperlukan untuk database dan context
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

// #2 proses: definisikan interface untuk operasi database student
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

// #3 proses: struct repository untuk operasi database student
type StudentRepository struct {
	db *sql.DB
}

// #4 proses: constructor untuk membuat instance StudentRepository baru
func NewStudentRepository(db *sql.DB) IStudentRepository {
	return &StudentRepository{db: db}
}

// #5 proses: ambil student ID berdasarkan user ID
func (r *StudentRepository) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	// #5a proses: query untuk ambil student ID dari database
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

// #6 proses: ambil data student lengkap berdasarkan user ID
func (r *StudentRepository) GetStudentByUserID(ctx context.Context, userID string) (*model.Student, error) {
	// #6a proses: query untuk ambil data student dengan join ke tabel users untuk ambil full_name
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.user_id = $1
	`

	// #6b proses: eksekusi query dan scan hasil ke struct student
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

// #7 proses: ambil data student lengkap berdasarkan student ID
func (r *StudentRepository) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	// #7a proses: query untuk ambil data student dengan join ke tabel users
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`

	// #7b proses: eksekusi query dan scan hasil ke struct student
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

// #8 proses: ambil semua student yang dibimbing oleh dosen wali tertentu
func (r *StudentRepository) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]model.Student, error) {
	// #8a proses: query untuk ambil semua student berdasarkan advisor_id
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.advisor_id = $1
		ORDER BY s.created_at DESC
	`

	// #8b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #8c proses: loop semua hasil dan masukkan ke slice students
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

// #9 proses: ambil semua student yang ada di database
func (r *StudentRepository) GetAllStudents(ctx context.Context) ([]model.Student, error) {
	// #9a proses: query untuk ambil semua student dengan join ke tabel users
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE(u.full_name, '') as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
	`

	// #9b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #9c proses: loop semua hasil dan masukkan ke slice students
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

// #10 proses: buat student baru di database
func (r *StudentRepository) CreateStudent(ctx context.Context, req model.CreateStudentRequest) (*model.Student, error) {
	// #10a proses: query untuk insert student baru dengan RETURNING
	query := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	// #10b proses: eksekusi query dan scan hasil, handle advisor_id yang bisa null
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

	// #10c proses: set advisor_id jika nilainya valid
	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

// #11 proses: buat student baru dalam transaction, dipakai untuk operasi multi-step
func (r *StudentRepository) CreateStudentWithTx(ctx context.Context, tx *sql.Tx, req model.CreateStudentRequest) (*model.Student, error) {
	// #11a proses: query untuk insert student baru dalam transaction
	query := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	// #11b proses: eksekusi query menggunakan transaction dan scan hasil
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

	// #11c proses: set advisor_id jika nilainya valid
	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

// #12 proses: update data student yang sudah ada
func (r *StudentRepository) UpdateStudent(ctx context.Context, id string, req model.UpdateStudentRequest) (*model.Student, error) {
	// #12a proses: query untuk update student dengan RETURNING
	query := `
		UPDATE students
		SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
		WHERE id = $5
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	// #12b proses: eksekusi query dan scan hasil, handle advisor_id yang bisa null
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

	// #12c proses: set advisor_id jika nilainya valid
	if advisorID.Valid {
		student.AdvisorID = advisorID.String
	}

	return student, nil
}

// #13 proses: update advisor untuk student tertentu
func (r *StudentRepository) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	// #13a proses: query untuk update advisor_id student
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`

	// #13b proses: handle advisor_id yang bisa kosong, set jadi nil jika kosong
	var advisorIDValue interface{}
	if advisorID == "" {
		advisorIDValue = nil
	} else {
		advisorIDValue = advisorID
	}

	// #13c proses: eksekusi query dan cek apakah ada baris yang terupdate
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
