package service

import (
	"database/sql"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNotificationsService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	page := helper.GetQueryInt(c, "page", 1)
	limit := helper.GetQueryInt(c, "limit", 10)
	page, limit = helper.ValidatePagination(page, limit)

	notifications, total, err := repositorypostgre.GetNotificationsByUserIDPaginated(postgresDB, userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil notifications. Detail: " + err.Error(),
			},
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	response := modelpostgre.GetNotificationsResponse{
		Status: "success",
		Data:   notifications,
	}
	response.Pagination.Page = page
	response.Pagination.Limit = limit
	response.Pagination.Total = total
	response.Pagination.TotalPages = totalPages

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetUnreadCountService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	count, err := repositorypostgre.GetUnreadCountByUserID(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil unread count. Detail: " + err.Error(),
			},
		})
	}

	response := modelpostgre.GetUnreadCountResponse{
		Status: "success",
	}
	response.Data.Count = count

	return c.Status(fiber.StatusOK).JSON(response)
}

func MarkAsReadService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID notification wajib diisi.",
			},
		})
	}

	err := repositorypostgre.MarkAsRead(postgresDB, notificationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Notifikasi tidak ditemukan atau bukan milik Anda.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menandai notification sebagai read. Detail: " + err.Error(),
			},
		})
	}

	var notification modelpostgre.Notification
	query := `SELECT id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at FROM notifications WHERE id = $1 AND user_id = $2`
	var achievementID sql.NullString
	var mongoAchievementID sql.NullString
	var readAt sql.NullTime
	err = postgresDB.QueryRow(query, notificationID, userID).Scan(
		&notification.ID, &notification.UserID, &notification.Type, &notification.Title, &notification.Message,
		&achievementID, &mongoAchievementID, &notification.IsRead, &readAt, &notification.CreatedAt, &notification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Notifikasi tidak ditemukan atau bukan milik Anda.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil notification. Detail: " + err.Error(),
			},
		})
	}

	if achievementID.Valid {
		notification.AchievementID = &achievementID.String
	}

	if mongoAchievementID.Valid {
		notification.MongoAchievementID = &mongoAchievementID.String
	}

	if readAt.Valid {
		notification.ReadAt = &readAt.Time
	}

	response := modelpostgre.MarkAsReadResponse{
		Status: "success",
		Data:   notification,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func MarkAllAsReadService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	err := repositorypostgre.MarkAllAsRead(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menandai semua notification sebagai read. Detail: " + err.Error(),
			},
		})
	}

	response := modelpostgre.MarkAllAsReadResponse{
		Status: "success",
	}
	response.Data.Message = "Semua notification telah ditandai sebagai read"

	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateAchievementNotification(postgresDB *sql.DB, mongoDB *mongo.Database, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error {
	achievement, err := repositorymongo.GetAchievementByID(mongoDB, mongoAchievementID)
	if err != nil {
		return err
	}

	title := achievement.Title
	if title == "" {
		title = "Prestasi"
	}

	message := "Prestasi \"" + title + "\" telah ditolak dengan catatan: " + rejectionNote

	req := modelpostgre.CreateNotificationRequest{
		UserID:             studentUserID,
		Type:               modelpostgre.NotificationTypeAchievementRejected,
		Title:              "Prestasi Ditolak",
		Message:            message,
		AchievementID:      &achievementRefID,
		MongoAchievementID: &mongoAchievementID,
	}

	_, err = repositorypostgre.CreateNotification(postgresDB, req)
	return err
}

func CreateSubmissionNotification(postgresDB *sql.DB, mongoDB *mongo.Database, studentID string, mongoAchievementID string, achievementRefID string) error {
	var student modelpostgre.Student
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id = $1`
	err := postgresDB.QueryRow(query, studentID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return err
	}

	if student.AdvisorID == "" {
		return nil
	}

	var lecturer modelpostgre.Lecturer
	lecturerQuery := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE id = $1`
	err = postgresDB.QueryRow(lecturerQuery, student.AdvisorID).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID,
		&lecturer.Department, &lecturer.CreatedAt,
	)
	if err != nil {
		return err
	}

	achievement, err := repositorymongo.GetAchievementByID(mongoDB, mongoAchievementID)
	if err != nil {
		return err
	}

	title := achievement.Title
	if title == "" {
		title = "Prestasi"
	}

	message := "Mahasiswa bimbingan Anda telah mengajukan prestasi \"" + title + "\" untuk diverifikasi."

	req := modelpostgre.CreateNotificationRequest{
		UserID:             lecturer.UserID,
		Type:               modelpostgre.NotificationTypeAchievementSubmitted,
		Title:              "Prestasi Baru Diajukan",
		Message:            message,
		AchievementID:      &achievementRefID,
		MongoAchievementID: &mongoAchievementID,
	}

	_, err = repositorypostgre.CreateNotification(postgresDB, req)
	return err
}

