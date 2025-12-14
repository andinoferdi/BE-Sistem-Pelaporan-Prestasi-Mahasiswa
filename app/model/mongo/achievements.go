package model

// #1 proses: import library yang diperlukan untuk MongoDB dan time
import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// #2 proses: definisikan konstanta tipe prestasi yang bisa dipilih
const (
	AchievementTypeAcademic      = "academic"
	AchievementTypeCompetition   = "competition"
	AchievementTypeOrganization  = "organization"
	AchievementTypePublication   = "publication"
	AchievementTypeCertification = "certification"
	AchievementTypeOther         = "other"
)

// #3 proses: definisikan konstanta level kompetisi untuk prestasi tipe competition
const (
	CompetitionLevelInternational = "international"
	CompetitionLevelNational      = "national"
	CompetitionLevelRegional      = "regional"
	CompetitionLevelLocal         = "local"
)

// #4 proses: definisikan konstanta tipe publikasi untuk prestasi tipe publication
const (
	PublicationTypeJournal    = "journal"
	PublicationTypeConference = "conference"
	PublicationTypeBook       = "book"
)

// #5 proses: struct untuk periode waktu, dipakai untuk prestasi tipe organization
type Period struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

// #6 proses: struct untuk attachment file yang diupload untuk prestasi
type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

// #7 proses: struct untuk detail prestasi yang dinamis, field berbeda tergantung tipe prestasi
type AchievementDetails struct {
	CompetitionName     *string                `bson:"competitionName,omitempty" json:"competitionName,omitempty"`
	CompetitionLevel    *string                `bson:"competitionLevel,omitempty" json:"competitionLevel,omitempty"`
	Rank                *int                   `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType           *string                `bson:"medalType,omitempty" json:"medalType,omitempty"`
	PublicationType     *string                `bson:"publicationType,omitempty" json:"publicationType,omitempty"`
	PublicationTitle    *string                `bson:"publicationTitle,omitempty" json:"publicationTitle,omitempty"`
	Authors             []string               `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher           *string                `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN                *string                `bson:"issn,omitempty" json:"issn,omitempty"`
	OrganizationName    *string                `bson:"organizationName,omitempty" json:"organizationName,omitempty"`
	Position            *string                `bson:"position,omitempty" json:"position,omitempty"`
	Period              *Period                `bson:"period,omitempty" json:"period,omitempty"`
	CertificationName   *string                `bson:"certificationName,omitempty" json:"certificationName,omitempty"`
	IssuedBy            *string                `bson:"issuedBy,omitempty" json:"issuedBy,omitempty"`
	CertificationNumber *string                `bson:"certificationNumber,omitempty" json:"certificationNumber,omitempty"`
	ValidUntil          *time.Time             `bson:"validUntil,omitempty" json:"validUntil,omitempty"`
	EventDate           *time.Time             `bson:"eventDate,omitempty" json:"eventDate,omitempty"`
	Location            *string                `bson:"location,omitempty" json:"location,omitempty"`
	Organizer           *string                `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score               *float64               `bson:"score,omitempty" json:"score,omitempty"`
	CustomFields        map[string]interface{} `bson:"customFields,omitempty" json:"customFields,omitempty"`
}

// #8 proses: struct utama untuk menyimpan data prestasi di MongoDB
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

// #9 proses: struct untuk request create prestasi baru
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

// #10 proses: struct untuk request update prestasi, semua field optional karena partial update
type UpdateAchievementRequest struct {
	AchievementType string              `bson:"achievementType,omitempty" json:"achievementType,omitempty" validate:"omitempty,oneof=academic competition organization publication certification other"`
	Title           string              `bson:"title,omitempty" json:"title,omitempty" validate:"omitempty"`
	Description     string              `bson:"description,omitempty" json:"description,omitempty" validate:"omitempty"`
	Details         *AchievementDetails `bson:"details,omitempty" json:"details,omitempty"`
	Attachments     []Attachment        `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string            `bson:"tags,omitempty" json:"tags,omitempty"`
	Points          *int                `bson:"points,omitempty" json:"points,omitempty"`
}

// #11 proses: struct response untuk get all achievements, return list semua prestasi
type GetAllAchievementsResponse struct {
	Status string        `json:"status"`
	Data   []Achievement `json:"data"`
}

// #12 proses: struct response untuk get achievement by ID, return satu prestasi
type GetAchievementByIDResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

// #13 proses: struct response untuk create achievement, return prestasi yang baru dibuat
type CreateAchievementResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

// #14 proses: struct response untuk update achievement, return prestasi yang sudah diupdate
type UpdateAchievementResponse struct {
	Status string      `json:"status"`
	Data   Achievement `json:"data"`
}

// #15 proses: struct response untuk delete achievement, hanya return status
type DeleteAchievementResponse struct {
	Status string `json:"status"`
}
