package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAchievementTypeConstants(t *testing.T) {
	if modelmongo.AchievementTypeAcademic != "academic" {
		t.Errorf("Expected AchievementTypeAcademic 'academic', got '%s'", modelmongo.AchievementTypeAcademic)
	}
	if modelmongo.AchievementTypeCompetition != "competition" {
		t.Errorf("Expected AchievementTypeCompetition 'competition', got '%s'", modelmongo.AchievementTypeCompetition)
	}
	if modelmongo.AchievementTypeOrganization != "organization" {
		t.Errorf("Expected AchievementTypeOrganization 'organization', got '%s'", modelmongo.AchievementTypeOrganization)
	}
	if modelmongo.AchievementTypePublication != "publication" {
		t.Errorf("Expected AchievementTypePublication 'publication', got '%s'", modelmongo.AchievementTypePublication)
	}
	if modelmongo.AchievementTypeCertification != "certification" {
		t.Errorf("Expected AchievementTypeCertification 'certification', got '%s'", modelmongo.AchievementTypeCertification)
	}
	if modelmongo.AchievementTypeOther != "other" {
		t.Errorf("Expected AchievementTypeOther 'other', got '%s'", modelmongo.AchievementTypeOther)
	}
}

func TestCompetitionLevelConstants(t *testing.T) {
	if modelmongo.CompetitionLevelInternational != "international" {
		t.Errorf("Expected CompetitionLevelInternational 'international', got '%s'", modelmongo.CompetitionLevelInternational)
	}
	if modelmongo.CompetitionLevelNational != "national" {
		t.Errorf("Expected CompetitionLevelNational 'national', got '%s'", modelmongo.CompetitionLevelNational)
	}
	if modelmongo.CompetitionLevelRegional != "regional" {
		t.Errorf("Expected CompetitionLevelRegional 'regional', got '%s'", modelmongo.CompetitionLevelRegional)
	}
	if modelmongo.CompetitionLevelLocal != "local" {
		t.Errorf("Expected CompetitionLevelLocal 'local', got '%s'", modelmongo.CompetitionLevelLocal)
	}
}

func TestPublicationTypeConstants(t *testing.T) {
	if modelmongo.PublicationTypeJournal != "journal" {
		t.Errorf("Expected PublicationTypeJournal 'journal', got '%s'", modelmongo.PublicationTypeJournal)
	}
	if modelmongo.PublicationTypeConference != "conference" {
		t.Errorf("Expected PublicationTypeConference 'conference', got '%s'", modelmongo.PublicationTypeConference)
	}
	if modelmongo.PublicationTypeBook != "book" {
		t.Errorf("Expected PublicationTypeBook 'book', got '%s'", modelmongo.PublicationTypeBook)
	}
}

func TestPeriod_StructCreation(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour)
	period := modelmongo.Period{
		Start: start,
		End:   end,
	}

	if period.Start != start {
		t.Errorf("Expected Start to match, got %v", period.Start)
	}
	if period.End != end {
		t.Errorf("Expected End to match, got %v", period.End)
	}
}

func TestPeriod_JSONMarshalling(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour)
	period := modelmongo.Period{
		Start: start,
		End:   end,
	}

	jsonData, err := json.Marshal(period)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["start"]; !exists {
		t.Error("Expected start field in JSON")
	}
	if _, exists := result["end"]; !exists {
		t.Error("Expected end field in JSON")
	}
}

