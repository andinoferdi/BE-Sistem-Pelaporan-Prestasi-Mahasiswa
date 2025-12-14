package service

// #1 proses: import library yang diperlukan untuk context, database, errors, fmt, helper, dan time
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	"time"
)

// #2 proses: definisikan interface untuk operasi achievement dengan integrasi PostgreSQL dan MongoDB
type IAchievementService interface {
	CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error)
	SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error)
	VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error)
	RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error)
	DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error)
	GetAchievements(ctx context.Context, userID string, roleID string, page, limit int, statusFilter string, achievementTypeFilter string, sortBy string, sortOrder string) (map[string]interface{}, error)
	GetAchievementsByStudentID(ctx context.Context, studentID string, page, limit int) (map[string]interface{}, error)
	GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error)
	UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error)
	GetAchievementStats(ctx context.Context) (map[string]interface{}, error)
	UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error)
	GetAchievementHistory(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error)
}

// #3 proses: struct service untuk achievement dengan dependency achievement MongoDB, achievement reference PostgreSQL, user, student, dan notification service
type AchievementService struct {
	achievementRepo     repositorymongo.IAchievementRepository
	achievementRefRepo  repositorypostgre.IAchievementReferenceRepository
	userRepo            repositorypostgre.IUserRepository
	studentRepo         repositorypostgre.IStudentRepository
	notificationService INotificationService
}

// #4 proses: constructor untuk membuat instance AchievementService baru
func NewAchievementService(
	achievementRepo repositorymongo.IAchievementRepository,
	achievementRefRepo repositorypostgre.IAchievementReferenceRepository,
	userRepo repositorypostgre.IUserRepository,
	studentRepo repositorypostgre.IStudentRepository,
	notificationService INotificationService,
) IAchievementService {
	return &AchievementService{
		achievementRepo:     achievementRepo,
		achievementRefRepo:  achievementRefRepo,
		userRepo:            userRepo,
		studentRepo:         studentRepo,
		notificationService: notificationService,
	}
}

// #5 proses: buat achievement baru di MongoDB dan reference di PostgreSQL
func (s *AchievementService) CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error) {
	// #5a proses: validasi user harus memiliki role Mahasiswa
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat membuat prestasi")
	}

	// #5b proses: ambil student ID dari user ID
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data mahasiswa tidak ditemukan. Pastikan user memiliki profil mahasiswa")
		}
		return nil, err
	}

	// #5c proses: validasi student ID format UUID
	if !helper.IsValidUUID(studentID) {
		return nil, errors.New("student ID tidak valid")
	}

	// #5d proses: validasi student ID ada di database
	_, err = s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data mahasiswa tidak ditemukan di database")
		}
		return nil, errors.New("error memvalidasi student ID: " + err.Error())
	}

	// #5e proses: validasi field wajib achievement type, title, dan description
	if req.AchievementType == "" {
		return nil, errors.New("achievement type wajib diisi")
	}

	if req.Title == "" {
		return nil, errors.New("title wajib diisi")
	}

	if req.Description == "" {
		return nil, errors.New("description wajib diisi")
	}

	// #5f proses: validasi achievement type harus salah satu dari tipe yang diizinkan
	validTypes := map[string]bool{
		"academic":      true,
		"competition":   true,
		"organization":  true,
		"publication":   true,
		"certification": true,
		"other":         true,
	}

	if !validTypes[req.AchievementType] {
		return nil, errors.New("achievement type tidak valid. Gunakan: academic, competition, organization, publication, certification, atau other")
	}

	// #5g proses: set student ID ke request
	req.StudentID = studentID

	// #5h proses: buat achievement object untuk disimpan ke MongoDB
	achievement := &modelmongo.Achievement{
		StudentID:       req.StudentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Attachments:     req.Attachments,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	// #5i proses: simpan achievement ke MongoDB
	createdAchievement, err := s.achievementRepo.CreateAchievement(ctx, achievement)
	if err != nil {
		return nil, errors.New("error menyimpan prestasi ke database: " + err.Error())
	}

	// #5j proses: buat reference di PostgreSQL dengan status draft
	refReq := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          studentID,
		MongoAchievementID: createdAchievement.ID.Hex(),
		Status:             modelpostgre.AchievementStatusDraft,
	}

	_, err = s.achievementRefRepo.CreateAchievementReference(ctx, refReq)
	if err != nil {
		return nil, errors.New("error membuat reference prestasi: " + err.Error())
	}

	// #5k proses: build response dengan achievement yang baru dibuat
	response := &modelmongo.CreateAchievementResponse{
		Status: "success",
		Data:   *createdAchievement,
	}

	return response, nil
}

