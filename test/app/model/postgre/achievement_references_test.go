package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestAchievementStatusConstants(t *testing.T) {
	if modelpostgre.AchievementStatusDraft != "draft" {
		t.Errorf("Expected AchievementStatusDraft 'draft', got '%s'", modelpostgre.AchievementStatusDraft)
	}
	if modelpostgre.AchievementStatusSubmitted != "submitted" {
		t.Errorf("Expected AchievementStatusSubmitted 'submitted', got '%s'", modelpostgre.AchievementStatusSubmitted)
	}
	if modelpostgre.AchievementStatusVerified != "verified" {
		t.Errorf("Expected AchievementStatusVerified 'verified', got '%s'", modelpostgre.AchievementStatusVerified)
	}
	if modelpostgre.AchievementStatusRejected != "rejected" {
		t.Errorf("Expected AchievementStatusRejected 'rejected', got '%s'", modelpostgre.AchievementStatusRejected)
	}
	if modelpostgre.AchievementStatusDeleted != "deleted" {
		t.Errorf("Expected AchievementStatusDeleted 'deleted', got '%s'", modelpostgre.AchievementStatusDeleted)
	}
}

func TestAchievementReference_StructCreation(t *testing.T) {
	now := time.Now()
	verifiedBy := "lecturer-id-1"
	rejectionNote := "Data tidak lengkap"
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		SubmittedAt:        &now,
		VerifiedAt:         &now,
		VerifiedBy:         &verifiedBy,
		RejectionNote:      &rejectionNote,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if ref.ID != "ref-id-1" {
		t.Errorf("Expected ID 'ref-id-1', got '%s'", ref.ID)
	}
	if ref.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", ref.StudentID)
	}
	if ref.MongoAchievementID != "mongo-id-1" {
		t.Errorf("Expected MongoAchievementID 'mongo-id-1', got '%s'", ref.MongoAchievementID)
	}
	if ref.Status != modelpostgre.AchievementStatusDraft {
		t.Errorf("Expected Status '%s', got '%s'", modelpostgre.AchievementStatusDraft, ref.Status)
	}
	if ref.SubmittedAt == nil {
		t.Error("Expected SubmittedAt to be non-nil")
	}
	if ref.VerifiedBy == nil || *ref.VerifiedBy != "lecturer-id-1" {
		t.Errorf("Expected VerifiedBy 'lecturer-id-1', got '%v'", ref.VerifiedBy)
	}
}

func TestAchievementReference_JSONMarshalling(t *testing.T) {
	now := time.Now()
	verifiedBy := "lecturer-id-1"
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		SubmittedAt:        &now,
		VerifiedAt:         &now,
		VerifiedBy:         &verifiedBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	jsonData, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "ref-id-1" {
		t.Errorf("Expected id 'ref-id-1', got '%v'", result["id"])
	}
	if result["student_id"] != "student-id-1" {
		t.Errorf("Expected student_id 'student-id-1', got '%v'", result["student_id"])
	}
	if result["mongo_achievement_id"] != "mongo-id-1" {
		t.Errorf("Expected mongo_achievement_id 'mongo-id-1', got '%v'", result["mongo_achievement_id"])
	}
	if result["status"] != "draft" {
		t.Errorf("Expected status 'draft', got '%v'", result["status"])
	}
}

func TestAchievementReference_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "ref-id-1",
		"student_id": "student-id-1",
		"mongo_achievement_id": "mongo-id-1",
		"status": "draft",
		"submitted_at": "2024-01-01T00:00:00Z",
		"verified_at": "2024-01-01T01:00:00Z",
		"verified_by": "lecturer-id-1",
		"rejection_note": "Data tidak lengkap",
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var ref modelpostgre.AchievementReference
	if err := json.Unmarshal([]byte(jsonStr), &ref); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ref.ID != "ref-id-1" {
		t.Errorf("Expected ID 'ref-id-1', got '%s'", ref.ID)
	}
	if ref.Status != "draft" {
		t.Errorf("Expected Status 'draft', got '%s'", ref.Status)
	}
	if ref.VerifiedBy == nil || *ref.VerifiedBy != "lecturer-id-1" {
		t.Errorf("Expected VerifiedBy 'lecturer-id-1', got '%v'", ref.VerifiedBy)
	}
}

func TestAchievementReference_ZeroValues(t *testing.T) {
	var ref modelpostgre.AchievementReference

	if ref.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", ref.ID)
	}
	if ref.StudentID != "" {
		t.Errorf("Expected empty StudentID, got '%s'", ref.StudentID)
	}
	if ref.Status != "" {
		t.Errorf("Expected empty Status, got '%s'", ref.Status)
	}
	if ref.SubmittedAt != nil {
		t.Error("Expected SubmittedAt to be nil")
	}
	if ref.VerifiedAt != nil {
		t.Error("Expected VerifiedAt to be nil")
	}
	if ref.VerifiedBy != nil {
		t.Error("Expected VerifiedBy to be nil")
	}
	if ref.RejectionNote != nil {
		t.Error("Expected RejectionNote to be nil")
	}
	if !ref.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", ref.CreatedAt)
	}
}

