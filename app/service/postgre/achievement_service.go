package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"time"
)

type IAchievementService interface {
	CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error)
	SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error)
	VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error)
	RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error)
	DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error)
	GetAchievements(ctx context.Context, userID string, roleID string, page, limit int) (map[string]interface{}, error)
	GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error)
	UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error)
	GetAchievementStats(ctx context.Context) (map[string]interface{}, error)
	UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error)
}

type AchievementService struct {
	achievementRepo      repositorymongo.IAchievementRepository
	achievementRefRepo   repositorypostgre.IAchievementReferenceRepository
	userRepo             repositorypostgre.IUserRepository
	studentRepo          repositorypostgre.IStudentRepository
	notificationService  INotificationService
}

func NewAchievementService(
	achievementRepo repositorymongo.IAchievementRepository,
	achievementRefRepo repositorypostgre.IAchievementReferenceRepository,
	userRepo repositorypostgre.IUserRepository,
	studentRepo repositorypostgre.IStudentRepository,
	notificationService INotificationService,
) IAchievementService {
	return &AchievementService{
		achievementRepo:     achievementRepo,
		achievementRefRepo:   achievementRefRepo,
		userRepo:            userRepo,
		studentRepo:         studentRepo,
		notificationService: notificationService,
	}
}

func (s *AchievementService) CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat membuat prestasi")
	}

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data mahasiswa tidak ditemukan. Pastikan user memiliki profil mahasiswa")
		}
		return nil, err
	}

	if req.AchievementType == "" {
		return nil, errors.New("achievement type wajib diisi")
	}

	if req.Title == "" {
		return nil, errors.New("title wajib diisi")
	}

	if req.Description == "" {
		return nil, errors.New("description wajib diisi")
	}

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

	req.StudentID = studentID

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

	createdAchievement, err := s.achievementRepo.CreateAchievement(ctx, achievement)
	if err != nil {
		return nil, errors.New("error menyimpan prestasi ke database: " + err.Error())
	}

	refReq := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          studentID,
		MongoAchievementID: createdAchievement.ID.Hex(),
		Status:             modelpostgre.AchievementStatusDraft,
	}

	_, err = s.achievementRefRepo.CreateAchievementReference(ctx, refReq)
	if err != nil {
		return nil, errors.New("error membuat reference prestasi: " + err.Error())
	}

	response := &modelmongo.CreateAchievementResponse{
		Status: "success",
		Data:   *createdAchievement,
	}

	return response, nil
}

func (s *AchievementService) SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat submit prestasi")
	}

	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat di-submit jika status adalah draft")
	}

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat submit prestasi milik Anda sendiri")
	}

	now := time.Now()
	err = s.achievementRefRepo.UpdateAchievementReferenceStatus(ctx, ref.ID, modelpostgre.AchievementStatusSubmitted, &now)
	if err != nil {
		return nil, errors.New("error mengupdate status prestasi: " + err.Error())
	}

	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	err = s.notificationService.CreateSubmissionNotification(ctx, ref.StudentID, ref.MongoAchievementID, ref.ID)
	if err != nil {
		fmt.Printf("Error creating notification for submitted achievement: %v\n", err)
	}

	response := &modelpostgre.UpdateAchievementReferenceResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

func (s *AchievementService) VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Dosen Wali" {
		return nil, errors.New("akses ditolak. Hanya dosen wali yang dapat memverifikasi prestasi")
	}

	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return nil, errors.New("prestasi hanya dapat diverifikasi jika status adalah submitted")
	}

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
		return nil, errors.New("akses ditolak. Anda hanya dapat memverifikasi prestasi mahasiswa bimbingan Anda")
	}

	err = s.achievementRefRepo.UpdateAchievementReferenceVerify(ctx, ref.ID, userID)
	if err != nil {
		return nil, errors.New("error memverifikasi prestasi: " + err.Error())
	}

	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	response := &modelpostgre.VerifyAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

