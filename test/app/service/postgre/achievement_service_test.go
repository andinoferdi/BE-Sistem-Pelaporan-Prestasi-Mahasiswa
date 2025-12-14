package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockAchievementRepo struct {
	byID                 *modelmongo.Achievement
	byStudentID          []modelmongo.Achievement
	byIDs                []modelmongo.Achievement
	byType               map[string]int
	competitionLevelDist map[string]int
	topStudents          []repositorymongo.TopStudentResult
	err                  error
	createErr            error
	updateErr            error
	deleteErr            error
	addAttachmentErr     error
}

func (m *mockAchievementRepo) CreateAchievement(ctx context.Context, achievement *modelmongo.Achievement) (*modelmongo.Achievement, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	achievement.ID = primitive.NewObjectID()
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	return achievement, nil
}

func (m *mockAchievementRepo) GetAchievementByID(ctx context.Context, id string) (*modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockAchievementRepo) UpdateAchievement(ctx context.Context, id string, req modelmongo.UpdateAchievementRequest) (*modelmongo.Achievement, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if m.err != nil {
		return nil, m.err
	}
	if m.byID == nil {
		return nil, nil
	}
	return m.byID, nil
}

func (m *mockAchievementRepo) DeleteAchievement(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return m.err
}

func (m *mockAchievementRepo) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byStudentID, nil
}

func (m *mockAchievementRepo) GetAchievementsByIDs(ctx context.Context, ids []string) ([]modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byIDs, nil
}

func (m *mockAchievementRepo) AddAttachmentToAchievement(ctx context.Context, id string, attachment modelmongo.Attachment) (*modelmongo.Achievement, error) {
	if m.addAttachmentErr != nil {
		return nil, m.addAttachmentErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockAchievementRepo) GetAchievementsByType(ctx context.Context) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byType, nil
}

func (m *mockAchievementRepo) GetCompetitionLevelDistribution(ctx context.Context) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.competitionLevelDist, nil
}

func (m *mockAchievementRepo) GetTopStudentsByPoints(ctx context.Context, limit int) ([]repositorymongo.TopStudentResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.topStudents, nil
}

type mockAchievementRefRepo struct {
	byMongoID       *modelpostgre.AchievementReference
	byID            *modelpostgre.AchievementReference
	byStudentID     []modelpostgre.AchievementReference
	byAdvisorID     []modelpostgre.AchievementReference
	allReferences   []modelpostgre.AchievementReference
	statsTotal      int
	statsVerified   int
	byPeriod        map[string]int
	allMongoIDs     []string
	err             error
	createErr       error
	updateErr       error
	updateVerifyErr error
	updateRejectErr error
}

func (m *mockAchievementRefRepo) CreateAchievementReference(ctx context.Context, req modelpostgre.CreateAchievementReferenceRequest) (*modelpostgre.AchievementReference, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	ref := &modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          req.StudentID,
		MongoAchievementID: req.MongoAchievementID,
		Status:             req.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	return ref, nil
}

func (m *mockAchievementRefRepo) GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*modelpostgre.AchievementReference, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byMongoID, nil
}

func (m *mockAchievementRefRepo) GetAchievementReferenceByID(ctx context.Context, id string) (*modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockAchievementRefRepo) UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, submittedAt *time.Time) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	return m.err
}

func (m *mockAchievementRefRepo) DeleteAchievementReference(ctx context.Context, id string) error {
	return m.err
}

func (m *mockAchievementRefRepo) GetAchievementReferenceByStudentID(ctx context.Context, studentID string) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byStudentID, nil
}

func (m *mockAchievementRefRepo) GetAchievementReferencesByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byAdvisorID, nil
}

func (m *mockAchievementRefRepo) GetAllAchievementReferences(ctx context.Context) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allReferences, nil
}

func (m *mockAchievementRefRepo) GetAchievementReferenceByStudentIDPaginated(ctx context.Context, studentID string, page, limit int) ([]modelpostgre.AchievementReference, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.byStudentID, len(m.byStudentID), nil
}