func TestPeriod_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"start": "2024-01-01T00:00:00Z",
		"end": "2024-01-02T00:00:00Z"
	}`

	var period modelmongo.Period
	if err := json.Unmarshal([]byte(jsonStr), &period); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if period.Start.IsZero() {
		t.Error("Expected Start to be non-zero")
	}
	if period.End.IsZero() {
		t.Error("Expected End to be non-zero")
	}
}

func TestPeriod_ZeroValues(t *testing.T) {
	var period modelmongo.Period

	if !period.Start.IsZero() {
		t.Errorf("Expected zero Start, got %v", period.Start)
	}
	if !period.End.IsZero() {
		t.Errorf("Expected zero End, got %v", period.End)
	}
}

func TestAttachment_StructCreation(t *testing.T) {
	now := time.Now()
	attachment := modelmongo.Attachment{
		FileName:   "document.pdf",
		FileURL:    "https://example.com/files/document.pdf",
		FileType:   "application/pdf",
		UploadedAt: now,
	}

	if attachment.FileName != "document.pdf" {
		t.Errorf("Expected FileName 'document.pdf', got '%s'", attachment.FileName)
	}
	if attachment.FileURL != "https://example.com/files/document.pdf" {
		t.Errorf("Expected FileURL 'https://example.com/files/document.pdf', got '%s'", attachment.FileURL)
	}
	if attachment.FileType != "application/pdf" {
		t.Errorf("Expected FileType 'application/pdf', got '%s'", attachment.FileType)
	}
}

func TestAttachment_JSONMarshalling(t *testing.T) {
	now := time.Now()
	attachment := modelmongo.Attachment{
		FileName:   "document.pdf",
		FileURL:    "https://example.com/files/document.pdf",
		FileType:   "application/pdf",
		UploadedAt: now,
	}

	jsonData, err := json.Marshal(attachment)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["fileName"] != "document.pdf" {
		t.Errorf("Expected fileName 'document.pdf', got '%v'", result["fileName"])
	}
	if result["fileUrl"] != "https://example.com/files/document.pdf" {
		t.Errorf("Expected fileUrl 'https://example.com/files/document.pdf', got '%v'", result["fileUrl"])
	}
	if result["fileType"] != "application/pdf" {
		t.Errorf("Expected fileType 'application/pdf', got '%v'", result["fileType"])
	}
}

func TestAttachment_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"fileName": "document.pdf",
		"fileUrl": "https://example.com/files/document.pdf",
		"fileType": "application/pdf",
		"uploadedAt": "2024-01-01T00:00:00Z"
	}`

	var attachment modelmongo.Attachment
	if err := json.Unmarshal([]byte(jsonStr), &attachment); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if attachment.FileName != "document.pdf" {
		t.Errorf("Expected FileName 'document.pdf', got '%s'", attachment.FileName)
	}
	if attachment.FileURL != "https://example.com/files/document.pdf" {
		t.Errorf("Expected FileURL 'https://example.com/files/document.pdf', got '%s'", attachment.FileURL)
	}
}

func TestAttachment_ZeroValues(t *testing.T) {
	var attachment modelmongo.Attachment

	if attachment.FileName != "" {
		t.Errorf("Expected empty FileName, got '%s'", attachment.FileName)
	}
	if attachment.FileURL != "" {
		t.Errorf("Expected empty FileURL, got '%s'", attachment.FileURL)
	}
	if attachment.FileType != "" {
		t.Errorf("Expected empty FileType, got '%s'", attachment.FileType)
	}
	if !attachment.UploadedAt.IsZero() {
		t.Errorf("Expected zero UploadedAt, got %v", attachment.UploadedAt)
	}
}

func TestAchievementDetails_StructCreation(t *testing.T) {
	competitionName := "Lomba Programming"
	competitionLevel := modelmongo.CompetitionLevelNational
	rank := 1
	medalType := "Gold"
	publicationType := modelmongo.PublicationTypeJournal
	publicationTitle := "Research Paper"
	authors := []string{"Author 1", "Author 2"}
	publisher := "Publisher Name"
	issn := "1234-5678"
	organizationName := "Organization Name"
	position := "President"
	start := time.Now()
	end := start.AddDate(0, 6, 0)
	period := modelmongo.Period{Start: start, End: end}
	certificationName := "Certification Name"
	issuedBy := "Issuer Name"
	certificationNumber := "CERT-123"
	validUntil := start.AddDate(1, 0, 0)
	eventDate := start
	location := "Jakarta"
	organizer := "Organizer Name"
	score := 95.5
	customFields := map[string]interface{}{
		"custom1": "value1",
		"custom2": 123,
	}

	details := modelmongo.AchievementDetails{
		CompetitionName:     &competitionName,
		CompetitionLevel:    &competitionLevel,
		Rank:                &rank,
		MedalType:           &medalType,
		PublicationType:     &publicationType,
		PublicationTitle:    &publicationTitle,
		Authors:             authors,
		Publisher:           &publisher,
		ISSN:                &issn,
		OrganizationName:    &organizationName,
		Position:            &position,
		Period:              &period,
		CertificationName:   &certificationName,
		IssuedBy:            &issuedBy,
		CertificationNumber: &certificationNumber,
		ValidUntil:          &validUntil,
		EventDate:           &eventDate,
		Location:            &location,
		Organizer:           &organizer,
		Score:               &score,
		CustomFields:        customFields,
	}

	if details.CompetitionName == nil || *details.CompetitionName != "Lomba Programming" {
		t.Errorf("Expected CompetitionName 'Lomba Programming', got '%v'", details.CompetitionName)
	}
	if details.CompetitionLevel == nil || *details.CompetitionLevel != modelmongo.CompetitionLevelNational {
		t.Errorf("Expected CompetitionLevel '%s', got '%v'", modelmongo.CompetitionLevelNational, details.CompetitionLevel)
	}
	if details.Rank == nil || *details.Rank != 1 {
		t.Errorf("Expected Rank 1, got '%v'", details.Rank)
	}
	if len(details.Authors) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(details.Authors))
	}
	if details.Score == nil || *details.Score != 95.5 {
		t.Errorf("Expected Score 95.5, got '%v'", details.Score)
	}
	if len(details.CustomFields) != 2 {
		t.Errorf("Expected 2 custom fields, got %d", len(details.CustomFields))
	}
}

