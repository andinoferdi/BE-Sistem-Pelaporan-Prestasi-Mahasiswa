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

type mockReportServiceAchievementRepo struct {
	byStudentID          []modelmongo.Achievement
	byIDs                []modelmongo.Achievement
	byType               map[string]int
	competitionLevelDist map[string]int
	topStudents          []repositorymongo.TopStudentResult
	err                  error
}

func (m *mockReportServiceAchievementRepo) CreateAchievement(ctx context.Context, achievement *modelmongo.Achievement) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRepo) GetAchievementByID(ctx context.Context, id string) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRepo) UpdateAchievement(ctx context.Context, id string, req modelmongo.UpdateAchievementRequest) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRepo) DeleteAchievement(ctx context.Context, id string) error {
	return m.err
}

func (m *mockReportServiceAchievementRepo) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byStudentID, nil
}

func (m *mockReportServiceAchievementRepo) GetAchievementsByIDs(ctx context.Context, ids []string) ([]modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byIDs, nil
}

func (m *mockReportServiceAchievementRepo) AddAttachmentToAchievement(ctx context.Context, id string, attachment modelmongo.Attachment) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRepo) GetAchievementsByType(ctx context.Context) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byType, nil
}

func (m *mockReportServiceAchievementRepo) GetCompetitionLevelDistribution(ctx context.Context) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.competitionLevelDist, nil
}

func (m *mockReportServiceAchievementRepo) GetTopStudentsByPoints(ctx context.Context, limit int) ([]repositorymongo.TopStudentResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.topStudents, nil
}

type mockReportServiceAchievementRefRepo struct {
	byStudentID   []modelpostgre.AchievementReference
	byAdvisorID   []modelpostgre.AchievementReference
	allReferences []modelpostgre.AchievementReference
	statsTotal    int
	statsVerified int
	byPeriod      map[string]int
	allMongoIDs   []string
	err           error
}

func (m *mockReportServiceAchievementRefRepo) CreateAchievementReference(ctx context.Context, req modelpostgre.CreateAchievementReferenceRequest) (*modelpostgre.AchievementReference, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*modelpostgre.AchievementReference, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferenceByID(ctx context.Context, id string) (*modelpostgre.AchievementReference, error) {
	return nil, m.err
}

func (m *mockReportServiceAchievementRefRepo) UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, submittedAt *time.Time) error {
	return m.err
}

func (m *mockReportServiceAchievementRefRepo) DeleteAchievementReference(ctx context.Context, id string) error {
	return m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferenceByStudentID(ctx context.Context, studentID string) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byStudentID, nil
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferencesByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byAdvisorID, nil
}

func (m *mockReportServiceAchievementRefRepo) GetAllAchievementReferences(ctx context.Context) ([]modelpostgre.AchievementReference, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allReferences, nil
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferenceByStudentIDPaginated(ctx context.Context, studentID string, page, limit int) ([]modelpostgre.AchievementReference, int, error) {
	return nil, 0, m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementReferencesByAdvisorIDPaginated(ctx context.Context, advisorID string, page, limit int) ([]modelpostgre.AchievementReference, int, error) {
	return nil, 0, m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAllAchievementReferencesPaginated(ctx context.Context, page, limit int, statusFilter string, sortBy string, sortOrder string) ([]modelpostgre.AchievementReference, int, error) {
	return nil, 0, m.err
}

func (m *mockReportServiceAchievementRefRepo) UpdateAchievementReferenceVerify(ctx context.Context, id string, verifiedBy string) error {
	return m.err
}

func (m *mockReportServiceAchievementRefRepo) UpdateAchievementReferenceReject(ctx context.Context, id string, verifiedBy string, rejectionNote string) error {
	return m.err
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementStats(ctx context.Context) (int, int, error) {
	if m.err != nil {
		return 0, 0, m.err
	}
	return m.statsTotal, m.statsVerified, nil
}

func (m *mockReportServiceAchievementRefRepo) GetAchievementsByPeriod(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byPeriod, nil
}

func (m *mockReportServiceAchievementRefRepo) GetAllAchievementMongoIDs(ctx context.Context) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allMongoIDs, nil
}

type mockReportServiceStudentRepo struct {
	byID        *modelpostgre.Student
	byUserID    *modelpostgre.Student
	byAdvisorID []modelpostgre.Student
	err         error
}

func (m *mockReportServiceStudentRepo) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	return "", m.err
}

func (m *mockReportServiceStudentRepo) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byUserID, nil
}

func (m *mockReportServiceStudentRepo) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockReportServiceStudentRepo) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byAdvisorID, nil
}

func (m *mockReportServiceStudentRepo) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockReportServiceStudentRepo) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockReportServiceStudentRepo) UpdateStudent(ctx context.Context, id string, req modelpostgre.UpdateStudentRequest) (*modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockReportServiceStudentRepo) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	return m.err
}

