package service

// #1 proses: import library yang diperlukan untuk context, database, errors, dan repository
import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
)

// #2 proses: definisikan interface untuk operasi notifikasi
type INotificationService interface {
	GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error)
	MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error)
	MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error)
	CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error
	CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error
}

// #3 proses: struct service untuk notifikasi dengan dependency notification, student, user, dan achievement repository
type NotificationService struct {
	notifRepo       repositorypostgre.INotificationRepository
	studentRepo     repositorypostgre.IStudentRepository
	userRepo        repositorypostgre.IUserRepository
	achievementRepo repositorymongo.IAchievementRepository
}

// #4 proses: constructor untuk membuat instance NotificationService baru
func NewNotificationService(
	notifRepo repositorypostgre.INotificationRepository,
	studentRepo repositorypostgre.IStudentRepository,
	userRepo repositorypostgre.IUserRepository,
	achievementRepo repositorymongo.IAchievementRepository,
) INotificationService {
	return &NotificationService{
		notifRepo:       notifRepo,
		studentRepo:     studentRepo,
		userRepo:        userRepo,
		achievementRepo: achievementRepo,
	}
}

// #5 proses: ambil notifikasi user dengan pagination
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error) {
	// #5a proses: validasi dan set default untuk page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// #5b proses: ambil notifikasi dengan pagination dari repository
	notifications, total, err := s.notifRepo.GetNotificationsByUserIDPaginated(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	// #5c proses: hitung total pages
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	// #5d proses: build response dengan data dan pagination info
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

// #6 proses: ambil jumlah notifikasi yang belum dibaca
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error) {
	// #6a proses: ambil count notifikasi belum dibaca dari repository
	count, err := s.notifRepo.GetUnreadCountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// #6b proses: build response dengan count
	response := &modelpostgre.GetUnreadCountResponse{
		Status: "success",
	}
	response.Data.Count = count

	return response, nil
}

// #7 proses: tandai notifikasi tertentu sebagai sudah dibaca
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error) {
	// #7a proses: validasi notification ID tidak kosong
	if notificationID == "" {
		return nil, errors.New("ID notification wajib diisi")
	}

	// #7b proses: update notifikasi jadi sudah dibaca di repository
	err := s.notifRepo.MarkAsRead(ctx, notificationID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("notifikasi tidak ditemukan atau bukan milik Anda")
		}
		return nil, err
	}

	// #7c proses: ambil notifikasi yang sudah diupdate untuk return
	notifications, _, err := s.notifRepo.GetNotificationsByUserIDPaginated(ctx, userID, 1, 1000)
	if err != nil {
		return nil, err
	}

	// #7d proses: cari notifikasi berdasarkan ID
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

	// #7e proses: build response dengan notifikasi yang sudah diupdate
	response := &modelpostgre.MarkAsReadResponse{
		Status: "success",
		Data:   *notification,
	}

	return response, nil
}

// #8 proses: tandai semua notifikasi user sebagai sudah dibaca
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error) {
	// #8a proses: update semua notifikasi jadi sudah dibaca di repository
	err := s.notifRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return nil, err
	}

	// #8b proses: build response dengan pesan sukses
	response := &modelpostgre.MarkAllAsReadResponse{
		Status: "success",
	}
	response.Data.Message = "Semua notification telah ditandai sebagai read"

	return response, nil
}

// #9 proses: buat notifikasi untuk student ketika prestasi ditolak
func (s *NotificationService) CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error {
	// #9a proses: ambil achievement dari MongoDB untuk ambil title
	achievement, err := s.achievementRepo.GetAchievementByID(ctx, mongoAchievementID)
	if err != nil {
		return err
	}
	if achievement == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// #9b proses: set title dari achievement atau gunakan default
	title := achievement.Title
	if title == "" {
		title = "Prestasi"
	}

	// #9c proses: buat message notifikasi dengan catatan penolakan
	message := "Prestasi \"" + title + "\" telah ditolak dengan catatan: " + rejectionNote

	// #9d proses: buat request notifikasi dan simpan ke database
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

// #10 proses: buat notifikasi untuk dosen wali ketika mahasiswa submit prestasi
func (s *NotificationService) CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error {
	// #10a proses: ambil student untuk dapat advisor ID
	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	// #10b proses: jika student tidak punya advisor, tidak perlu buat notifikasi
	if student.AdvisorID == "" {
		return nil
	}

	// #10c proses: ambil lecturer berdasarkan advisor ID untuk dapat user ID
	lecturer, err := s.userRepo.GetLecturerByID(ctx, student.AdvisorID)
	if err != nil {
		return err
	}

	// #10d proses: ambil achievement dari MongoDB untuk ambil title
	achievement, err := s.achievementRepo.GetAchievementByID(ctx, mongoAchievementID)
	if err != nil {
		return err
	}
	if achievement == nil {
		return errors.New("prestasi tidak ditemukan")
	}

	// #10e proses: set title dari achievement atau gunakan default
	title := achievement.Title
	if title == "" {
		title = "Prestasi"
	}

	// #10f proses: buat message notifikasi untuk dosen wali
	message := "Mahasiswa bimbingan Anda telah mengajukan prestasi \"" + title + "\" untuk diverifikasi."

	// #10g proses: buat request notifikasi dan simpan ke database untuk dosen wali
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