func TestAchievementDetails_JSONMarshalling(t *testing.T) {
	competitionName := "Lomba Programming"
	competitionLevel := modelmongo.CompetitionLevelNational
	rank := 1
	details := modelmongo.AchievementDetails{
		CompetitionName:  &competitionName,
		CompetitionLevel: &competitionLevel,
		Rank:             &rank,
		Authors:          []string{"Author 1"},
	}

	jsonData, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["competitionName"] != "Lomba Programming" {
		t.Errorf("Expected competitionName 'Lomba Programming', got '%v'", result["competitionName"])
	}
	if result["competitionLevel"] != "national" {
		t.Errorf("Expected competitionLevel 'national', got '%v'", result["competitionLevel"])
	}
	if result["rank"] != float64(1) {
		t.Errorf("Expected rank 1, got '%v'", result["rank"])
	}
}

func TestAchievementDetails_JSONMarshallingOmitempty(t *testing.T) {
	details := modelmongo.AchievementDetails{}

	jsonData, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["competitionName"]; exists {
		t.Error("Expected competitionName to be omitted when nil")
	}
	if _, exists := result["rank"]; exists {
		t.Error("Expected rank to be omitted when nil")
	}
	if _, exists := result["authors"]; exists {
		t.Error("Expected authors to be omitted when empty")
	}
}

func TestAchievementDetails_ZeroValues(t *testing.T) {
	var details modelmongo.AchievementDetails

	if details.CompetitionName != nil {
		t.Error("Expected CompetitionName to be nil")
	}
	if details.CompetitionLevel != nil {
		t.Error("Expected CompetitionLevel to be nil")
	}
	if details.Rank != nil {
		t.Error("Expected Rank to be nil")
	}
	if details.Authors != nil {
		t.Error("Expected Authors to be nil")
	}
	if details.CustomFields != nil {
		t.Error("Expected CustomFields to be nil")
	}
}

func TestAchievement_StructCreation(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	competitionName := "Lomba Programming"
	competitionLevel := modelmongo.CompetitionLevelNational
	rank := 1
	details := modelmongo.AchievementDetails{
		CompetitionName:  &competitionName,
		CompetitionLevel: &competitionLevel,
		Rank:             &rank,
	}
	attachment := modelmongo.Attachment{
		FileName:   "document.pdf",
		FileURL:    "https://example.com/files/document.pdf",
		FileType:   "application/pdf",
		UploadedAt: now,
	}
	tags := []string{"programming", "competition"}
	deletedAt := now.Add(1 * time.Hour)

	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeCompetition,
		Title:           "Juara 1 Lomba Programming",
		Description:     "Menjadi juara 1 dalam lomba programming nasional",
		Details:         details,
		Attachments:     []modelmongo.Attachment{attachment},
		Tags:            tags,
		Points:          100,
		DeletedAt:       &deletedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if achievement.ID != id {
		t.Errorf("Expected ID to match, got %v", achievement.ID)
	}
	if achievement.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", achievement.StudentID)
	}
	if achievement.AchievementType != modelmongo.AchievementTypeCompetition {
		t.Errorf("Expected AchievementType '%s', got '%s'", modelmongo.AchievementTypeCompetition, achievement.AchievementType)
	}
	if achievement.Title != "Juara 1 Lomba Programming" {
		t.Errorf("Expected Title 'Juara 1 Lomba Programming', got '%s'", achievement.Title)
	}
	if len(achievement.Attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(achievement.Attachments))
	}
	if len(achievement.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(achievement.Tags))
	}
	if achievement.Points != 100 {
		t.Errorf("Expected Points 100, got %d", achievement.Points)
	}
}

