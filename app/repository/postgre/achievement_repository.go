package repository

import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	"time"
)

type IAchievementReferenceRepository interface {
	CreateAchievementReference(ctx context.Context, req model.CreateAchievementReferenceRequest) (*model.AchievementReference, error)
	GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*model.AchievementReference, error)
	GetAchievementReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error)
	UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, submittedAt *time.Time) error
	DeleteAchievementReference(ctx context.Context, id string) error
	GetAchievementReferenceByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error)
	GetAchievementReferencesByAdvisorID(ctx context.Context, advisorID string) ([]model.AchievementReference, error)
	GetAllAchievementReferences(ctx context.Context) ([]model.AchievementReference, error)
	GetAchievementReferenceByStudentIDPaginated(ctx context.Context, studentID string, page, limit int) ([]model.AchievementReference, int, error)
	GetAchievementReferencesByAdvisorIDPaginated(ctx context.Context, advisorID string, page, limit int) ([]model.AchievementReference, int, error)
	GetAllAchievementReferencesPaginated(ctx context.Context, page, limit int, statusFilter string, sortBy string, sortOrder string) ([]model.AchievementReference, int, error)
	UpdateAchievementReferenceVerify(ctx context.Context, id string, verifiedBy string) error
	UpdateAchievementReferenceReject(ctx context.Context, id string, verifiedBy string, rejectionNote string) error
	GetAchievementStats(ctx context.Context) (int, int, error)
	GetAchievementsByPeriod(ctx context.Context, startDate, endDate time.Time) (map[string]int, error)
	GetAllAchievementMongoIDs(ctx context.Context) ([]string, error)
}

type AchievementReferenceRepository struct {
	db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) IAchievementReferenceRepository {
	return &AchievementReferenceRepository{db: db}
}

func (r *AchievementReferenceRepository) CreateAchievementReference(ctx context.Context, req model.CreateAchievementReferenceRequest) (*model.AchievementReference, error) {
	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, student_id, mongo_achievement_id, status, submitted_at, 
		          verified_at, verified_by, rejection_note, created_at, updated_at
	`

	ref := new(model.AchievementReference)
	err := r.db.QueryRowContext(ctx, query, req.StudentID, req.MongoAchievementID, req.Status).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (r *AchievementReferenceRepository) GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1 AND status != 'deleted'
	`

	ref := new(model.AchievementReference)
	err := r.db.QueryRowContext(ctx, query, mongoID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (r *AchievementReferenceRepository) GetAchievementReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	ref := new(model.AchievementReference)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (r *AchievementReferenceRepository) UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, submittedAt *time.Time) error {
	var query string
	var err error

	if submittedAt != nil {
		query = `
			UPDATE achievement_references
			SET status = $1, submitted_at = $2, updated_at = NOW()
			WHERE id = $3
		`
		_, err = r.db.ExecContext(ctx, query, status, submittedAt, id)
	} else {
		query = `
			UPDATE achievement_references
			SET status = $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = r.db.ExecContext(ctx, query, status, id)
	}

	return err
}

func (r *AchievementReferenceRepository) DeleteAchievementReference(ctx context.Context, id string) error {
	query := `DELETE FROM achievement_references WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *AchievementReferenceRepository) GetAchievementReferenceByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
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

func (r *AchievementReferenceRepository) GetAchievementReferencesByAdvisorID(ctx context.Context, advisorID string) ([]model.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at,
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, advisorID)
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

func (r *AchievementReferenceRepository) GetAllAchievementReferences(ctx context.Context) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
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

func (r *AchievementReferenceRepository) GetAchievementReferenceByStudentIDPaginated(ctx context.Context, studentID string, page, limit int) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, studentID).Scan(&total)
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

	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
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

func (r *AchievementReferenceRepository) GetAchievementReferencesByAdvisorIDPaginated(ctx context.Context, advisorID string, page, limit int) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, advisorID).Scan(&total)
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

	rows, err := r.db.QueryContext(ctx, query, advisorID, limit, offset)
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

func (r *AchievementReferenceRepository) GetAllAchievementReferencesPaginated(ctx context.Context, page, limit int, statusFilter string, sortBy string, sortOrder string) ([]model.AchievementReference, int, error) {
	offset := (page - 1) * limit

	var countQuery string
	var total int
	var err error

	if statusFilter != "" {
		countQuery = `
			SELECT COUNT(*)
			FROM achievement_references
			WHERE status = $1 AND status != 'deleted'
		`
		err = r.db.QueryRowContext(ctx, countQuery, statusFilter).Scan(&total)
	} else {
		countQuery = `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE status != 'deleted'
	`
		err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}

	if err != nil {
		return nil, 0, err
	}

	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" || (sortOrder != "ASC" && sortOrder != "DESC") {
		sortOrder = "DESC"
	}

	allowedSortBy := map[string]bool{
		"created_at":   true,
		"updated_at":   true,
		"submitted_at": true,
		"status":       true,
	}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}

	orderBy := sortBy + " " + sortOrder

	var query string
	var queryArgs []interface{}

	if statusFilter != "" {
		query = `
			SELECT id, student_id, mongo_achievement_id, status, submitted_at,
			       verified_at, verified_by, rejection_note, created_at, updated_at
			FROM achievement_references
			WHERE status = $1 AND status != 'deleted'
			ORDER BY ` + orderBy + `
			LIMIT $2 OFFSET $3
		`
		queryArgs = []interface{}{statusFilter, limit, offset}
	} else {
		query = `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
			ORDER BY ` + orderBy + `
		LIMIT $1 OFFSET $2
	`
		queryArgs = []interface{}{limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
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

func (r *AchievementReferenceRepository) UpdateAchievementReferenceVerify(ctx context.Context, id string, verifiedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, verified_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, model.AchievementStatusVerified, verifiedBy, id)
	return err
}

func (r *AchievementReferenceRepository) UpdateAchievementReferenceReject(ctx context.Context, id string, verifiedBy string, rejectionNote string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, model.AchievementStatusRejected, verifiedBy, rejectionNote, id)
	return err
}

func (r *AchievementReferenceRepository) GetAchievementStats(ctx context.Context) (int, int, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'verified') as verified
		FROM achievement_references
		WHERE status != 'deleted'
	`
	var total, verified int
	err := r.db.QueryRowContext(ctx, query).Scan(&total, &verified)
	if err != nil {
		return 0, 0, err
	}
	return total, verified, nil
}

func (r *AchievementReferenceRepository) GetAchievementsByPeriod(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	query := `
		SELECT 
			TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') as period,
			COUNT(*) as count
		FROM achievement_references
		WHERE status != 'deleted'
			AND created_at >= $1
			AND created_at <= $2
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY period
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var period string
		var count int
		if err := rows.Scan(&period, &count); err != nil {
			continue
		}
		result[period] = count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AchievementReferenceRepository) GetAllAchievementMongoIDs(ctx context.Context) ([]string, error) {
	query := `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE status != 'deleted'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mongoIDs []string
	for rows.Next() {
		var mongoID string
		if err := rows.Scan(&mongoID); err != nil {
			continue
		}
		mongoIDs = append(mongoIDs, mongoID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mongoIDs, nil
}
