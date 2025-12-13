package repository

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

type ILecturerRepository interface {
	GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error)
	GetAllLecturers(ctx context.Context) ([]model.Lecturer, error)
	CreateLecturer(ctx context.Context, req model.CreateLecturerRequest) (*model.Lecturer, error)
	UpdateLecturer(ctx context.Context, id string, req model.UpdateLecturerRequest) (*model.Lecturer, error)
}

type LecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) ILecturerRepository {
	return &LecturerRepository{db: db}
}

func (r *LecturerRepository) GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.user_id = $1
	`

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

func (r *LecturerRepository) GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error) {
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.id = $1
	`

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

func (r *LecturerRepository) GetAllLecturers(ctx context.Context) ([]model.Lecturer, error) {
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, COALESCE(u.full_name, '') as full_name, l.created_at
		FROM lecturers l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *LecturerRepository) CreateLecturer(ctx context.Context, req model.CreateLecturerRequest) (*model.Lecturer, error) {
	query := `
		INSERT INTO lecturers (user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, lecturer_id, department, created_at
	`

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

func (r *LecturerRepository) UpdateLecturer(ctx context.Context, id string, req model.UpdateLecturerRequest) (*model.Lecturer, error) {
	query := `
		UPDATE lecturers
		SET lecturer_id = $1, department = $2
		WHERE id = $3
		RETURNING id, user_id, lecturer_id, department, created_at
	`

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