func TestAchievement_JSONMarshalling(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeAcademic,
		Title:           "Juara 1 Lomba Programming",
		Description:     "Menjadi juara 1 dalam lomba programming nasional",
		Details:         modelmongo.AchievementDetails{},
		Attachments:     []modelmongo.Attachment{},
		Tags:            []string{"programming"},
		Points:          100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	jsonData, err := json.Marshal(achievement)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != id.Hex() {
		t.Errorf("Expected id '%s', got '%v'", id.Hex(), result["id"])
	}
	if result["studentId"] != "student-id-1" {
		t.Errorf("Expected studentId 'student-id-1', got '%v'", result["studentId"])
	}
	if result["achievementType"] != "academic" {
		t.Errorf("Expected achievementType 'academic', got '%v'", result["achievementType"])
	}
	if result["title"] != "Juara 1 Lomba Programming" {
		t.Errorf("Expected title 'Juara 1 Lomba Programming', got '%v'", result["title"])
	}
	if result["points"] != float64(100) {
		t.Errorf("Expected points 100, got '%v'", result["points"])
	}
}

func TestAchievement_JSONUnmarshalling(t *testing.T) {
	id := primitive.NewObjectID()
	jsonStr := `{
		"id": "` + id.Hex() + `",
		"studentId": "student-id-1",
		"achievementType": "academic",
		"title": "Juara 1 Lomba Programming",
		"description": "Menjadi juara 1 dalam lomba programming nasional",
		"details": {},
		"attachments": [],
		"tags": ["programming"],
		"points": 100,
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-01T00:00:00Z"
	}`

	var achievement modelmongo.Achievement
	if err := json.Unmarshal([]byte(jsonStr), &achievement); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if achievement.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", achievement.StudentID)
	}
	if achievement.AchievementType != "academic" {
		t.Errorf("Expected AchievementType 'academic', got '%s'", achievement.AchievementType)
	}
	if achievement.Title != "Juara 1 Lomba Programming" {
		t.Errorf("Expected Title 'Juara 1 Lomba Programming', got '%s'", achievement.Title)
	}
	if achievement.Points != 100 {
		t.Errorf("Expected Points 100, got %d", achievement.Points)
	}
}

func TestAchievement_ZeroValues(t *testing.T) {
	var achievement modelmongo.Achievement

	if achievement.StudentID != "" {
		t.Errorf("Expected empty StudentID, got '%s'", achievement.StudentID)
	}
	if achievement.AchievementType != "" {
		t.Errorf("Expected empty AchievementType, got '%s'", achievement.AchievementType)
	}
	if achievement.Title != "" {
		t.Errorf("Expected empty Title, got '%s'", achievement.Title)
	}
	if achievement.Points != 0 {
		t.Errorf("Expected Points 0, got %d", achievement.Points)
	}
	if achievement.Attachments != nil {
		t.Error("Expected Attachments to be nil")
	}
	if achievement.Tags != nil {
		t.Error("Expected Tags to be nil")
	}
	if achievement.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil")
	}
	if !achievement.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", achievement.CreatedAt)
	}
}

func TestAchievement_JSONMarshallingOmitempty(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeAcademic,
		Title:           "Juara 1 Lomba Programming",
		Description:     "Menjadi juara 1 dalam lomba programming nasional",
		Details:         modelmongo.AchievementDetails{},
		Points:          100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	jsonData, err := json.Marshal(achievement)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["attachments"]; exists {
		t.Error("Expected attachments to be omitted when empty")
	}
	if _, exists := result["tags"]; exists {
		t.Error("Expected tags to be omitted when empty")
	}
	if _, exists := result["deletedAt"]; exists {
		t.Error("Expected deletedAt to be omitted when nil")
	}
}