// #6 proses: submit achievement untuk verifikasi, ubah status dari draft ke submitted
func (s *AchievementService) SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error) {
	// #6a proses: validasi user harus memiliki role Mahasiswa
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat submit prestasi")
	}

	// #6b proses: ambil achievement reference berdasarkan mongo ID
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #6c proses: validasi status harus draft untuk bisa di-submit
	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat di-submit jika status adalah draft")
	}

	// #6d proses: ambil student ID dan validasi ownership
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat submit prestasi milik Anda sendiri")
	}

	// #6e proses: update status jadi submitted dan set submitted_at
	now := time.Now()
	err = s.achievementRefRepo.UpdateAchievementReferenceStatus(ctx, ref.ID, modelpostgre.AchievementStatusSubmitted, &now)
	if err != nil {
		return nil, errors.New("error mengupdate status prestasi: " + err.Error())
	}

	// #6f proses: ambil reference yang sudah diupdate
	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	// #6g proses: buat notifikasi untuk dosen wali tentang submission
	err = s.notificationService.CreateSubmissionNotification(ctx, ref.StudentID, ref.MongoAchievementID, ref.ID)
	if err != nil {
		fmt.Printf("Error creating notification for submitted achievement: %v\n", err)
	}

	// #6h proses: build response dengan reference yang sudah diupdate
	response := &modelpostgre.UpdateAchievementReferenceResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

// #7 proses: verifikasi achievement oleh dosen wali, ubah status jadi verified
func (s *AchievementService) VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error) {
	// #7a proses: validasi user harus memiliki role Dosen Wali
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Dosen Wali" {
		return nil, errors.New("akses ditolak. Hanya dosen wali yang dapat memverifikasi prestasi")
	}

	// #7b proses: ambil achievement reference berdasarkan mongo ID
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #7c proses: validasi status harus submitted untuk bisa diverifikasi
	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return nil, errors.New("prestasi hanya dapat diverifikasi jika status adalah submitted")
	}

	// #7d proses: ambil lecturer dan student untuk validasi advisor relationship
	lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
		}
		return nil, err
	}

	student, err := s.studentRepo.GetStudentByID(ctx, ref.StudentID)
	if err != nil {
		return nil, errors.New("error mengambil data student: " + err.Error())
	}

	// #7e proses: validasi dosen wali adalah advisor dari student
	if student.AdvisorID != lecturer.ID {
		return nil, errors.New("akses ditolak. Anda hanya dapat memverifikasi prestasi mahasiswa bimbingan Anda")
	}

	// #7f proses: update status jadi verified dengan set verified_by dan verified_at
	err = s.achievementRefRepo.UpdateAchievementReferenceVerify(ctx, ref.ID, userID)
	if err != nil {
		return nil, errors.New("error memverifikasi prestasi: " + err.Error())
	}

	// #7g proses: ambil reference yang sudah diupdate
	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	// #7h proses: build response dengan reference yang sudah diupdate
	response := &modelpostgre.VerifyAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