func (s *AchievementService) RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Dosen Wali" {
		return nil, errors.New("akses ditolak. Hanya dosen wali yang dapat menolak prestasi")
	}

	if req.RejectionNote == "" {
		return nil, errors.New("rejection note wajib diisi")
	}

	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return nil, errors.New("prestasi hanya dapat ditolak jika status adalah submitted")
	}

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
		return nil, errors.New("akses ditolak. Anda hanya dapat menolak prestasi mahasiswa bimbingan Anda")
	}

	err = s.achievementRefRepo.UpdateAchievementReferenceReject(ctx, ref.ID, userID, req.RejectionNote)
	if err != nil {
		return nil, errors.New("error menolak prestasi: " + err.Error())
	}

	updatedRef, err := s.achievementRefRepo.GetAchievementReferenceByID(ctx, ref.ID)
	if err != nil {
		return nil, errors.New("error mengambil data prestasi yang diupdate: " + err.Error())
	}

	err = s.notificationService.CreateAchievementNotification(ctx, student.UserID, ref.MongoAchievementID, ref.ID, req.RejectionNote)
	if err != nil {
		fmt.Printf("Error creating notification for rejected achievement: %v\n", err)
	}

	response := &modelpostgre.RejectAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return response, nil
}

func (s *AchievementService) DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat menghapus prestasi")
	}

	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat dihapus jika status adalah draft")
	}

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat menghapus prestasi milik Anda sendiri")
	}

	err = s.achievementRepo.DeleteAchievement(ctx, mongoID)
	if err != nil {
		return nil, errors.New("error menghapus prestasi dari database: " + err.Error())
	}

	err = s.achievementRefRepo.UpdateAchievementReferenceStatus(ctx, ref.ID, modelpostgre.AchievementStatusDeleted, nil)
	if err != nil {
		return nil, errors.New("error mengupdate status prestasi menjadi deleted: " + err.Error())
	}

	response := &modelmongo.DeleteAchievementResponse{
		Status: "success",
	}

	return response, nil
}

func (s *AchievementService) GetAchievements(ctx context.Context, userID string, roleID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	var references []modelpostgre.AchievementReference
	var total int

	if roleName == "Mahasiswa" {
		student, err := s.studentRepo.GetStudentByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data student: " + err.Error())
		}

		references, total, err = s.achievementRefRepo.GetAchievementReferenceByStudentIDPaginated(ctx, student.ID, page, limit)
		if err != nil {
			return nil, errors.New("error mengambil achievement references: " + err.Error())
		}
	} else if roleName == "Dosen Wali" {
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
		references, total, err = s.achievementRefRepo.GetAllAchievementReferencesPaginated(ctx, page, limit)
		if err != nil {
			return nil, errors.New("error mengambil achievement references: " + err.Error())
		}
	} else {
		return nil, errors.New("akses ditolak. Role tidak memiliki akses untuk melihat prestasi")
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

func (s *AchievementService) GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName == "Mahasiswa" {
		studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
		}

		if ref.StudentID != studentID {
			return nil, errors.New("akses ditolak. Anda hanya dapat melihat prestasi milik Anda sendiri")
		}
	} else if roleName == "Dosen Wali" {
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

func (s *AchievementService) UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error) {
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("akses ditolak. Hanya mahasiswa yang dapat mengupdate prestasi")
	}

	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("prestasi hanya dapat diupdate jika status adalah draft")
	}

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat mengupdate prestasi milik Anda sendiri")
	}

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

	updatedAchievement, err := s.achievementRepo.UpdateAchievement(ctx, mongoID, req)
	if err != nil {
		return nil, errors.New("error mengupdate prestasi di database: " + err.Error())
	}

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

	return map[string]interface{}{
		"status": "success",
		"data":   result,
	}, nil
}

func (s *AchievementService) GetAchievementStats(ctx context.Context) (map[string]interface{}, error) {
	total, verified, err := s.achievementRefRepo.GetAchievementStats(ctx)
	if err != nil {
		return nil, errors.New("error mengambil statistik prestasi: " + err.Error())
	}

	percentage := 0
	if total > 0 {
		percentage = int((float64(verified) / float64(total)) * 100)
	}

	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"total":      total,
			"verified":   verified,
			"percentage": percentage,
		},
	}, nil
}

func (s *AchievementService) UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error) {
	ref, err := s.achievementRefRepo.GetAchievementReferenceByMongoID(ctx, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("prestasi tidak ditemukan")
		}
		return nil, err
	}

	studentID, err := s.studentRepo.GetStudentIDByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error mengambil data mahasiswa: " + err.Error())
	}

	if ref.StudentID != studentID {
		return nil, errors.New("akses ditolak. Anda hanya dapat menambahkan attachment ke prestasi milik Anda sendiri")
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return nil, errors.New("attachment hanya dapat ditambahkan jika status prestasi adalah draft")
	}

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