func (m *mockAchievementRefRepo) GetAchievementReferencesByAdvisorIDPaginated(ctx context.Context, advisorID string, page, limit int) ([]modelpostgre.AchievementReference, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.byAdvisorID, len(m.byAdvisorID), nil
}

func (m *mockAchievementRefRepo) GetAllAchievementReferencesPaginated(ctx context.Context, page, limit int, statusFilter string, sortBy string, sortOrder string) ([]modelpostgre.AchievementReference, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.allReferences, len(m.allReferences), nil
}

func (m *mockAchievementRefRepo) UpdateAchievementReferenceVerify(ctx context.Context, id string, verifiedBy string) error {
	if m.updateVerifyErr != nil {
		return m.updateVerifyErr
	}
	return m.err
}

func (m *mockAchievementRefRepo) UpdateAchievementReferenceReject(ctx context.Context, id string, verifiedBy string, rejectionNote string) error {
	if m.updateRejectErr != nil {
		return m.updateRejectErr
	}
	return m.err
}

func (m *mockAchievementRefRepo) GetAchievementStats(ctx context.Context) (int, int, error) {
	if m.err != nil {
		return 0, 0, m.err
	}
	return m.statsTotal, m.statsVerified, nil
}

func (m *mockAchievementRefRepo) GetAchievementsByPeriod(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byPeriod, nil
}

func (m *mockAchievementRefRepo) GetAllAchievementMongoIDs(ctx context.Context) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allMongoIDs, nil
}

type mockUserRepo struct {
	byID              *modelpostgre.User
	byEmail           *modelpostgre.User
	byUsernameOrEmail *modelpostgre.User
	allUsers          []modelpostgre.User
	allRoles          []modelpostgre.Role
	roleName          string
	permissions       []string
	lecturerByUserID  *modelpostgre.Lecturer
	lecturerByID      *modelpostgre.Lecturer
	err               error
}

func (m *mockUserRepo) FindUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockUserRepo) FindUserByEmail(ctx context.Context, email string) (*modelpostgre.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byEmail, nil
}

func (m *mockUserRepo) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byUsernameOrEmail, nil
}

func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allUsers, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user modelpostgre.User) (*modelpostgre.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &user, nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, id string, user modelpostgre.User) (*modelpostgre.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &user, nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.err
}

func (m *mockUserRepo) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return m.err
}

func (m *mockUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.permissions, nil
}

func (m *mockUserRepo) GetRoleName(ctx context.Context, roleID string) (string, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", m.err
	}
	return m.roleName, nil
}

func (m *mockUserRepo) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allRoles, nil
}

func (m *mockUserRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockUserRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

type mockStudentRepo struct {
	byID              *modelpostgre.Student
	byUserID          *modelpostgre.Student
	studentIDByUserID string
	byAdvisorID       []modelpostgre.Student
	allStudents       []modelpostgre.Student
	err               error
}

func (m *mockStudentRepo) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockStudentRepo) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byUserID, nil
}

func (m *mockStudentRepo) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentRepo) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byAdvisorID, nil
}

func (m *mockStudentRepo) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allStudents, nil
}

func (m *mockStudentRepo) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	student := &modelpostgre.Student{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
		CreatedAt:    time.Now(),
	}
	return student, nil
}

func (m *mockStudentRepo) UpdateStudent(ctx context.Context, id string, req modelpostgre.UpdateStudentRequest) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentRepo) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	return m.err
}

type mockNotificationService struct {
	err error
}

func (m *mockNotificationService) GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error) {
	return nil, m.err
}

func (m *mockNotificationService) GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error) {
	return nil, m.err
}

func (m *mockNotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error) {
	return nil, m.err
}

func (m *mockNotificationService) MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error) {
	return nil, m.err
}

func (m *mockNotificationService) CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error {
	return m.err
}

func (m *mockNotificationService) CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error {
	return m.err
}