// #8 proses: tolak achievement oleh dosen wali dengan catatan penolakan
func (s *AchievementService) RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error) {
	// #8a proses: validasi user harus memiliki role Dosen Wali
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Dosen Wali" {
		return nil, errors.New("akses ditolak. Hanya dosen wali yang dapat menolak prestasi")
	}

	// #8b proses: validasi rejection note tidak kosong
	if req.RejectionNote == "" {
		return nil, errors.New("rejection note wajib diisi")
	}

	// #8c proses: ambil achievement reference berdasarkan mongo ID
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #8d proses: validasi status harus submitted untuk bisa ditolak
	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return nil, errors.New("prestasi hanya dapat ditolak jika status adalah submitted")
	}

	// #8e proses: ambil lecturer dan student untuk validasi advisor relationship
	lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
		}
		return nil, err
	}

	student, err := s.studentRepo.GetStudentByID(ctx, ref.StudentID)
	if err != nil {
		return nil, errors.New("error mengambil data student: " + err.Error())
	}

	// #8f proses: validasi dosen wali adalah advisor dari student
	if student.AdvisorID != lecturer.ID {
		return nil, errors.New("akses ditolak. Anda hanya dapat menolak prestasi mahasiswa bimbingan Anda")
	}

	// #8g proses: update status jadi rejected dengan set rejection note
	err = s.achievementRefRepo.UpdateAchievementReferenceReject(ctx, ref.ID, userID, req.RejectionNote)
	if err != nil {
		return nil, errors.New("error menolak prestasi: " + err.Error())
	}

	// #8h proses: ambil reference yang sudah diupdate
	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	// #8i proses: buat notifikasi untuk student tentang penolakan
	err = s.notificationService.CreateAchievementNotification(ctx, student.UserID, ref.MongoAchievementID, ref.ID, req.RejectionNote)
	if err != nil {
		fmt.Printf("Error creating notification for rejected achievement: %v\n", err)
	}

	// #8j proses: build response dengan reference yang sudah diupdate
	response := &modelpostgre.RejectAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

// #9 proses: hapus achievement, hanya bisa jika status draft dan milik user sendiri
func (s *AchievementService) DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error) {
	// #9a proses: validasi user harus memiliki role Mahasiswa
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat menghapus prestasi")
	}

	// #9b proses: ambil achievement reference berdasarkan mongo ID
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #9c proses: validasi status harus draft untuk bisa dihapus
	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat dihapus jika status adalah draft")
	}

	// #9d proses: ambil student ID dan validasi ownership
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat menghapus prestasi milik Anda sendiri")
	}

	// #9e proses: soft delete achievement di MongoDB
	err = s.achievementRepo.DeleteAchievement(ctx, mongoID)
	if err != nil {
		return nil, errors.New("error menghapus prestasi dari database: " + err.Error())
	}

	// #9f proses: update status reference jadi deleted
	err = s.achievementRefRepo.UpdateAchievementReferenceStatus(ctx, ref.ID, modelpostgre.AchievementStatusDeleted, nil)
	if err != nil {
		return nil, errors.New("error mengupdate status prestasi menjadi deleted: " + err.Error())
	}

	// #9g proses: build response sukses
	response := &modelmongo.DeleteAchievementResponse{
		Status: "success",
	}

	return response, nil
}

