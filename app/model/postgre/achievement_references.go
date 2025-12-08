package model

import "time"

const (
	AchievementStatusDraft     = "draft"
	AchievementStatusSubmitted = "submitted"
	AchievementStatusVerified   = "verified"
	AchievementStatusRejected   = "rejected"
	AchievementStatusDeleted   = "deleted"
)

type AchievementReference struct {
	ID                 string     `json:"id"`
	StudentID           string     `json:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerifiedBy         *string    `json:"verified_by"`
	RejectionNote      *string    `json:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type CreateAchievementReferenceRequest struct {
	StudentID           string `json:"student_id" validate:"required"`
	MongoAchievementID  string `json:"mongo_achievement_id" validate:"required"`
	Status              string `json:"status" validate:"required"`
}

type UpdateAchievementReferenceRequest struct {
	Status         string `json:"status" validate:"required"`
	RejectionNote  string `json:"rejection_note"`
}

type VerifyAchievementRequest struct {
}

type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note" validate:"required"`
}

type GetAllAchievementReferencesResponse struct {
	Status string                 `json:"status"`
	Data   []AchievementReference `json:"data"`
}

type GetAchievementReferenceByIDResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

type CreateAchievementReferenceResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

type UpdateAchievementReferenceResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

type DeleteAchievementReferenceResponse struct {
	Status string `json:"status"`
}

type VerifyAchievementResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

type RejectAchievementResponse struct {
	Status string               `json:"status"`
	Data   AchievementReference `json:"data"`
}

