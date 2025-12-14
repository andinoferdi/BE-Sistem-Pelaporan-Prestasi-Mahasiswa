package route_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/DATA-DOG/go-sqlmock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockAchievementService struct {
	createResponse      *modelmongo.CreateAchievementResponse
	createErr           error
	submitResponse      *modelpostgre.UpdateAchievementReferenceResponse
	submitErr           error
	verifyResponse      *modelpostgre.VerifyAchievementResponse
	verifyErr           error
	rejectResponse      *modelpostgre.RejectAchievementResponse
	rejectErr           error
	deleteResponse      *modelmongo.DeleteAchievementResponse
	deleteErr           error
	getAchievementsResp map[string]interface{}
	getAchievementsErr  error
	getByStudentIDResp  map[string]interface{}
	getByStudentIDErr   error
	getByIDResp         map[string]interface{}
	getByIDErr          error
	updateResp          map[string]interface{}
	updateErr           error
	statsResp           map[string]interface{}
	statsErr            error
	uploadFileResp      *modelmongo.Attachment
	uploadFileErr       error
	historyResp         map[string]interface{}
	historyErr          error
}

func (m *mockAchievementService) CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.createResponse, nil
}

func (m *mockAchievementService) SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error) {
	if m.submitErr != nil {
		return nil, m.submitErr
	}
	return m.submitResponse, nil
}

func (m *mockAchievementService) VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	return m.verifyResponse, nil
}

func (m *mockAchievementService) RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error) {
	if m.rejectErr != nil {
		return nil, m.rejectErr
	}
	return m.rejectResponse, nil
}

func (m *mockAchievementService) DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error) {
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	return m.deleteResponse, nil
}

func (m *mockAchievementService) GetAchievements(ctx context.Context, userID string, roleID string, page, limit int, statusFilter string, achievementTypeFilter string, sortBy string, sortOrder string) (map[string]interface{}, error) {
	if m.getAchievementsErr != nil {
		return nil, m.getAchievementsErr
	}
	return m.getAchievementsResp, nil
}

func (m *mockAchievementService) GetAchievementsByStudentID(ctx context.Context, studentID string, page, limit int) (map[string]interface{}, error) {
	if m.getByStudentIDErr != nil {
		return nil, m.getByStudentIDErr
	}
	return m.getByStudentIDResp, nil
}

func (m *mockAchievementService) GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.getByIDResp, nil
}

func (m *mockAchievementService) UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.updateResp, nil
}

func (m *mockAchievementService) GetAchievementStats(ctx context.Context) (map[string]interface{}, error) {
	if m.statsErr != nil {
		return nil, m.statsErr
	}
	return m.statsResp, nil
}

func (m *mockAchievementService) UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error) {
	if m.uploadFileErr != nil {
		return nil, m.uploadFileErr
	}
	return m.uploadFileResp, nil
}

func (m *mockAchievementService) GetAchievementHistory(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	if m.historyErr != nil {
		return nil, m.historyErr
	}
	return m.historyResp, nil
}

func TestGetAchievementStatsRoute_Success(t *testing.T) {
	app := setupTestApp()
	mockService := &mockAchievementService{
		statsResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"total":      100,
				"verified":   75,
				"percentage": 75,
			},
		},
	}

	routepostgre.AchievementRoutes(app, mockService, nil)

	req := httptest.NewRequest("GET", "/api/v1/achievements/stats", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetAchievementsRoute_Success(t *testing.T) {
	db, _ := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockAchievementService{
		getAchievementsResp: map[string]interface{}{
			"status": "success",
			"data":   []interface{}{},
			"pagination": map[string]interface{}{
				"page":        1,
				"limit":       10,
				"total":       0,
				"total_pages": 0,
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("GET", "/api/v1/achievements", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetAchievementByIDRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:read").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	mockService := &mockAchievementService{
		getByIDResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"id":    mongoID,
				"title": "Test Achievement",
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("GET", "/api/v1/achievements/"+mongoID, nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateAchievementRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:create").
		WillReturnRows(rows)

	mockService := &mockAchievementService{
		createResponse: &modelmongo.CreateAchievementResponse{
			Status: "success",
			Data: modelmongo.Achievement{
				ID:              primitive.NewObjectID(),
				StudentID:       "student-id-1",
				AchievementType: "academic",
				Title:           "Test Achievement",
				Description:     "Test Description",
				Points:          100,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	reqBody := modelmongo.CreateAchievementRequest{
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	req := createRequestWithToken("POST", "/api/v1/achievements", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateAchievementRoute_InvalidBody(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:create").
		WillReturnRows(rows)

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, &mockAchievementService{}, db)

	req := createRequestWithToken("POST", "/api/v1/achievements", "invalid json", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusBadRequest)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSubmitAchievementRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:update").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	now := time.Now()
	mockService := &mockAchievementService{
		submitResponse: &modelpostgre.UpdateAchievementReferenceResponse{
			Status: "success",
			Data: modelpostgre.AchievementReference{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: mongoID,
				Status:             modelpostgre.AchievementStatusSubmitted,
				SubmittedAt:        &now,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("POST", "/api/v1/achievements/"+mongoID+"/submit", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestVerifyAchievementRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:verify").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	now := time.Now()
	mockService := &mockAchievementService{
		verifyResponse: &modelpostgre.VerifyAchievementResponse{
			Status: "success",
			Data: modelpostgre.AchievementReference{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: mongoID,
				Status:             modelpostgre.AchievementStatusVerified,
				VerifiedAt:         &now,
				VerifiedBy:         stringPtr(userID),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("POST", "/api/v1/achievements/"+mongoID+"/verify", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRejectAchievementRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:verify").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	rejectionNote := "Data tidak lengkap"
	mockService := &mockAchievementService{
		rejectResponse: &modelpostgre.RejectAchievementResponse{
			Status: "success",
			Data: modelpostgre.AchievementReference{
				ID:                 "ref-id-1",
				StudentID:          "student-id-1",
				MongoAchievementID: mongoID,
				Status:             modelpostgre.AchievementStatusRejected,
				RejectionNote:      &rejectionNote,
				VerifiedBy:         stringPtr(userID),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	reqBody := modelpostgre.RejectAchievementRequest{
		RejectionNote: rejectionNote,
	}

	req := createRequestWithToken("POST", "/api/v1/achievements/"+mongoID+"/reject", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteAchievementRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:delete").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	mockService := &mockAchievementService{
		deleteResponse: &modelmongo.DeleteAchievementResponse{
			Status: "success",
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("DELETE", "/api/v1/achievements/"+mongoID, nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAchievementHistoryRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:read").
		WillReturnRows(rows)

	mongoID := primitive.NewObjectID().Hex()
	mockService := &mockAchievementService{
		historyResp: map[string]interface{}{
			"status": "success",
			"data": []map[string]interface{}{
				{
					"status":     "draft",
					"changed_at": time.Now().Format(time.RFC3339),
				},
			},
		},
	}

	app := setupTestApp()
	routepostgre.AchievementRoutes(app, mockService, db)

	req := createRequestWithToken("GET", "/api/v1/achievements/"+mongoID+"/history", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
