package repository

// #1 proses: import library yang diperlukan untuk database dan context
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

// #2 proses: definisikan interface untuk operasi database notification
type INotificationRepository interface {
	CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (*model.Notification, error)
	GetNotificationsByUserIDPaginated(ctx context.Context, userID string, page, limit int) ([]model.Notification, int, error)
	GetUnreadCountByUserID(ctx context.Context, userID string) (int, error)
	MarkAsRead(ctx context.Context, notificationID string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

// #3 proses: struct repository untuk operasi database notification
type NotificationRepository struct {
	db *sql.DB
}

// #4 proses: constructor untuk membuat instance NotificationRepository baru
func NewNotificationRepository(db *sql.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

// #5 proses: buat notifikasi baru di database
func (r *NotificationRepository) CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (*model.Notification, error) {
	// #5a proses: query untuk insert notifikasi baru dengan RETURNING
	query := `
		INSERT INTO notifications (user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, false, NOW(), NOW())
		RETURNING id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
	`

	// #5b proses: eksekusi query dan scan hasil, handle field yang bisa null
	notif := new(model.Notification)
	var achievementID sql.NullString
	var mongoAchievementID sql.NullString
	var readAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, req.UserID, req.Type, req.Title, req.Message, req.AchievementID, req.MongoAchievementID).Scan(
		&notif.ID, &notif.UserID, &notif.Type, &notif.Title, &notif.Message,
		&achievementID, &mongoAchievementID, &notif.IsRead, &readAt, &notif.CreatedAt, &notif.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// #5c proses: set achievement_id jika nilainya valid
	if achievementID.Valid {
		notif.AchievementID = &achievementID.String
	}

	// #5d proses: set mongo_achievement_id jika nilainya valid
	if mongoAchievementID.Valid {
		notif.MongoAchievementID = &mongoAchievementID.String
	}

	// #5e proses: set read_at jika nilainya valid
	if readAt.Valid {
		notif.ReadAt = &readAt.Time
	}

	return notif, nil
}

// #6 proses: ambil notifikasi user dengan pagination
func (r *NotificationRepository) GetNotificationsByUserIDPaginated(ctx context.Context, userID string, page, limit int) ([]model.Notification, int, error) {
	// #6a proses: hitung offset untuk pagination
	offset := (page - 1) * limit

	// #6b proses: query untuk hitung total notifikasi user
	countQuery := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// #6c proses: query untuk ambil notifikasi dengan limit dan offset
	query := `
		SELECT id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	// #6d proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// #6e proses: loop semua hasil dan masukkan ke slice notifications
	var notifications []model.Notification
	for rows.Next() {
		var notif model.Notification
		var achievementID sql.NullString
		var mongoAchievementID sql.NullString
		var readAt sql.NullTime

		err := rows.Scan(
			&notif.ID, &notif.UserID, &notif.Type, &notif.Title, &notif.Message,
			&achievementID, &mongoAchievementID, &notif.IsRead, &readAt, &notif.CreatedAt, &notif.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// #6f proses: set field yang bisa null jika nilainya valid
		if achievementID.Valid {
			notif.AchievementID = &achievementID.String
		}

		if mongoAchievementID.Valid {
			notif.MongoAchievementID = &mongoAchievementID.String
		}

		if readAt.Valid {
			notif.ReadAt = &readAt.Time
		}

		notifications = append(notifications, notif)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// #7 proses: hitung jumlah notifikasi yang belum dibaca oleh user
func (r *NotificationRepository) GetUnreadCountByUserID(ctx context.Context, userID string) (int, error) {
	// #7a proses: query untuk hitung notifikasi yang is_read = false
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// #8 proses: tandai satu notifikasi sebagai sudah dibaca
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	// #8a proses: query untuk update is_read jadi true dan set read_at
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	// #8b proses: eksekusi query dan cek apakah ada baris yang terupdate
	result, err := r.db.ExecContext(ctx, query, notificationID, userID)
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

// #9 proses: tandai semua notifikasi user sebagai sudah dibaca
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	// #9a proses: query untuk update semua notifikasi yang belum dibaca jadi sudah dibaca
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND is_read = false
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
