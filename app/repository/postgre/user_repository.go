package repository

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

type IUserRepository interface {
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetRoleName(ctx context.Context, roleID string) (string, error)
	GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error)
	SaveRefreshToken(ctx context.Context, userID string, token string, expiresAt string) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1
	`

	user := new(model.User)
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`

	user := new(model.User)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.username = $1 OR u.email = $1
	`

	user := new(model.User)
	err := r.db.QueryRowContext(ctx, query, usernameOrEmail).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt string
	CreatedAt string
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID string, token string, expiresAt string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`
	rt := new(RefreshToken)
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *UserRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1
		ORDER BY p.name
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *UserRepository) GetRoleName(ctx context.Context, roleID string) (string, error) {
	query := `SELECT name FROM roles WHERE id = $1`
	var roleName string
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&roleName)
	if err != nil {
		return "", err
	}
	return roleName, nil
}

func (r *UserRepository) GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
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

func (r *UserRepository) GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error) {
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