func TestCreateAchievementRequest_StructCreation(t *testing.T) {
	competitionName := "Lomba Programming"
	competitionLevel := modelmongo.CompetitionLevelNational
	rank := 1
	details := modelmongo.AchievementDetails{
		CompetitionName:  &competitionName,
		CompetitionLevel: &competitionLevel,
		Rank:             &rank,
	}
	attachment := modelmongo.Attachment{
		FileName:   "document.pdf",
		FileURL:    "https://example.com/files/document.pdf",
		FileType:   "application/pdf",
		UploadedAt: time.Now(),
	}
	tags := []string{"programming", "competition"}

	req := modelmongo.CreateAchievementRequest{
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeCompetition,
		Title:           "Juara 1 Lomba Programming",
		Description:     "Menjadi juara 1 dalam lomba programming nasional",
		Details:         details,
		Attachments:     []modelmongo.Attachment{attachment},
		Tags:            tags,
		Points:          100,
	}

	if req.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", req.StudentID)
	}
	if req.AchievementType != modelmongo.AchievementTypeCompetition {
		t.Errorf("Expected AchievementType '%s', got '%s'", modelmongo.AchievementTypeCompetition, req.AchievementType)
	}
	if req.Title != "Juara 1 Lomba Programming" {
		t.Errorf("Expected Title 'Juara 1 Lomba Programming', got '%s'", req.Title)
	}
	if req.Points != 100 {
		t.Errorf("Expected Points 100, got %d", req.Points)
	}
	if len(req.Attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(req.Attachments))
	}
}

