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
		WHERE mongo_achievement_id = $1 AND status != 'deleted'
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
		WHERE student_id = $1 AND status != 'deleted'
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

func GetAchievementReferencesByAdvisorID(db *sql.DB, advisorID string) ([]model.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at,
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`

	rows, err := db.Query(query, advisorID)
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

func GetAllAchievementReferences(db *sql.DB) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
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

func GetAchievementReferenceByStudentIDPaginated(db *sql.DB, studentID string, page, limit int) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
	`
	var total int
	err := db.QueryRow(countQuery, studentID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		references = append(references, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return references, total, nil
}

func GetAchievementReferencesByAdvisorIDPaginated(db *sql.DB, advisorID string, page, limit int) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
	`
	var total int
	err := db.QueryRow(countQuery, advisorID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at,
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, advisorID, limit, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		references = append(references, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return references, total, nil
}

func GetAllAchievementReferencesPaginated(db *sql.DB, page, limit int) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE status != 'deleted'
	`
	var total int
	err := db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		references = append(references, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return references, total, nil
}

func UpdateAchievementReferenceVerify(db *sql.DB, id string, verifiedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, verified_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`
	_, err := db.Exec(query, model.AchievementStatusVerified, verifiedBy, id)
	return err
}

func UpdateAchievementReferenceReject(db *sql.DB, id string, verifiedBy string, rejectionNote string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := db.Exec(query, model.AchievementStatusRejected, verifiedBy, rejectionNote, id)
	return err
}

func GetAchievementStats(db *sql.DB) (int, int, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'verified') as verified
		FROM achievement_references
		WHERE status != 'deleted'
	`
	var total, verified int
	err := db.QueryRow(query).Scan(&total, &verified)
	if err != nil {
		return 0, 0, err
	}
	return total, verified, nil
}

