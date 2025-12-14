package repository

// #1 proses: import library yang diperlukan untuk database, context, dan time
import (
	"context"
	"database/sql"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	"time"
)

// #2 proses: definisikan interface untuk operasi database achievement reference
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

// #3 proses: struct repository untuk operasi database achievement reference
type AchievementReferenceRepository struct {
	db *sql.DB
}

// #4 proses: constructor untuk membuat instance AchievementReferenceRepository baru
func NewAchievementReferenceRepository(db *sql.DB) IAchievementReferenceRepository {
	return &AchievementReferenceRepository{db: db}
}

// #5 proses: buat achievement reference baru di database
func (r *AchievementReferenceRepository) CreateAchievementReference(ctx context.Context, req model.CreateAchievementReferenceRequest) (*model.AchievementReference, error) {
	// #5a proses: query untuk insert achievement reference baru dengan RETURNING
	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, student_id, mongo_achievement_id, status, submitted_at, 
		          verified_at, verified_by, rejection_note, created_at, updated_at
	`

	// #5b proses: eksekusi query dan scan hasil ke struct ref
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

// #6 proses: ambil achievement reference berdasarkan mongo achievement ID
func (r *AchievementReferenceRepository) GetAchievementReferenceByMongoID(ctx context.Context, mongoID string) (*model.AchievementReference, error) {
	// #6a proses: query untuk ambil achievement reference, filter yang status bukan deleted
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1 AND status != 'deleted'
	`

	// #6b proses: eksekusi query dan scan hasil ke struct ref
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

// #7 proses: ambil achievement reference berdasarkan ID
func (r *AchievementReferenceRepository) GetAchievementReferenceByID(ctx context.Context, id string) (*model.AchievementReference, error) {
	// #7a proses: query untuk ambil achievement reference berdasarkan ID
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	// #7b proses: eksekusi query dan scan hasil ke struct ref
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

// #8 proses: update status achievement reference, bisa sekaligus set submitted_at
func (r *AchievementReferenceRepository) UpdateAchievementReferenceStatus(ctx context.Context, id string, status string, submittedAt *time.Time) error {
	// #8a proses: cek apakah submitted_at perlu diupdate juga
	var query string
	var err error

	if submittedAt != nil {
		// #8b proses: query untuk update status dan submitted_at
		query = `
			UPDATE achievement_references
			SET status = $1, submitted_at = $2, updated_at = NOW()
			WHERE id = $3
		`
		_, err = r.db.ExecContext(ctx, query, status, submittedAt, id)
	} else {
		// #8c proses: query untuk update status saja
		query = `
			UPDATE achievement_references
			SET status = $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = r.db.ExecContext(ctx, query, status, id)
	}

	return err
}

// #9 proses: hapus achievement reference dari database
func (r *AchievementReferenceRepository) DeleteAchievementReference(ctx context.Context, id string) error {
	// #9a proses: query untuk delete achievement reference berdasarkan ID
	query := `DELETE FROM achievement_references WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// #10 proses: ambil semua achievement reference milik student tertentu
func (r *AchievementReferenceRepository) GetAchievementReferenceByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	// #10a proses: query untuk ambil achievement reference berdasarkan student_id, filter yang status bukan deleted
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
	`

	// #10b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #10c proses: loop semua hasil dan masukkan ke slice references
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

// #11 proses: ambil semua achievement reference dari mahasiswa yang dibimbing dosen wali tertentu
func (r *AchievementReferenceRepository) GetAchievementReferencesByAdvisorID(ctx context.Context, advisorID string) ([]model.AchievementReference, error) {
	// #11a proses: query untuk ambil achievement reference dengan join ke tabel students
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at,
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`

	// #11b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #11c proses: loop semua hasil dan masukkan ke slice references
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

// #12 proses: ambil semua achievement reference yang ada di database
func (r *AchievementReferenceRepository) GetAllAchievementReferences(ctx context.Context) ([]model.AchievementReference, error) {
	// #12a proses: query untuk ambil semua reference, filter yang status bukan deleted
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
	`

	// #12b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #12c proses: loop semua hasil dan masukkan ke slice references
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

// #13 proses: ambil achievement reference student dengan pagination
func (r *AchievementReferenceRepository) GetAchievementReferenceByStudentIDPaginated(ctx context.Context, studentID string, page, limit int) ([]model.AchievementReference, int, error) {
	// #13a proses: hitung offset untuk pagination
	offset := (page - 1) * limit

	// #13b proses: query untuk hitung total reference student
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

	// #13c proses: query untuk ambil reference dengan limit dan offset
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	// #13d proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// #13e proses: loop semua hasil dan masukkan ke slice references
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

// #14 proses: ambil achievement reference dari mahasiswa bimbingan dosen dengan pagination
func (r *AchievementReferenceRepository) GetAchievementReferencesByAdvisorIDPaginated(ctx context.Context, advisorID string, page, limit int) ([]model.AchievementReference, int, error) {
	// #14a proses: hitung offset untuk pagination
	offset := (page - 1) * limit

	// #14b proses: query untuk hitung total reference dengan join ke students
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

	// #14c proses: query untuk ambil reference dengan limit dan offset
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at,
		       ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
		LIMIT $2 OFFSET $3
	`

	// #14d proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, advisorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// #14e proses: loop semua hasil dan masukkan ke slice references
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

// #15 proses: ambil semua achievement reference dengan pagination, filtering, dan sorting
func (r *AchievementReferenceRepository) GetAllAchievementReferencesPaginated(ctx context.Context, page, limit int, statusFilter string, sortBy string, sortOrder string) ([]model.AchievementReference, int, error) {
	// #15a proses: hitung offset untuk pagination
	offset := (page - 1) * limit

	// #15b proses: hitung total dengan atau tanpa status filter
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

	// #15c proses: validasi dan set default untuk sortBy dan sortOrder
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" || (sortOrder != "ASC" && sortOrder != "DESC") {
		sortOrder = "DESC"
	}

	// #15d proses: validasi sortBy hanya boleh field yang diizinkan
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

	// #15e proses: buat query dengan atau tanpa status filter
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

	// #15f proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// #15g proses: loop semua hasil dan masukkan ke slice references
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

// #16 proses: update status achievement reference jadi verified
func (r *AchievementReferenceRepository) UpdateAchievementReferenceVerify(ctx context.Context, id string, verifiedBy string) error {
	// #16a proses: query untuk update status jadi verified dan set verified_by serta verified_at
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, verified_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, model.AchievementStatusVerified, verifiedBy, id)
	return err
}

// #17 proses: update status achievement reference jadi rejected dengan catatan penolakan
func (r *AchievementReferenceRepository) UpdateAchievementReferenceReject(ctx context.Context, id string, verifiedBy string, rejectionNote string) error {
	// #17a proses: query untuk update status jadi rejected dan set verified_by serta rejection_note
	query := `
		UPDATE achievement_references
		SET status = $1, verified_by = $2, rejection_note = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, model.AchievementStatusRejected, verifiedBy, rejectionNote, id)
	return err
}

// #18 proses: ambil statistik achievement, hitung total dan yang sudah verified
func (r *AchievementReferenceRepository) GetAchievementStats(ctx context.Context) (int, int, error) {
	// #18a proses: query untuk hitung total dan verified sekaligus dengan FILTER
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

// #19 proses: ambil statistik achievement per periode bulan dalam rentang waktu tertentu
func (r *AchievementReferenceRepository) GetAchievementsByPeriod(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	// #19a proses: query untuk group by bulan dan hitung jumlah per periode
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

	// #19b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #19c proses: loop semua hasil dan masukkan ke map result
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

// #20 proses: ambil semua mongo achievement ID dari reference yang belum dihapus
func (r *AchievementReferenceRepository) GetAllAchievementMongoIDs(ctx context.Context) ([]string, error) {
	// #20a proses: query untuk ambil semua mongo_achievement_id
	query := `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE status != 'deleted'
	`

	// #20b proses: eksekusi query dan ambil semua baris hasil
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// #20c proses: loop semua hasil dan masukkan ke slice mongoIDs
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