// #10 proses: ambil achievements dengan pagination dan filtering berdasarkan role user
func (s *AchievementService) GetAchievements(ctx context.Context, userID string, roleID string, page, limit int, statusFilter string, achievementTypeFilter string, sortBy string, sortOrder string) (map[string]interface{}, error) {
	// #10a proses: validasi dan set default untuk page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// #10b proses: ambil role name untuk menentukan query yang sesuai
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	var references []modelpostgre.AchievementReference
	var total int

	// #10c proses: ambil references berdasarkan role, Mahasiswa hanya lihat milik sendiri
	if roleName == "Mahasiswa" {
		// #10d proses: ambil student dan query references milik student
		student, err := s.studentRepo.GetStudentByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data student: " + err.Error())
		}

		references, total, err = s.achievementRefRepo.GetAchievementReferenceByStudentIDPaginated(ctx, student.ID, page, limit)
		if err != nil {
			return nil, errors.New("error mengambil achievement references: " + err.Error())
		}
	} else if roleName == "Dosen Wali" {
		// #10e proses: ambil lecturer dan query references dari mahasiswa bimbingan
		lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
			}
			return nil, err
		}

		references, total, err = s.achievementRefRepo.GetAchievementReferencesByAdvisorIDPaginated(ctx, lecturer.ID, page, limit)
		if err != nil {
			return nil, errors.New("error mengambil achievement references: " + err.Error())
		}
	} else if roleName == "Admin" {
		// #10f proses: admin bisa lihat semua dengan filter dan sorting
		references, total, err = s.achievementRefRepo.GetAllAchievementReferencesPaginated(ctx, page, limit, statusFilter, sortBy, sortOrder)
		if err != nil {
			return nil, errors.New("error mengambil achievement references: " + err.Error())
		}
	} else {
		return nil, errors.New("akses ditolak. Role tidak memiliki akses untuk melihat prestasi")
	}

	// #10g proses: jika tidak ada references, return response kosong
	if len(references) == 0 {
		totalPages := 0
		if total > 0 {
			totalPages = (total + limit - 1) / limit
		}
		return map[string]interface{}{
			"status": "success",
			"data":   []modelmongo.Achievement{},
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		}, nil
	}

	// #10h proses: kumpulkan semua mongo IDs dari references
	var mongoIDs []string
	for _, ref := range references {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}

	// #10i proses: ambil achievements dari MongoDB berdasarkan IDs
	achievements, err := s.achievementRepo.GetAchievementsByIDs(ctx, mongoIDs)
	if err != nil {
		return nil, errors.New("error mengambil achievements dari MongoDB: " + err.Error())
	}

	// #10j proses: buat map reference dan student untuk lookup cepat
	referenceMap := make(map[string]modelpostgre.AchievementReference)
	for _, ref := range references {
		referenceMap[ref.MongoAchievementID] = ref
	}

	studentMap := make(map[string]modelpostgre.Student)
	for _, ref := range references {
		if _, exists := studentMap[ref.StudentID]; !exists {
			student, err := s.studentRepo.GetStudentByID(ctx, ref.StudentID)
			if err == nil && student != nil {
				studentMap[ref.StudentID] = *student
			}
		}
	}

	// #10k proses: build result dengan gabungkan data MongoDB dan PostgreSQL
	var result []map[string]interface{}
	for _, achievement := range achievements {
		// #10l proses: filter berdasarkan achievement type jika ada
		if achievementTypeFilter != "" && achievement.AchievementType != achievementTypeFilter {
			continue
		}

		ref := referenceMap[achievement.ID.Hex()]
		item := map[string]interface{}{
			"id":              achievement.ID.Hex(),
			"studentId":       achievement.StudentID,
			"achievementType": achievement.AchievementType,
			"title":           achievement.Title,
			"description":     achievement.Description,
			"details":         achievement.Details,
			"attachments":     achievement.Attachments,
			"tags":            achievement.Tags,
			"points":          achievement.Points,
			"createdAt":       achievement.CreatedAt.Format(time.RFC3339),
			"updatedAt":       achievement.UpdatedAt.Format(time.RFC3339),
			"status":          ref.Status,
		}

		if student, exists := studentMap[ref.StudentID]; exists {
			item["student_name"] = student.FullName
		}

		if ref.SubmittedAt != nil {
			item["submitted_at"] = ref.SubmittedAt.Format(time.RFC3339)
		}
		if ref.VerifiedAt != nil {
			item["verified_at"] = ref.VerifiedAt.Format(time.RFC3339)
		}
		if ref.VerifiedBy != nil {
			verifiedByUser, err := s.userRepo.FindUserByID(ctx, *ref.VerifiedBy)
			if err == nil && verifiedByUser != nil {
				item["verified_by"] = verifiedByUser.FullName
			} else {
				item["verified_by"] = *ref.VerifiedBy
			}
		}
		if ref.RejectionNote != nil {
			item["rejection_note"] = *ref.RejectionNote
		}

		result = append(result, item)
	}

	// #10m proses: jika ada filter achievement type, adjust total
	if achievementTypeFilter != "" {
		total = len(result)
	}

	// #10n proses: hitung total pages dan build response dengan pagination
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return map[string]interface{}{
		"status": "success",
		"data":   result,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}, nil
}

