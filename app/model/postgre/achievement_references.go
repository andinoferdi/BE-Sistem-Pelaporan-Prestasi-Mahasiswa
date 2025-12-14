package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: definisikan konstanta status prestasi untuk workflow approval
const (
	AchievementStatusDraft     = "draft"
	AchievementStatusSubmitted = "submitted"
	AchievementStatusVerified  = "verified"
	AchievementStatusRejected  = "rejected"
	AchievementStatusDeleted   = "deleted"
)

// #3 proses: struct untuk menyimpan referensi prestasi di PostgreSQL, link ke MongoDB achievement
type AchievementReference struct {
	ID                 string     `json:"id"`
	StudentID          string     `json:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerifiedBy         *string    `json:"verified_by"`
	RejectionNote      *string    `json:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// #4 proses: struct untuk request create referensi prestasi baru
type CreateAchievementReferenceRequest struct {
	StudentID          string `json:"student_id" validate:"required"`
	MongoAchievementID string `json:"mongo_achievement_id" validate:"required"`
	Status             string `json:"status" validate:"required"`
}

// #5 proses: struct untuk request update status referensi prestasi
type UpdateAchievementReferenceRequest struct {
	Status        string `json:"status" validate:"required"`
	RejectionNote string `json:"rejection_note"`
}

// #6 proses: struct untuk request verify prestasi, tidak perlu field tambahan
type VerifyAchievementRequest struct {
}

// #7 proses: struct untuk request reject prestasi, wajib ada alasan penolakan
type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note" validate:"required"`
}

// #8 proses: struct response untuk get all achievement references, return list semua referensi
type GetAllAchievementReferencesResponse struct {
	Status string                 `json:"status"`
	Data   []AchievementReference `json:"data"`
}

// #9 proses: struct response untuk get achievement reference by ID, return satu referensi
type GetAchievementReferenceByIDResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

// #10 proses: struct response untuk create achievement reference, return referensi yang baru dibuat
type CreateAchievementReferenceResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

// #11 proses: struct response untuk update achievement reference, return referensi yang sudah diupdate
type UpdateAchievementReferenceResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

// #12 proses: struct response untuk delete achievement reference, hanya return status
type DeleteAchievementReferenceResponse struct {
	Status string `json:"status"`
}

// #13 proses: struct response untuk verify achievement, return referensi yang sudah diverifikasi
type VerifyAchievementResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

// #14 proses: struct response untuk reject achievement, return referensi yang sudah ditolak
type RejectAchievementResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}
