package repository

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

type INotificationRepository interface {
	CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (*model.Notification, error)
	GetNotificationsByUserIDPaginated(ctx context.Context, userID string, page, limit int) ([]model.Notification, int, error)
	GetUnreadCountByUserID(ctx context.Context, userID string) (int, error)
	MarkAsRead(ctx context.Context, notificationID string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, req model.CreateNotificationRequest) (*model.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, false, NOW(), NOW())
		RETURNING id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
	`

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

	if achievementID.Valid {
		notif.AchievementID = &achievementID.String
	}

	if mongoAchievementID.Valid {
		notif.MongoAchievementID = &mongoAchievementID.String
	}

	if readAt.Valid {
		notif.ReadAt = &readAt.Time
	}

	return notif, nil
}

func (r *NotificationRepository) GetNotificationsByUserIDPaginated(ctx context.Context, userID string, page, limit int) ([]model.Notification, int, error) {
	offset := (page - 1) * limit

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

	query := `
		SELECT id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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

func (r *NotificationRepository) GetUnreadCountByUserID(ctx context.Context, userID string) (int, error) {
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

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

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

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND is_read = false
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
