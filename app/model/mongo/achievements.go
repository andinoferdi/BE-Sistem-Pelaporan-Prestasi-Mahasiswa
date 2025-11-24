package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AchievementTypeAcademic      = "academic"
	AchievementTypeCompetition   = "competition"
	AchievementTypeOrganization  = "organization"
	AchievementTypePublication   = "publication"
	AchievementTypeCertification = "certification"
	AchievementTypeOther         = "other"
)

const (
	CompetitionLevelInternational = "international"
	CompetitionLevelNational       = "national"
	CompetitionLevelRegional       = "regional"
	CompetitionLevelLocal          = "local"
)

const (
	PublicationTypeJournal    = "journal"
	PublicationTypeConference = "conference"
	PublicationTypeBook       = "book"
)

type Period struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

type Attachment struct {
	FileName    string    `bson:"file_name" json:"file_name"`
	FileURL     string    `bson:"file_url" json:"file_url"`
	FileType    string    `bson:"file_type" json:"file_type"`
	UploadedAt  time.Time `bson:"uploaded_at" json:"uploaded_at"`
}

type AchievementDetails struct {
	CompetitionName    *string                `bson:"competition_name,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel   *string                `bson:"competition_level,omitempty" json:"competition_level,omitempty"`
	Rank               *int                   `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType          *string                `bson:"medal_type,omitempty" json:"medal_type,omitempty"`
	PublicationType    *string                `bson:"publication_type,omitempty" json:"publication_type,omitempty"`
	PublicationTitle   *string                `bson:"publication_title,omitempty" json:"publication_title,omitempty"`
	Authors            []string               `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher          *string                `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN               *string                `bson:"issn,omitempty" json:"issn,omitempty"`
	OrganizationName   *string                `bson:"organization_name,omitempty" json:"organization_name,omitempty"`
	Position           *string                `bson:"position,omitempty" json:"position,omitempty"`
	Period             *Period                `bson:"period,omitempty" json:"period,omitempty"`
	CertificationName  *string                `bson:"certification_name,omitempty" json:"certification_name,omitempty"`
	IssuedBy           *string                `bson:"issued_by,omitempty" json:"issued_by,omitempty"`
	CertificationNumber *string               `bson:"certification_number,omitempty" json:"certification_number,omitempty"`
	ValidUntil         *time.Time             `bson:"valid_until,omitempty" json:"valid_until,omitempty"`
	EventDate          *time.Time             `bson:"event_date,omitempty" json:"event_date,omitempty"`
	Location           *string                `bson:"location,omitempty" json:"location,omitempty"`
	Organizer          *string                `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score              *float64                `bson:"score,omitempty" json:"score,omitempty"`
	CustomFields       map[string]interface{} `bson:"custom_fields,omitempty" json:"custom_fields,omitempty"`
}

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID       string             `bson:"student_id" json:"student_id"`
	AchievementType string             `bson:"achievement_type" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          int                `bson:"points" json:"points"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateAchievementRequest struct {
	StudentID       string             `bson:"student_id" json:"student_id" validate:"required"`
	AchievementType string             `bson:"achievement_type" json:"achievement_type" validate:"required,oneof=academic competition organization publication certification other"`
	Title           string             `bson:"title" json:"title" validate:"required"`
	Description     string             `bson:"description" json:"description" validate:"required"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          int                `bson:"points" json:"points" validate:"required"`
}

type UpdateAchievementRequest struct {
	AchievementType string             `bson:"achievement_type,omitempty" json:"achievement_type,omitempty" validate:"omitempty,oneof=academic competition organization publication certification other"`
	Title           string             `bson:"title,omitempty" json:"title,omitempty" validate:"omitempty"`
	Description     string             `bson:"description,omitempty" json:"description,omitempty" validate:"omitempty"`
	Details         *AchievementDetails `bson:"details,omitempty" json:"details,omitempty"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          *int               `bson:"points,omitempty" json:"points,omitempty"`
}

type GetAllAchievementsResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []Achievement `json:"data"`
}

type GetAchievementByIDResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Achievement `json:"data"`
}

type CreateAchievementResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Achievement `json:"data"`
}

type UpdateAchievementResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    Achievement `json:"data"`
}

type DeleteAchievementResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

