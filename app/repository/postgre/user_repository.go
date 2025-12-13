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
	GetAllUsers(ctx context.Context) ([]model.User, error)
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	UpdateUser(ctx context.Context, id string, user model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserRole(ctx context.Context, id string, roleID string) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetRoleName(ctx context.Context, roleID string) (string, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error)
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

func (r *UserRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
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

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FullName, &user.RoleID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	newUser := new(model.User)
	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive,
	).Scan(
		&newUser.ID, &newUser.Username, &newUser.Email, &newUser.PasswordHash,
		&newUser.FullName, &newUser.RoleID, &newUser.IsActive,
		&newUser.CreatedAt, &newUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *UserRepository) CreateUserWithTx(ctx context.Context, tx *sql.Tx, user model.User) (*model.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	newUser := new(model.User)
	err := tx.QueryRowContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive,
	).Scan(
		&newUser.ID, &newUser.Username, &newUser.Email, &newUser.PasswordHash,
		&newUser.FullName, &newUser.RoleID, &newUser.IsActive,
		&newUser.CreatedAt, &newUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, user model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, full_name = $3, role_id = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	updatedUser := new(model.User)
	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.FullName, user.RoleID, user.IsActive, id,
	).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.PasswordHash,
		&updatedUser.FullName, &updatedUser.RoleID, &updatedUser.IsActive,
		&updatedUser.CreatedAt, &updatedUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *UserRepository) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	query := `
		UPDATE users
		SET role_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, roleID, id)
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
