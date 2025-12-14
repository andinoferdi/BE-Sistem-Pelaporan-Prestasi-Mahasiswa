package repository

// #1 proses: import library yang diperlukan untuk database dan context
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

// #2 proses: definisikan interface untuk operasi database lecturer
type ILecturerRepository interface {
	GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error)
	GetAllLecturers(ctx context.Context) ([]model.Lecturer, error)
	CreateLecturer(ctx context.Context, req model.CreateLecturerRequest) (*model.Lecturer, error)
	UpdateLecturer(ctx context.Context, id string, req model.UpdateLecturerRequest) (*model.Lecturer, error)
}

// #3 proses: struct repository untuk operasi database lecturer
type LecturerRepository struct {
	db *sql.DB
}

// #4 proses: constructor untuk membuat instance LecturerRepository baru
func NewLecturerRepository(db *sql.DB) ILecturerRepository {
	return &LecturerRepository{db: db}
}

// #5 proses: ambil data lecturer berdasarkan user ID
func (r *LecturerRepository) GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	// #5a proses: query untuk ambil data lecturer dari database
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.user_id = $1
	`

	// #5b proses: eksekusi query dan scan hasil ke struct lecturer
	lecturer := new(model.Lecturer)
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

// #6 proses: ambil data lecturer berdasarkan lecturer ID
func (r *LecturerRepository) GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error) {
	// #6a proses: query untuk ambil data lecturer berdasarkan ID
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.id = $1
	`

	// #6b proses: eksekusi query dan scan hasil ke struct lecturer
	lecturer := new(model.Lecturer)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

// #7 proses: ambil semua lecturer yang ada di database
func (r *LecturerRepository) GetAllLecturers(ctx context.Context) ([]model.Lecturer, error) {
	// #7a proses: query untuk ambil semua lecturer dengan join ke tabel users untuk ambil full_name
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, COALESCE(u.full_name, '') as full_name, l.created_at
		FROM lecturers l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
	`

	// #7b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #7c proses: loop semua hasil dan masukkan ke slice lecturers
	var lecturers []model.Lecturer
	for rows.Next() {
		var lecturer model.Lecturer
		err := rows.Scan(
			&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
			&lecturer.Department, &lecturer.FullName, &lecturer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		lecturers = append(lecturers, lecturer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lecturers, nil
}

// #8 proses: buat lecturer baru di database
func (r *LecturerRepository) CreateLecturer(ctx context.Context, req model.CreateLecturerRequest) (*model.Lecturer, error) {
	// #8a proses: query untuk insert lecturer baru dengan RETURNING
	query := `
		INSERT INTO lecturers (user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, lecturer_id, department, created_at
	`

	// #8b proses: eksekusi query dan scan hasil ke struct lecturer
	lecturer := new(model.Lecturer)
	err := r.db.QueryRowContext(ctx, query,
		req.UserID, req.LecturerID, req.Department,
	).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

// #9 proses: buat lecturer baru dalam transaction, dipakai untuk operasi multi-step
func (r *LecturerRepository) CreateLecturerWithTx(ctx context.Context, tx *sql.Tx, req model.CreateLecturerRequest) (*model.Lecturer, error) {
	// #9a proses: query untuk insert lecturer baru dalam transaction
	query := `
		INSERT INTO lecturers (user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, lecturer_id, department, created_at
	`

	// #9b proses: eksekusi query menggunakan transaction dan scan hasil
	lecturer := new(model.Lecturer)
	err := tx.QueryRowContext(ctx, query,
		req.UserID, req.LecturerID, req.Department,
	).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

// #10 proses: update data lecturer yang sudah ada
func (r *LecturerRepository) UpdateLecturer(ctx context.Context, id string, req model.UpdateLecturerRequest) (*model.Lecturer, error) {
	// #10a proses: query untuk update lecturer dengan RETURNING
	query := `
		UPDATE lecturers
		SET lecturer_id = $1, department = $2
		WHERE id = $3
		RETURNING id, user_id, lecturer_id, department, created_at
	`

	// #10b proses: eksekusi query dan scan hasil ke struct lecturer
	lecturer := new(model.Lecturer)
	err := r.db.QueryRowContext(ctx, query,
		req.LecturerID, req.Department, id,
	).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}
