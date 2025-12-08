package repository

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func CreateNotification(db *sql.DB, req model.CreateNotificationRequest) (*model.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, false, NOW(), NOW())
		RETURNING id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
	`

	notif := new(model.Notification)
	var achievementID sql.NullString
	var mongoAchievementID sql.NullString
	var readAt sql.NullTime

	err := db.QueryRow(query, req.UserID, req.Type, req.Title, req.Message, req.AchievementID, req.MongoAchievementID).Scan(
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

func GetNotificationsByUserIDPaginated(db *sql.DB, userID string, page, limit int) ([]model.Notification, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1
	`
	var total int
	err := db.QueryRow(countQuery, userID).Scan(&total)
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

	rows, err := db.Query(query, userID, limit, offset)
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

func GetUnreadCountByUserID(db *sql.DB, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false
	`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func MarkAsRead(db *sql.DB, notificationID string, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	result, err := db.Exec(query, notificationID, userID)
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

func MarkAllAsRead(db *sql.DB, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND is_read = false
	`

	_, err := db.Exec(query, userID)
	return err
}