func TestCreateAchievementRequest_JSONMarshalling(t *testing.T) {
	req := modelmongo.CreateAchievementRequest{
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeAcademic,
		Title:           "Juara 1 Lomba Programming",
		Description:     "Menjadi juara 1 dalam lomba programming nasional",
		Details:         modelmongo.AchievementDetails{},
		Attachments:     []modelmongo.Attachment{},
		Tags:            []string{"programming"},
		Points:          100,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["studentId"] != "student-id-1" {
		t.Errorf("Expected studentId 'student-id-1', got '%v'", result["studentId"])
	}
	if result["achievementType"] != "academic" {
		t.Errorf("Expected achievementType 'academic', got '%v'", result["achievementType"])
	}
	if result["title"] != "Juara 1 Lomba Programming" {
		t.Errorf("Expected title 'Juara 1 Lomba Programming', got '%v'", result["title"])
	}
	if result["points"] != float64(100) {
		t.Errorf("Expected points 100, got '%v'", result["points"])
	}
}

func TestCreateAchievementRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"studentId": "student-id-1",
		"achievementType": "academic",
		"title": "Juara 1 Lomba Programming",
		"description": "Menjadi juara 1 dalam lomba programming nasional",
		"details": {},
		"attachments": [],
		"tags": ["programming"],
		"points": 100
	}`

	var req modelmongo.CreateAchievementRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.StudentID != "student-id-1" {
		t.Errorf("Expected StudentID 'student-id-1', got '%s'", req.StudentID)
	}
	if req.AchievementType != "academic" {
		t.Errorf("Expected AchievementType 'academic', got '%s'", req.AchievementType)
	}
	if req.Points != 100 {
		t.Errorf("Expected Points 100, got %d", req.Points)
	}
}

func TestUpdateAchievementRequest_StructCreation(t *testing.T) {
	points := 150
	competitionName := "Lomba Programming"
	details := &modelmongo.AchievementDetails{
		CompetitionName: &competitionName,
	}
	attachment := modelmongo.Attachment{
		FileName:   "new-document.pdf",
		FileURL:    "https://example.com/files/new-document.pdf",
		FileType:   "application/pdf",
		UploadedAt: time.Now(),
	}
	tags := []string{"programming", "updated"}

	req := modelmongo.UpdateAchievementRequest{
		AchievementType: modelmongo.AchievementTypeCompetition,
		Title:           "Updated Title",
		Description:     "Updated Description",
		Details:         details,
		Attachments:     []modelmongo.Attachment{attachment},
		Tags:            tags,
		Points:          &points,
	}

	if req.AchievementType != modelmongo.AchievementTypeCompetition {
		t.Errorf("Expected AchievementType '%s', got '%s'", modelmongo.AchievementTypeCompetition, req.AchievementType)
	}
	if req.Title != "Updated Title" {
		t.Errorf("Expected Title 'Updated Title', got '%s'", req.Title)
	}
	if req.Points == nil || *req.Points != 150 {
		t.Errorf("Expected Points 150, got '%v'", req.Points)
	}
	if req.Details == nil {
		t.Error("Expected Details to be non-nil")
	}
}

func TestUpdateAchievementRequest_JSONMarshalling(t *testing.T) {
	points := 150
	req := modelmongo.UpdateAchievementRequest{
		AchievementType: modelmongo.AchievementTypeCompetition,
		Title:           "Updated Title",
		Description:     "Updated Description",
		Points:          &points,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["achievementType"] != "competition" {
		t.Errorf("Expected achievementType 'competition', got '%v'", result["achievementType"])
	}
	if result["title"] != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%v'", result["title"])
	}
	if result["points"] != float64(150) {
		t.Errorf("Expected points 150, got '%v'", result["points"])
	}
}

func TestUpdateAchievementRequest_JSONMarshallingOmitempty(t *testing.T) {
	req := modelmongo.UpdateAchievementRequest{
		Title: "Updated Title",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["achievementType"]; exists {
		t.Error("Expected achievementType to be omitted when empty")
	}
	if _, exists := result["points"]; exists {
		t.Error("Expected points to be omitted when nil")
	}
	if _, exists := result["details"]; exists {
		t.Error("Expected details to be omitted when nil")
	}
}

func TestGetAllAchievementsResponse_StructCreation(t *testing.T) {
	now := time.Now()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	achievements := []modelmongo.Achievement{
		{
			ID:              id1,
			StudentID:       "student-id-1",
			AchievementType: modelmongo.AchievementTypeAcademic,
			Title:           "Achievement 1",
			Description:     "Description 1",
			Details:         modelmongo.AchievementDetails{},
			Points:          100,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			ID:              id2,
			StudentID:       "student-id-2",
			AchievementType: modelmongo.AchievementTypeCompetition,
			Title:           "Achievement 2",
			Description:     "Description 2",
			Details:         modelmongo.AchievementDetails{},
			Points:          150,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	resp := modelmongo.GetAllAchievementsResponse{
		Status: "success",
		Data:   achievements,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 achievements, got %d", len(resp.Data))
	}
	if resp.Data[0].Title != "Achievement 1" {
		t.Errorf("Expected first achievement Title 'Achievement 1', got '%s'", resp.Data[0].Title)
	}
}

func TestGetAllAchievementsResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievements := []modelmongo.Achievement{
		{
			ID:              id,
			StudentID:       "student-id-1",
			AchievementType: modelmongo.AchievementTypeAcademic,
			Title:           "Achievement 1",
			Description:     "Description 1",
			Details:         modelmongo.AchievementDetails{},
			Points:          100,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	resp := modelmongo.GetAllAchievementsResponse{
		Status: "success",
		Data:   achievements,
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
			t.Errorf("Expected 1 achievement in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetAchievementByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeAcademic,
		Title:           "Achievement 1",
		Description:     "Description 1",
		Details:         modelmongo.AchievementDetails{},
		Points:          100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	resp := modelmongo.GetAchievementByIDResponse{
		Status: "success",
		Data:   achievement,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != id {
		t.Errorf("Expected Data.ID to match, got %v", resp.Data.ID)
	}
}

func TestCreateAchievementResponse_StructCreation(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeAcademic,
		Title:           "Achievement 1",
		Description:     "Description 1",
		Details:         modelmongo.AchievementDetails{},
		Points:          100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	resp := modelmongo.CreateAchievementResponse{
		Status: "success",
		Data:   achievement,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Title != "Achievement 1" {
		t.Errorf("Expected Data.Title 'Achievement 1', got '%s'", resp.Data.Title)
	}
}

func TestUpdateAchievementResponse_StructCreation(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	achievement := modelmongo.Achievement{
		ID:              id,
		StudentID:       "student-id-1",
		AchievementType: modelmongo.AchievementTypeCompetition,
		Title:           "Updated Achievement",
		Description:     "Updated Description",
		Details:         modelmongo.AchievementDetails{},
		Points:          150,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	resp := modelmongo.UpdateAchievementResponse{
		Status: "success",
		Data:   achievement,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Points != 150 {
		t.Errorf("Expected Data.Points 150, got %d", resp.Data.Points)
	}
}

func TestDeleteAchievementResponse_StructCreation(t *testing.T) {
	resp := modelmongo.DeleteAchievementResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteAchievementResponse_JSONMarshalling(t *testing.T) {
	resp := modelmongo.DeleteAchievementResponse{
		Status: "success",
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
}
