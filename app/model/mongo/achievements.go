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
	FileName    string    `bson:"fileName" json:"fileName"`
	FileURL     string    `bson:"fileUrl" json:"fileUrl"`
	FileType    string    `bson:"fileType" json:"fileType"`
	UploadedAt  time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

type AchievementDetails struct {
	CompetitionName    *string                `bson:"competitionName,omitempty" json:"competitionName,omitempty"`
	CompetitionLevel   *string                `bson:"competitionLevel,omitempty" json:"competitionLevel,omitempty"`
	Rank               *int                   `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType          *string                `bson:"medalType,omitempty" json:"medalType,omitempty"`
	PublicationType    *string                `bson:"publicationType,omitempty" json:"publicationType,omitempty"`
	PublicationTitle   *string                `bson:"publicationTitle,omitempty" json:"publicationTitle,omitempty"`
	Authors            []string               `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher          *string                `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN               *string                `bson:"issn,omitempty" json:"issn,omitempty"`
	OrganizationName   *string                `bson:"organizationName,omitempty" json:"organizationName,omitempty"`
	Position           *string                `bson:"position,omitempty" json:"position,omitempty"`
	Period             *Period                `bson:"period,omitempty" json:"period,omitempty"`
	CertificationName  *string                `bson:"certificationName,omitempty" json:"certificationName,omitempty"`
	IssuedBy           *string                `bson:"issuedBy,omitempty" json:"issuedBy,omitempty"`
	CertificationNumber *string               `bson:"certificationNumber,omitempty" json:"certificationNumber,omitempty"`
	ValidUntil         *time.Time             `bson:"validUntil,omitempty" json:"validUntil,omitempty"`
	EventDate          *time.Time             `bson:"eventDate,omitempty" json:"eventDate,omitempty"`
	Location           *string                `bson:"location,omitempty" json:"location,omitempty"`
	Organizer          *string                `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score              *float64                `bson:"score,omitempty" json:"score,omitempty"`
	CustomFields       map[string]interface{} `bson:"customFields,omitempty" json:"customFields,omitempty"`
}

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID       string             `bson:"studentId" json:"studentId"`
	AchievementType string             `bson:"achievementType" json:"achievementType"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          int                `bson:"points" json:"points"`
	DeletedAt       *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateAchievementRequest struct {
	StudentID       string             `bson:"studentId" json:"studentId" validate:"required"`
	AchievementType string             `bson:"achievementType" json:"achievementType" validate:"required,oneof=academic competition organization publication certification other"`
	Title           string             `bson:"title" json:"title" validate:"required"`
	Description     string             `bson:"description" json:"description" validate:"required"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          int                `bson:"points" json:"points" validate:"required"`
}

type UpdateAchievementRequest struct {
	AchievementType string             `bson:"achievementType,omitempty" json:"achievementType,omitempty" validate:"omitempty,oneof=academic competition organization publication certification other"`
	Title           string             `bson:"title,omitempty" json:"title,omitempty" validate:"omitempty"`
	Description     string             `bson:"description,omitempty" json:"description,omitempty" validate:"omitempty"`
	Details         *AchievementDetails `bson:"details,omitempty" json:"details,omitempty"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          *int               `bson:"points,omitempty" json:"points,omitempty"`
}

type GetAllAchievementsResponse struct {
	Status string        `json:"status"`
	Data   []Achievement `json:"data"`
}

type GetAchievementByIDResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

type CreateAchievementResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

type UpdateAchievementResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

type DeleteAchievementResponse struct {
	Status string `json:"status"`
}