func TestCreateAchievement_Success(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRepo := &mockAchievementRepo{
		byID: &modelmongo.Achievement{
			ID:              primitive.NewObjectID(),
			StudentID:       "550e8400-e29b-41d4-a716-446655440000",
			AchievementType: "academic",
			Title:           "Test Achievement",
			Description:     "Test Description",
			Points:          100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockAchievementRefRepo := &mockAchievementRefRepo{}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
		byID: &modelpostgre.Student{
			ID:     "550e8400-e29b-41d4-a716-446655440000",
			UserID: "user-id-1",
		},
	}

	mockNotificationService := &mockNotificationService{}

	service := servicepostgre.NewAchievementService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		mockNotificationService,
	)

	req := modelmongo.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	result, err := service.CreateAchievement(ctx, "user-id-1", "role-id-1", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestCreateAchievement_InvalidRole(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockUserRepo{
		roleName: "Admin",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		&mockAchievementRefRepo{},
		mockUserRepo,
		&mockStudentRepo{},
		&mockNotificationService{},
	)

	req := modelmongo.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "Test",
		Description:     "Test",
		Points:          100,
	}

	_, err := service.CreateAchievement(ctx, "user-id-1", "role-id-1", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "akses ditolak. Hanya mahasiswa yang dapat membuat prestasi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateAchievement_MissingStudentProfile(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		&mockAchievementRefRepo{},
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	req := modelmongo.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "Test",
		Description:     "Test",
		Points:          100,
	}

	_, err := service.CreateAchievement(ctx, "user-id-1", "role-id-1", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "data mahasiswa tidak ditemukan. Pastikan user memiliki profil mahasiswa" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateAchievement_ValidationErrors(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
		byID: &modelpostgre.Student{
			ID: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		&mockAchievementRefRepo{},
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	testCases := []struct {
		name string
		req  modelmongo.CreateAchievementRequest
		want string
	}{
		{
			name: "empty achievement type",
			req: modelmongo.CreateAchievementRequest{
				Title:       "Test",
				Description: "Test",
				Points:      100,
			},
			want: "achievement type wajib diisi",
		},
		{
			name: "empty title",
			req: modelmongo.CreateAchievementRequest{
				AchievementType: "academic",
				Description:     "Test",
				Points:          100,
			},
			want: "title wajib diisi",
		},
		{
			name: "empty description",
			req: modelmongo.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "Test",
				Points:          100,
			},
			want: "description wajib diisi",
		},
		{
			name: "invalid achievement type",
			req: modelmongo.CreateAchievementRequest{
				AchievementType: "invalid",
				Title:           "Test",
				Description:     "Test",
				Points:          100,
			},
			want: "achievement type tidak valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateAchievement(ctx, "user-id-1", "role-id-1", tc.req)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tc.want && !contains(err.Error(), tc.want) {
				t.Errorf("Expected error containing '%s', got: %v", tc.want, err)
			}
		})
	}
}

func TestSubmitAchievement_Success(t *testing.T) {
	ctx := setupTestContext()

	now := time.Now()
	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusDraft,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		byID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusSubmitted,
			SubmittedAt:        &now,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	result, err := service.SubmitAchievement(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data.Status != modelpostgre.AchievementStatusSubmitted {
		t.Errorf("Expected status 'submitted', got '%s'", result.Data.Status)
	}
}

func TestSubmitAchievement_WrongStatus(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusSubmitted,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	_, err := service.SubmitAchievement(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "prestasi hanya dapat di-submit jika status adalah draft" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestVerifyAchievement_Success(t *testing.T) {
	ctx := setupTestContext()

	now := time.Now()
	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusSubmitted,
		},
		byID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusVerified,
			VerifiedAt:         &now,
			VerifiedBy:         stringPtr("lecturer-id-1"),
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Dosen Wali",
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-1",
			UserID: "lecturer-user-id-1",
		},
	}

	mockStudentRepo := &mockStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			AdvisorID: "lecturer-id-1",
		},
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	result, err := service.VerifyAchievement(ctx, "lecturer-user-id-1", "role-id-1", "mongo-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestRejectAchievement_Success(t *testing.T) {
	ctx := setupTestContext()

	now := time.Now()
	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusSubmitted,
		},
		byID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusRejected,
			VerifiedBy:         stringPtr("lecturer-id-1"),
			RejectionNote:      stringPtr("Data tidak lengkap"),
			UpdatedAt:          now,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Dosen Wali",
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-1",
			UserID: "lecturer-user-id-1",
		},
	}

	mockStudentRepo := &mockStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			AdvisorID: "lecturer-id-1",
		},
	}

	mockNotificationService := &mockNotificationService{}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		mockNotificationService,
	)

	req := modelpostgre.RejectAchievementRequest{
		RejectionNote: "Data tidak lengkap",
	}

	result, err := service.RejectAchievement(ctx, "lecturer-user-id-1", "role-id-1", "mongo-id-1", req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestRejectAchievement_MissingRejectionNote(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockUserRepo{
		roleName: "Dosen Wali",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		&mockAchievementRefRepo{},
		mockUserRepo,
		&mockStudentRepo{},
		&mockNotificationService{},
	)

	req := modelpostgre.RejectAchievementRequest{
		RejectionNote: "",
	}

	_, err := service.RejectAchievement(ctx, "lecturer-user-id-1", "role-id-1", "mongo-id-1", req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "rejection note wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestDeleteAchievement_Success(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusDraft,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	result, err := service.DeleteAchievement(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestGetAchievementByID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRepo := &mockAchievementRepo{
		byID: &modelmongo.Achievement{
			ID:              primitive.NewObjectID(),
			StudentID:       "550e8400-e29b-41d4-a716-446655440000",
			AchievementType: "academic",
			Title:           "Test Achievement",
			Description:     "Test Description",
			Points:          100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusDraft,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
	}

	service := servicepostgre.NewAchievementService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	result, err := service.GetAchievementByID(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
}

func TestGetAchievementHistory_Success(t *testing.T) {
	ctx := setupTestContext()

	now := time.Now()
	submittedAt := now.Add(1 * time.Hour)
	verifiedAt := now.Add(2 * time.Hour)
	verifiedBy := "lecturer-id-1"

	mockAchievementRefRepo := &mockAchievementRefRepo{
		byMongoID: &modelpostgre.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "550e8400-e29b-41d4-a716-446655440000",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusVerified,
			SubmittedAt:        &submittedAt,
			VerifiedAt:         &verifiedAt,
			VerifiedBy:         &verifiedBy,
			CreatedAt:          now,
			UpdatedAt:          verifiedAt,
		},
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
		byID: &modelpostgre.User{
			ID:       "lecturer-id-1",
			FullName: "Dosen Test",
		},
	}

	mockStudentRepo := &mockStudentRepo{
		studentIDByUserID: "550e8400-e29b-41d4-a716-446655440000",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		mockStudentRepo,
		&mockNotificationService{},
	)

	result, err := service.GetAchievementHistory(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}

	history, ok := result["data"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be array of history items")
	}

	if len(history) < 2 {
		t.Errorf("Expected at least 2 history entries, got %d", len(history))
	}
}

func TestGetAchievementHistory_NotFound(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRefRepo := &mockAchievementRefRepo{
		err: sql.ErrNoRows,
	}

	mockUserRepo := &mockUserRepo{
		roleName: "Mahasiswa",
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		mockUserRepo,
		&mockStudentRepo{},
		&mockNotificationService{},
	)

	_, err := service.GetAchievementHistory(ctx, "user-id-1", "role-id-1", "mongo-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "prestasi tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetAchievementStats_Success(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRefRepo := &mockAchievementRefRepo{
		statsTotal:    100,
		statsVerified: 75,
	}

	service := servicepostgre.NewAchievementService(
		&mockAchievementRepo{},
		mockAchievementRefRepo,
		&mockUserRepo{},
		&mockStudentRepo{},
		&mockNotificationService{},
	)

	result, err := service.GetAchievementStats(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
}

func stringPtr(s string) *string {
	return &s
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