func TestAchievementReference_PointerFieldsNil(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if ref.SubmittedAt != nil {
		t.Error("Expected SubmittedAt to be nil")
	}
	if ref.VerifiedAt != nil {
		t.Error("Expected VerifiedAt to be nil")
	}
	if ref.VerifiedBy != nil {
		t.Error("Expected VerifiedBy to be nil")
	}
	if ref.RejectionNote != nil {
		t.Error("Expected RejectionNote to be nil")
	}
}

func TestAchievementReference_JSONMarshallingNilPointers(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	jsonData, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if submittedAt, exists := result["submitted_at"]; exists {
		if submittedAt != nil {
			t.Errorf("Expected submitted_at to be null when nil, got '%v'", submittedAt)
		}
	} else {
		t.Error("Expected submitted_at field to exist in JSON")
	}
	if verifiedAt, exists := result["verified_at"]; exists {
		if verifiedAt != nil {
			t.Errorf("Expected verified_at to be null when nil, got '%v'", verifiedAt)
		}
	} else {
		t.Error("Expected verified_at field to exist in JSON")
	}
	if verifiedBy, exists := result["verified_by"]; exists {
		if verifiedBy != nil {
			t.Errorf("Expected verified_by to be null when nil, got '%v'", verifiedBy)
		}
	} else {
		t.Error("Expected verified_by field to exist in JSON")
	}
	if rejectionNote, exists := result["rejection_note"]; exists {
		if rejectionNote != nil {
			t.Errorf("Expected rejection_note to be null when nil, got '%v'", rejectionNote)
		}
	} else {
		t.Error("Expected rejection_note field to exist in JSON")
	}
}

func TestCreateAchievementReferenceRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
	}

	if req.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", req.StudentID)
	}
	if req.MongoAchievementID != "mongo-id-1" {
		t.Errorf("Expected MongoAchievementID 'mongo-id-1', got '%s'", req.MongoAchievementID)
	}
	if req.Status != modelpostgre.AchievementStatusDraft {
		t.Errorf("Expected Status '%s', got '%s'", modelpostgre.AchievementStatusDraft, req.Status)
	}
}

func TestCreateAchievementReferenceRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["student_id"] != "student-id-1" {
		t.Errorf("Expected student_id 'student-id-1', got '%v'", result["student_id"])
	}
	if result["mongo_achievement_id"] != "mongo-id-1" {
		t.Errorf("Expected mongo_achievement_id 'mongo-id-1', got '%v'", result["mongo_achievement_id"])
	}
	if result["status"] != "draft" {
		t.Errorf("Expected status 'draft', got '%v'", result["status"])
	}
}

func TestCreateAchievementReferenceRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"student_id": "student-id-1",
		"mongo_achievement_id": "mongo-id-1",
		"status": "draft"
	}`

	var req modelpostgre.CreateAchievementReferenceRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", req.StudentID)
	}
	if req.Status != "draft" {
		t.Errorf("Expected Status 'draft', got '%s'", req.Status)
	}
}

func TestUpdateAchievementReferenceRequest_StructCreation(t *testing.T) {
	req := modelpostgre.UpdateAchievementReferenceRequest{
		Status:        modelpostgre.AchievementStatusSubmitted,
		RejectionNote: "Data tidak lengkap",
	}

	if req.Status != modelpostgre.AchievementStatusSubmitted {
		t.Errorf("Expected Status '%s', got '%s'", modelpostgre.AchievementStatusSubmitted, req.Status)
	}
	if req.RejectionNote != "Data tidak lengkap" {
		t.Errorf("Expected RejectionNote 'Data tidak lengkap', got '%s'", req.RejectionNote)
	}
}

func TestUpdateAchievementReferenceRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.UpdateAchievementReferenceRequest{
		Status:        modelpostgre.AchievementStatusSubmitted,
		RejectionNote: "Data tidak lengkap",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "submitted" {
		t.Errorf("Expected status 'submitted', got '%v'", result["status"])
	}
	if result["rejection_note"] != "Data tidak lengkap" {
		t.Errorf("Expected rejection_note 'Data tidak lengkap', got '%v'", result["rejection_note"])
	}
}

func TestVerifyAchievementRequest_StructCreation(t *testing.T) {
	req := modelpostgre.VerifyAchievementRequest{}

	if req == (modelpostgre.VerifyAchievementRequest{}) {
		// Empty struct is valid
	}
}

func TestVerifyAchievementRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.VerifyAchievementRequest{}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}

func TestRejectAchievementRequest_StructCreation(t *testing.T) {
	req := modelpostgre.RejectAchievementRequest{
		RejectionNote: "Data tidak lengkap",
	}

	if req.RejectionNote != "Data tidak lengkap" {
		t.Errorf("Expected RejectionNote 'Data tidak lengkap', got '%s'", req.RejectionNote)
	}
}

func TestRejectAchievementRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.RejectAchievementRequest{
		RejectionNote: "Data tidak lengkap",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["rejection_note"] != "Data tidak lengkap" {
		t.Errorf("Expected rejection_note 'Data tidak lengkap', got '%v'", result["rejection_note"])
	}
}

func TestRejectAchievementRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"rejection_note": "Data tidak lengkap"
	}`

	var req modelpostgre.RejectAchievementRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.RejectionNote != "Data tidak lengkap" {
		t.Errorf("Expected RejectionNote 'Data tidak lengkap', got '%s'", req.RejectionNote)
	}
}