type mockReportServiceUserRepo struct {
	byID             *modelpostgre.User
	byIDMap          map[string]*modelpostgre.User
	roleName         string
	lecturerByUserID *modelpostgre.Lecturer
	lecturerByID     *modelpostgre.Lecturer
	err              error
}

func (m *mockReportServiceUserRepo) FindUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	if m.byIDMap != nil {
		if user, exists := m.byIDMap[id]; exists {
			return user, nil
		}
	}
	return m.byID, nil
}

func (m *mockReportServiceUserRepo) FindUserByEmail(ctx context.Context, email string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) CreateUser(ctx context.Context, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) UpdateUser(ctx context.Context, id string, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.err
}

func (m *mockReportServiceUserRepo) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return m.err
}

func (m *mockReportServiceUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) GetRoleName(ctx context.Context, roleID string) (string, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", m.err
	}
	return m.roleName, nil
}

func (m *mockReportServiceUserRepo) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	return nil, m.err
}

func (m *mockReportServiceUserRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockReportServiceUserRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

type mockReportServiceLecturerRepo struct {
	byID *modelpostgre.Lecturer
	err  error
}

func (m *mockReportServiceLecturerRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockReportServiceLecturerRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockReportServiceLecturerRepo) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockReportServiceLecturerRepo) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockReportServiceLecturerRepo) UpdateLecturer(ctx context.Context, id string, req modelpostgre.UpdateLecturerRequest) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func TestGetStatistics_Admin(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byType: map[string]int{
			"academic":     10,
			"competition":  5,
			"organization": 3,
		},
		competitionLevelDist: map[string]int{
			"international": 2,
			"national":      3,
		},
		topStudents: []repositorymongo.TopStudentResult{
			{
				StudentID:        "student-id-1",
				TotalPoints:      100,
				AchievementCount: 5,
			},
		},
	}

	mockAchievementRefRepo := &mockReportServiceAchievementRefRepo{
		statsTotal:    100,
		statsVerified: 75,
	}

	mockStudentRepo := &mockReportServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:     "student-id-1",
			UserID: "user-id-1",
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		roleName: "Admin",
		byID: &modelpostgre.User{
			ID:       "user-id-1",
			Username: "testuser",
			FullName: "Test User",
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockStudentRepo,
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	result, err := service.GetStatistics(ctx, "user-id-1", "role-id-1")

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

func TestGetStatistics_Mahasiswa(t *testing.T) {
	ctx := setupTestContext()

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              primitive.NewObjectID(),
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockReportServiceStudentRepo{
		byUserID: &modelpostgre.Student{
			ID:     "student-id-1",
			UserID: "user-id-1",
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		&mockReportServiceAchievementRefRepo{},
		mockStudentRepo,
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	result, err := service.GetStatistics(ctx, "user-id-1", "role-id-1")

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

func TestGetStatistics_DosenWali(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockReportServiceUserRepo{
		roleName: "Dosen Wali",
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-1",
			UserID: "user-id-1",
		},
	}

	mockStudentRepo := &mockReportServiceStudentRepo{
		byAdvisorID: []modelpostgre.Student{
			{
				ID:        "student-id-1",
				AdvisorID: "lecturer-id-1",
			},
		},
	}

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              primitive.NewObjectID(),
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		&mockReportServiceAchievementRefRepo{},
		mockStudentRepo,
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	result, err := service.GetStatistics(ctx, "user-id-1", "role-id-1")

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

func TestGetStatistics_InvalidRole(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockReportServiceUserRepo{
		roleName: "InvalidRole",
	}

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetStatistics(ctx, "user-id-1", "role-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "akses ditolak. Role tidak memiliki akses untuk melihat statistik" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetStudentReport_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockReportServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2023",
		},
	}

	achievementID := primitive.NewObjectID()
	achievementIDHex := achievementID.Hex()

	mockAchievementRefRepo := &mockReportServiceAchievementRefRepo{
		byStudentID: []modelpostgre.AchievementReference{
			{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: achievementIDHex,
				Status:             modelpostgre.AchievementStatusVerified,
				CreatedAt:          time.Now(),
			},
		},
	}

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              achievementID,
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		byID: &modelpostgre.User{
			ID:       "user-id-1",
			FullName: "Test User",
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockStudentRepo,
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	result, err := service.GetStudentReport(ctx, "student-id-1")

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

func TestGetStudentReport_EmptyStudentID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetStudentReport(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetStudentReport_NotFound(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockReportServiceStudentRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		mockStudentRepo,
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetStudentReport(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetLecturerReport_Success(t *testing.T) {
	ctx := setupTestContext()

	mockLecturerRepo := &mockReportServiceLecturerRepo{
		byID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
		},
	}

	achievementID := primitive.NewObjectID()
	achievementIDHex := achievementID.Hex()

	mockAchievementRefRepo := &mockReportServiceAchievementRefRepo{
		byStudentID: []modelpostgre.AchievementReference{
			{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: achievementIDHex,
				Status:             modelpostgre.AchievementStatusVerified,
				CreatedAt:          time.Now(),
			},
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		byIDMap: map[string]*modelpostgre.User{
			"user-id-1": {
				ID:       "user-id-1",
				FullName: "Test Lecturer",
			},
			"student-user-id-1": {
				ID:       "student-user-id-1",
				FullName: "Test Student",
			},
		},
	}

	mockStudentRepo := &mockReportServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "student-id-1",
			UserID:    "student-user-id-1",
			AdvisorID: "lecturer-id-1",
		},
		byAdvisorID: []modelpostgre.Student{
			{
				ID:        "student-id-1",
				UserID:    "student-user-id-1",
				AdvisorID: "lecturer-id-1",
			},
		},
	}

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              achievementID,
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockStudentRepo,
		mockUserRepo,
		mockLecturerRepo,
	)

	result, err := service.GetLecturerReport(ctx, "lecturer-id-1")

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

func TestGetLecturerReport_EmptyLecturerID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetLecturerReport(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "lecturer ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetCurrentStudentReport_Success(t *testing.T) {
	ctx := setupTestContext()

	achievementID := primitive.NewObjectID()
	achievementIDHex := achievementID.Hex()

	mockStudentRepo := &mockReportServiceStudentRepo{
		byUserID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2023",
		},
		byID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2023",
		},
	}

	mockAchievementRefRepo := &mockReportServiceAchievementRefRepo{
		byStudentID: []modelpostgre.AchievementReference{
			{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: achievementIDHex,
				Status:             modelpostgre.AchievementStatusVerified,
				CreatedAt:          time.Now(),
			},
		},
	}

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              achievementID,
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		byID: &modelpostgre.User{
			ID:       "user-id-1",
			FullName: "Test Student",
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockStudentRepo,
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	result, err := service.GetCurrentStudentReport(ctx, "user-id-1")

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

func TestGetCurrentStudentReport_EmptyUserID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetCurrentStudentReport(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetCurrentStudentReport_NoStudentProfile(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockReportServiceStudentRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		mockStudentRepo,
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetCurrentStudentReport(ctx, "user-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student profile tidak ditemukan untuk user ini" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetCurrentLecturerReport_Success(t *testing.T) {
	ctx := setupTestContext()

	achievementID := primitive.NewObjectID()
	achievementIDHex := achievementID.Hex()

	mockLecturerRepo := &mockReportServiceLecturerRepo{
		byID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
		},
	}

	mockUserRepo := &mockReportServiceUserRepo{
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
		},
		byIDMap: map[string]*modelpostgre.User{
			"user-id-1": {
				ID:       "user-id-1",
				FullName: "Test Lecturer",
			},
			"student-user-id-1": {
				ID:       "student-user-id-1",
				FullName: "Test Student",
			},
		},
	}

	mockStudentRepo := &mockReportServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "student-id-1",
			UserID:    "student-user-id-1",
			AdvisorID: "lecturer-id-1",
		},
		byAdvisorID: []modelpostgre.Student{
			{
				ID:        "student-id-1",
				UserID:    "student-user-id-1",
				AdvisorID: "lecturer-id-1",
			},
		},
	}

	mockAchievementRefRepo := &mockReportServiceAchievementRefRepo{
		byStudentID: []modelpostgre.AchievementReference{
			{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: achievementIDHex,
				Status:             modelpostgre.AchievementStatusVerified,
				CreatedAt:          time.Now(),
			},
		},
	}

	mockAchievementRepo := &mockReportServiceAchievementRepo{
		byStudentID: []modelmongo.Achievement{
			{
				ID:              achievementID,
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	service := servicepostgre.NewReportService(
		mockAchievementRepo,
		mockAchievementRefRepo,
		mockStudentRepo,
		mockUserRepo,
		mockLecturerRepo,
	)

	result, err := service.GetCurrentLecturerReport(ctx, "user-id-1")

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

func TestGetCurrentLecturerReport_EmptyUserID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		&mockReportServiceUserRepo{},
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetCurrentLecturerReport(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetCurrentLecturerReport_NoLecturerProfile(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockReportServiceUserRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewReportService(
		&mockReportServiceAchievementRepo{},
		&mockReportServiceAchievementRefRepo{},
		&mockReportServiceStudentRepo{},
		mockUserRepo,
		&mockReportServiceLecturerRepo{},
	)

	_, err := service.GetCurrentLecturerReport(ctx, "user-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "lecturer profile tidak ditemukan untuk user ini" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
