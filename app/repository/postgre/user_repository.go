package repository

// #1 proses: import library yang diperlukan untuk database dan context
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

// #2 proses: definisikan interface untuk operasi database user
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

// #3 proses: struct repository untuk operasi database user
type UserRepository struct {
	db *sql.DB
}

// #4 proses: constructor untuk membuat instance UserRepository baru
func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

// #5 proses: cari user berdasarkan email
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// #5a proses: query untuk ambil data user dari database
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1
	`

	// #5b proses: eksekusi query dan scan hasil ke struct user
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

// #6 proses: cari user berdasarkan ID
func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	// #6a proses: query untuk ambil data user berdasarkan ID
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`

	// #6b proses: eksekusi query dan scan hasil ke struct user
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

// #7 proses: cari user berdasarkan username atau email, dipakai untuk login
func (r *UserRepository) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	// #7a proses: query untuk cari user dengan username atau email
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.username = $1 OR u.email = $1
	`

	// #7b proses: eksekusi query dan scan hasil ke struct user
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

// #8 proses: ambil semua permission yang dimiliki user berdasarkan role
func (r *UserRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// #8a proses: query untuk ambil permission dari role user
	query := `
		SELECT p.name
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1
		ORDER BY p.name
	`
	// #8b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #8c proses: loop semua hasil dan masukkan ke slice permissions
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

// #9 proses: ambil nama role berdasarkan role ID
func (r *UserRepository) GetRoleName(ctx context.Context, roleID string) (string, error) {
	// #9a proses: query untuk ambil nama role dari database
	query := `SELECT name FROM roles WHERE id = $1`
	var roleName string
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&roleName)
	if err != nil {
		return "", err
	}
	return roleName, nil
}

// #10 proses: ambil semua role yang ada di database
func (r *UserRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	// #10a proses: query untuk ambil semua role, diurutkan berdasarkan nama
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY name
	`

	// #10b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #10c proses: loop semua hasil dan masukkan ke slice roles
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

// #11 proses: ambil data lecturer berdasarkan user ID
func (r *UserRepository) GetLecturerByUserID(ctx context.Context, userID string) (*model.Lecturer, error) {
	// #11a proses: query untuk ambil data lecturer dari database
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.user_id = $1
	`

	// #11b proses: eksekusi query dan scan hasil ke struct lecturer
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

// #12 proses: ambil data lecturer berdasarkan lecturer ID
func (r *UserRepository) GetLecturerByID(ctx context.Context, id string) (*model.Lecturer, error) {
	// #12a proses: query untuk ambil data lecturer berdasarkan ID
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.id = $1
	`

	// #12b proses: eksekusi query dan scan hasil ke struct lecturer
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

// #13 proses: ambil semua user yang ada di database
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	// #13a proses: query untuk ambil semua user, diurutkan dari yang terbaru
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		ORDER BY u.created_at DESC
	`

	// #13b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #13c proses: loop semua hasil dan masukkan ke slice users
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

// #14 proses: buat user baru di database
func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	// #14a proses: query untuk insert user baru dengan RETURNING untuk ambil data yang baru dibuat
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	// #14b proses: eksekusi query dan scan hasil ke struct newUser
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

// #15 proses: buat user baru dalam transaction, dipakai untuk operasi multi-step
func (r *UserRepository) CreateUserWithTx(ctx context.Context, tx *sql.Tx, user model.User) (*model.User, error) {
	// #15a proses: query untuk insert user baru dalam transaction
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	// #15b proses: eksekusi query menggunakan transaction dan scan hasil
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

// #16 proses: update data user yang sudah ada
func (r *UserRepository) UpdateUser(ctx context.Context, id string, user model.User) (*model.User, error) {
	// #16a proses: query untuk update user dengan RETURNING untuk ambil data yang sudah diupdate
	query := `
		UPDATE users
		SET username = $1, email = $2, full_name = $3, role_id = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
	`

	// #16b proses: eksekusi query dan scan hasil ke struct updatedUser
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

// #17 proses: hapus user dari database
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	// #17a proses: query untuk delete user berdasarkan ID
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// #17b proses: cek apakah ada baris yang terhapus, jika tidak return error
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// #18 proses: update role user, ganti role_id user tertentu
func (r *UserRepository) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	// #18a proses: query untuk update role_id user
	query := `
		UPDATE users
		SET role_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	// #18b proses: eksekusi query dan cek apakah ada baris yang terupdate
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