func TestGetAllAchievementReferencesResponse_StructCreation(t *testing.T) {
	now := time.Now()
	refs := []modelpostgre.AchievementReference{
		{
			ID:                 "ref-id-1",
			StudentID:          "student-id-1",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusDraft,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 "ref-id-2",
			StudentID:          "student-id-2",
			MongoAchievementID: "mongo-id-2",
			Status:             modelpostgre.AchievementStatusSubmitted,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	resp := modelpostgre.GetAllAchievementReferencesResponse{
		Status: "success",
		Data:   refs,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 achievement references, got %d", len(resp.Data))
	}
	if resp.Data[0].Status != modelpostgre.AchievementStatusDraft {
		t.Errorf("Expected first ref Status '%s', got '%s'", modelpostgre.AchievementStatusDraft, resp.Data[0].Status)
	}
}

func TestGetAllAchievementReferencesResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	refs := []modelpostgre.AchievementReference{
		{
			ID:                 "ref-id-1",
			StudentID:          "student-id-1",
			MongoAchievementID: "mongo-id-1",
			Status:             modelpostgre.AchievementStatusDraft,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	resp := modelpostgre.GetAllAchievementReferencesResponse{
		Status: "success",
		Data:   refs,
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].([]interface{}); ok {
		if len(data) != 1 {
			t.Errorf("Expected 1 achievement reference in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetAchievementReferenceByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	resp := modelpostgre.GetAchievementReferenceByIDResponse{
		Status: "success",
		Data:   ref,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "ref-id-1" {
		t.Errorf("Expected Data.ID 'ref-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreateAchievementReferenceResponse_StructCreation(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	resp := modelpostgre.CreateAchievementReferenceResponse{
		Status: "success",
		Data:   ref,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Status != modelpostgre.AchievementStatusDraft {
		t.Errorf("Expected Data.Status '%s', got '%s'", modelpostgre.AchievementStatusDraft, resp.Data.Status)
	}
}

func TestUpdateAchievementReferenceResponse_StructCreation(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusSubmitted,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	resp := modelpostgre.UpdateAchievementReferenceResponse{
		Status: "success",
		Data:   ref,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Status != modelpostgre.AchievementStatusSubmitted {
		t.Errorf("Expected Data.Status '%s', got '%s'", modelpostgre.AchievementStatusSubmitted, resp.Data.Status)
	}
}

func TestDeleteAchievementReferenceResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteAchievementReferenceResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestVerifyAchievementResponse_StructCreation(t *testing.T) {
	now := time.Now()
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusVerified,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	resp := modelpostgre.VerifyAchievementResponse{
		Status: "success",
		Data:   ref,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Status != modelpostgre.AchievementStatusVerified {
		t.Errorf("Expected Data.Status '%s', got '%s'", modelpostgre.AchievementStatusVerified, resp.Data.Status)
	}
}

func TestRejectAchievementResponse_StructCreation(t *testing.T) {
	now := time.Now()
	verifiedBy := "lecturer-id-1"
	rejectionNote := "Data tidak lengkap"
	ref := modelpostgre.AchievementReference{
		ID:                 "ref-id-1",
		StudentID:          "student-id-1",
		MongoAchievementID: "mongo-id-1",
		Status:             modelpostgre.AchievementStatusRejected,
		VerifiedBy:         &verifiedBy,
		RejectionNote:      &rejectionNote,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	resp := modelpostgre.RejectAchievementResponse{
		Status: "success",
		Data:   ref,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Status != modelpostgre.AchievementStatusRejected {
		t.Errorf("Expected Data.Status '%s', got '%s'", modelpostgre.AchievementStatusRejected, resp.Data.Status)
	}
	if resp.Data.RejectionNote == nil || *resp.Data.RejectionNote != "Data tidak lengkap" {
		t.Errorf("Expected Data.RejectionNote 'Data tidak lengkap', got '%v'", resp.Data.RejectionNote)
	}
}
