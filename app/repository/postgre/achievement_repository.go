package repository

import (
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	"time"
)

func CreateAchievementReference(db *sql.DB, req model.CreateAchievementReferenceRequest) (*model.AchievementReference, error) {
	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, student_id, mongo_achievement_id, status, submitted_at, 
		          verified_at, verified_by, rejection_note, created_at, updated_at
	`

	ref := new(model.AchievementReference)
	err := db.QueryRow(query, req.StudentID, req.MongoAchievementID, req.Status).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func GetAchievementReferenceByMongoID(db *sql.DB, mongoID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`

	ref := new(model.AchievementReference)
	err := db.QueryRow(query, mongoID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func GetAchievementReferenceByID(db *sql.DB, id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	ref := new(model.AchievementReference)
	err := db.QueryRow(query, id).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func UpdateAchievementReferenceStatus(db *sql.DB, id string, status string, submittedAt *time.Time) error {
	var query string
	var err error

	if submittedAt != nil {
		query = `
			UPDATE achievement_references
			SET status = $1, submitted_at = $2, updated_at = NOW()
			WHERE id = $3
		`
		_, err = db.Exec(query, status, submittedAt, id)
	} else {
		query = `
			UPDATE achievement_references
			SET status = $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = db.Exec(query, status, id)
	}

	return err
}

func DeleteAchievementReference(db *sql.DB, id string) error {
	query := `DELETE FROM achievement_references WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func GetAchievementReferenceByStudentID(db *sql.DB, studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		references = append(references, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return references, nil
}