// #11 proses: ambil achievements berdasarkan student ID dengan pagination
func (s *AchievementService) GetAchievementsByStudentID(ctx context.Context, studentID string, page, limit int) (map[string]interface{}, error) {
	// #11a proses: validasi dan set default untuk page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	references, total, err := s.achievementRefRepo.GetAchievementReferenceByStudentIDPaginated(ctx, studentID, page, limit)
	if err != nil {
		return nil, errors.New("error mengambil achievement references: " + err.Error())
	}

	if len(references) == 0 {
		totalPages := 0
		if total > 0 {
			totalPages = (total + limit - 1) / limit
		}
		return map[string]interface{}{
			"status": "success",
			"data":   []modelmongo.Achievement{},
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		}, nil
	}

	var mongoIDs []string
	for _, ref := range references {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}

	achievements, err := s.achievementRepo.GetAchievementsByIDs(ctx, mongoIDs)
	if err != nil {
		return nil, errors.New("error mengambil achievements dari MongoDB: " + err.Error())
	}

	referenceMap := make(map[string]modelpostgre.AchievementReference)
	for _, ref := range references {
		referenceMap[ref.MongoAchievementID] = ref
	}

	var result []map[string]interface{}
	for _, achievement := range achievements {
		ref := referenceMap[achievement.ID.Hex()]
		item := map[string]interface{}{
			"id":              achievement.ID.Hex(),
			"studentId":       achievement.StudentID,
			"achievementType": achievement.AchievementType,
			"title":           achievement.Title,
			"description":     achievement.Description,
			"details":         achievement.Details,
			"attachments":     achievement.Attachments,
			"tags":            achievement.Tags,
			"points":          achievement.Points,
			"createdAt":       achievement.CreatedAt.Format(time.RFC3339),
			"updatedAt":       achievement.UpdatedAt.Format(time.RFC3339),
			"status":          ref.Status,
		}

		if ref.SubmittedAt != nil {
			item["submitted_at"] = ref.SubmittedAt.Format(time.RFC3339)
		}
		if ref.VerifiedAt != nil {
			item["verified_at"] = ref.VerifiedAt.Format(time.RFC3339)
		}
		if ref.VerifiedBy != nil {
			verifiedByUser, err := s.userRepo.FindUserByID(ctx, *ref.VerifiedBy)
			if err == nil && verifiedByUser != nil {
				item["verified_by"] = verifiedByUser.FullName
			} else {
				item["verified_by"] = *ref.VerifiedBy
			}
		}
		if ref.RejectionNote != nil {
			item["rejection_note"] = *ref.RejectionNote
		}

		result = append(result, item)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return map[string]interface{}{
		"status": "success",
		"data":   result,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}, nil
}

// #15 proses: ambil achievement detail berdasarkan ID dengan validasi akses
func (s *AchievementService) GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	// #15a proses: ambil achievement reference untuk validasi akses
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #15b proses: validasi akses berdasarkan role
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName == "Mahasiswa" {
		// #15c proses: mahasiswa hanya bisa lihat prestasi sendiri
		studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
		}

		if ref.StudentID != studentID {
			return nil, errors.New("akses ditolak. Anda hanya dapat melihat prestasi milik Anda sendiri")
		}
	} else if roleName == "Dosen Wali" {
		// #15d proses: dosen wali hanya bisa lihat prestasi mahasiswa bimbingan
		lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
			}
			return nil, err
		}

		student, err := s.studentRepo.GetStudentByID(ctx, ref.StudentID)
		if err != nil {
			return nil, errors.New("error mengambil data student: " + err.Error())
		}

		if student.AdvisorID != lecturer.ID {
			return nil, errors.New("akses ditolak. Anda hanya dapat melihat prestasi mahasiswa bimbingan Anda")
		}
	} else if roleName != "Admin" {
		return nil, errors.New("akses ditolak. Role tidak memiliki akses untuk melihat prestasi")
	}

	// #15e proses: ambil achievement dari MongoDB
	achievement, err := s.achievementRepo.GetAchievementByID(ctx, mongoID)
	if err != nil {
		return nil, errors.New("error mengambil achievement dari database: " + err.Error())
	}
	if achievement == nil {
		return nil, errors.New("prestasi tidak ditemukan")
	}

	result := map[string]interface{}{
		"id":              achievement.ID.Hex(),
		"studentId":       achievement.StudentID,
		"achievementType": achievement.AchievementType,
		"title":           achievement.Title,
		"description":     achievement.Description,
		"details":         achievement.Details,
		"attachments":     achievement.Attachments,
		"tags":            achievement.Tags,
		"points":          achievement.Points,
		"createdAt":       achievement.CreatedAt.Format(time.RFC3339),
		"updatedAt":       achievement.UpdatedAt.Format(time.RFC3339),
		"status":          ref.Status,
	}

	if ref.SubmittedAt != nil {
		result["submitted_at"] = ref.SubmittedAt.Format(time.RFC3339)
	}

	if ref.VerifiedAt != nil {
		result["verified_at"] = ref.VerifiedAt.Format(time.RFC3339)
	}

	if ref.VerifiedBy != nil {
		verifiedByUser, err := s.userRepo.FindUserByID(ctx, *ref.VerifiedBy)
		if err == nil && verifiedByUser != nil {
			result["verified_by"] = verifiedByUser.FullName
		} else {
			result["verified_by"] = *ref.VerifiedBy
		}
	}

	if ref.RejectionNote != nil {
		result["rejection_note"] = *ref.RejectionNote
	}

	return map[string]interface{}{
		"status": "success",
		"data":   result,
	}, nil
}

