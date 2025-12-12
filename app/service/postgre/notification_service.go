package service

import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
)

type INotificationService interface {
	GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error)
	MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error)
	MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error)
	CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error
	CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error
}

type NotificationService struct {
	notifRepo repositorypostgre.INotificationRepository
	studentRepo repositorypostgre.IStudentRepository
	userRepo repositorypostgre.IUserRepository
	achievementRepo repositorymongo.IAchievementRepository
}

func NewNotificationService(
	notifRepo repositorypostgre.INotificationRepository,
	studentRepo repositorypostgre.IStudentRepository,
	userRepo repositorypostgre.IUserRepository,
	achievementRepo repositorymongo.IAchievementRepository,
) INotificationService {
	return &NotificationService{
		notifRepo:      notifRepo,
		studentRepo:    studentRepo,
		userRepo:       userRepo,
		achievementRepo: achievementRepo,
	}
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	notifications, total, err := s.notifRepo.GetNotificationsByUserIDPaginated(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	response := &modelpostgre.GetNotificationsResponse{
		Status: "success",
		Data:   notifications,
	}
	response.Pagination.Page = page
	response.Pagination.Limit = limit
	response.Pagination.Total = total
	response.Pagination.TotalPages = totalPages

	return response, nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error) {
	count, err := s.notifRepo.GetUnreadCountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &modelpostgre.GetUnreadCountResponse{
		Status: "success",
	}
	response.Data.Count = count

	return response, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error) {
	if notificationID == "" {
		return nil, errors.New("ID notification wajib diisi")
	}

	err := s.notifRepo.MarkAsRead(ctx, notificationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("notifikasi tidak ditemukan atau bukan milik Anda")
		}
		return nil, err
	}

	notifications, _, err := s.notifRepo.GetNotificationsByUserIDPaginated(ctx, userID, 1, 1000)
	if err != nil {
		return nil, err
	}

	var notification *modelpostgre.Notification
	for i := range notifications {
		if notifications[i].ID == notificationID {
			notification = &notifications[i]
			break
		}
	}

	if notification == nil {
		return nil, errors.New("notifikasi tidak ditemukan")
	}

	response := &modelpostgre.MarkAsReadResponse{
		Status: "success",
		Data:   *notification,
	}

	return response, nil
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error) {
	err := s.notifRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &modelpostgre.MarkAllAsReadResponse{
		Status: "success",
	}
	response.Data.Message = "Semua notification telah ditandai sebagai read"

	return response, nil
}

func (s *NotificationService) CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error {
	achievement, err := s.achievementRepo.GetAchievementByID(ctx, mongoAchievementID)
	if err != nil {
		return err
	}
	if achievement == nil {
		return errors.New("prestasi tidak ditemukan")
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

	_, err = s.notifRepo.CreateNotification(ctx, req)
	return err
}

func (s *NotificationService) CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error {
	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	if student.AdvisorID == "" {
		return nil
	}

	lecturer, err := s.userRepo.GetLecturerByID(ctx, student.AdvisorID)
	if err != nil {
		return err
	}

	achievement, err := s.achievementRepo.GetAchievementByID(ctx, mongoAchievementID)
	if err != nil {
		return err
	}
	if achievement == nil {
		return errors.New("prestasi tidak ditemukan")
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

	_, err = s.notifRepo.CreateNotification(ctx, req)
	return err
}