// #16 proses: update achievement, hanya bisa jika status draft dan milik user sendiri
func (s *AchievementService) UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error) {
	// #16a proses: validasi user harus memiliki role Mahasiswa
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat mengupdate prestasi")
	}

	// #16b proses: ambil achievement reference untuk validasi status dan ownership
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #16c proses: validasi status harus draft untuk bisa diupdate
	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat diupdate jika status adalah draft")
	}

	// #16d proses: validasi ownership, hanya student pemilik yang bisa update
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat mengupdate prestasi milik Anda sendiri")
	}

	// #16e proses: jika ada achievement type, validasi harus salah satu tipe yang diizinkan
	if req.AchievementType != "" {
		validTypes := map[string]bool{
			"academic":      true,
			"competition":   true,
			"organization":  true,
			"publication":   true,
			"certification": true,
			"other":         true,
		}

		if !validTypes[req.AchievementType] {
			return nil, errors.New("achievement type tidak valid. Gunakan: academic, competition, organization, publication, certification, atau other")
		}
	}

	// #16f proses: update achievement di MongoDB
	updatedAchievement, err := s.achievementRepo.UpdateAchievement(ctx, mongoID, req)
	if err != nil {
		return nil, errors.New("error mengupdate prestasi di database: " + err.Error())
	}

	// #16g proses: build result dengan gabungkan data MongoDB dan reference
	result := map[string]interface{}{
		"id":              updatedAchievement.ID.Hex(),
		"studentId":       updatedAchievement.StudentID,
		"achievementType": updatedAchievement.AchievementType,
		"title":           updatedAchievement.Title,
		"description":     updatedAchievement.Description,
		"details":         updatedAchievement.Details,
		"attachments":     updatedAchievement.Attachments,
		"tags":            updatedAchievement.Tags,
		"points":          updatedAchievement.Points,
		"createdAt":       updatedAchievement.CreatedAt.Format(time.RFC3339),
		"updatedAt":       updatedAchievement.UpdatedAt.Format(time.RFC3339),
		"status":          ref.Status,
	}

	// #16h proses: build response dengan achievement yang sudah diupdate
	return map[string]interface{}{
		"status": "success",
		"data":   result,
	}, nil
}

// #12 proses: ambil statistik achievement total dan verified
func (s *AchievementService) GetAchievementStats(ctx context.Context) (map[string]interface{}, error) {
	// #12a proses: ambil total dan verified dari repository
	total, verified, err := s.achievementRefRepo.GetAchievementStats(ctx)
	if err != nil {
		return nil, errors.New("error mengambil statistik prestasi: " + err.Error())
	}

	// #12b proses: hitung persentase verified
	percentage := 0
	if total > 0 {
		percentage = int((float64(verified) / float64(total)) * 100)
	}

	// #12c proses: build response dengan statistik
	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"total":      total,
			"verified":   verified,
			"percentage": percentage,
		},
	}, nil
}

// #13 proses: upload file attachment ke achievement
func (s *AchievementService) UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error) {
	// #13a proses: ambil achievement reference untuk validasi
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #13b proses: validasi ownership, hanya student pemilik yang bisa upload
	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat menambahkan attachment ke prestasi milik Anda sendiri")
	}

	// #13c proses: validasi status harus draft untuk bisa upload attachment
	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("attachment hanya dapat ditambahkan jika status prestasi adalah draft")
	}

	// #13d proses: buat attachment object dan tambahkan ke achievement
	attachment := modelmongo.Attachment{
		FileName:   fileName,
		FileURL:    fileURL,
		FileType:   fileType,
		UploadedAt: time.Now(),
	}

	_, err = s.achievementRepo.AddAttachmentToAchievement(ctx, mongoID, attachment)
	if err != nil {
		return nil, errors.New("error menambahkan attachment ke prestasi: " + err.Error())
	}

	return &attachment, nil
}

// #14 proses: ambil history perubahan status achievement
func (s *AchievementService) GetAchievementHistory(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	// #14a proses: ambil achievement reference
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	// #14b proses: validasi akses berdasarkan role
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName == "Mahasiswa" {
		// #14c proses: mahasiswa hanya bisa lihat history prestasi sendiri
		studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
		}

		if ref.StudentID != studentID {
			return nil, errors.New("akses ditolak. Anda hanya dapat melihat history prestasi milik Anda sendiri")
		}
	} else if roleName == "Dosen Wali" {
		// #14d proses: dosen wali hanya bisa lihat history prestasi mahasiswa bimbingan
		lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
			}
			return nil, err
		}

		student, err := s.studentRepo.GetStudentByID(ctx, ref.StudentID)
		if err != nil {
			return nil, errors.New("error mengambil data student: " + err.Error())
		}

		if student.AdvisorID != lecturer.ID {
			return nil, errors.New("akses ditolak. Anda hanya dapat melihat history prestasi mahasiswa bimbingan Anda")
		}
	} else if roleName != "Admin" {
		return nil, errors.New("akses ditolak. Role tidak memiliki akses untuk melihat history prestasi")
	}

	// #14e proses: build history dari status changes berdasarkan timestamp di reference
	var history []map[string]interface{}

	// #14f proses: tambahkan entry untuk status draft saat pertama dibuat
	draftEntry := map[string]interface{}{
		"status":          modelpostgre.AchievementStatusDraft,
		"changed_at":      ref.CreatedAt.Format(time.RFC3339),
		"changed_by":      nil,
		"changed_by_name": nil,
		"note":            nil,
	}
	history = append(history, draftEntry)

	// #14g proses: jika ada submitted_at, tambahkan entry untuk status submitted
	if ref.SubmittedAt != nil {
		submittedEntry := map[string]interface{}{
			"status":          modelpostgre.AchievementStatusSubmitted,
			"changed_at":      ref.SubmittedAt.Format(time.RFC3339),
			"changed_by":      nil,
			"changed_by_name": nil,
			"note":            nil,
		}
		history = append(history, submittedEntry)
	}

	// #14h proses: jika status verified, tambahkan entry dengan info verified_by
	if ref.Status == modelpostgre.AchievementStatusVerified && ref.VerifiedAt != nil {
		var verifiedByName *string
		if ref.VerifiedBy != nil {
			verifiedByUser, err := s.userRepo.FindUserByID(ctx, *ref.VerifiedBy)
			if err == nil && verifiedByUser != nil {
				name := verifiedByUser.FullName
				verifiedByName = &name
			}
		}

		verifiedEntry := map[string]interface{}{
			"status":          modelpostgre.AchievementStatusVerified,
			"changed_at":      ref.VerifiedAt.Format(time.RFC3339),
			"changed_by":      ref.VerifiedBy,
			"changed_by_name": verifiedByName,
			"note":            nil,
		}
		history = append(history, verifiedEntry)
	}

	// #14i proses: jika status rejected, tambahkan entry dengan rejection note
	if ref.Status == modelpostgre.AchievementStatusRejected {
		var verifiedByName *string
		if ref.VerifiedBy != nil {
			verifiedByUser, err := s.userRepo.FindUserByID(ctx, *ref.VerifiedBy)
			if err == nil && verifiedByUser != nil {
				name := verifiedByUser.FullName
				verifiedByName = &name
			}
		}

		rejectedEntry := map[string]interface{}{
			"status":          modelpostgre.AchievementStatusRejected,
			"changed_at":      ref.UpdatedAt.Format(time.RFC3339),
			"changed_by":      ref.VerifiedBy,
			"changed_by_name": verifiedByName,
			"note":            ref.RejectionNote,
		}
		history = append(history, rejectedEntry)
	}

	// #14j proses: jika status deleted, tambahkan entry untuk status deleted
	if ref.Status == modelpostgre.AchievementStatusDeleted {
		deletedEntry := map[string]interface{}{
			"status":          modelpostgre.AchievementStatusDeleted,
			"changed_at":      ref.UpdatedAt.Format(time.RFC3339),
			"changed_by":      nil,
			"changed_by_name": nil,
			"note":            nil,
		}
		history = append(history, deletedEntry)
	}

	// #14k proses: build response dengan history
	return map[string]interface{}{
		"status": "success",
		"data":   history,
	}, nil
}
