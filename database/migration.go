package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const postgresSchemaSQL = `DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;

DROP TABLE IF EXISTS achievement_references CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS lecturers CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

DROP TYPE IF EXISTS achievement_status CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE achievement_status AS ENUM ('draft', 'submitted', 'verified', 'rejected');

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE lecturers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID REFERENCES lecturers(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE achievement_references (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status achievement_status NOT NULL DEFAULT 'draft',
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_advisor_id ON students(advisor_id);
CREATE INDEX idx_lecturers_user_id ON lecturers(user_id);
CREATE INDEX idx_achievement_references_student_id ON achievement_references(student_id);
CREATE INDEX idx_achievement_references_status ON achievement_references(status);
CREATE INDEX idx_achievement_references_verified_by ON achievement_references(verified_by);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_achievement_references_updated_at BEFORE UPDATE ON achievement_references
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`

const postgresSampleDataSQL = `-- Sample Data untuk PostgreSQL
-- Jalankan file ini setelah menjalankan postgre_schema.sql

-- Hapus data yang sudah ada (jika ada)
DELETE FROM achievement_references;
DELETE FROM students;
DELETE FROM lecturers;
DELETE FROM role_permissions;
DELETE FROM users;
DELETE FROM permissions;
DELETE FROM roles;

-- Insert Roles
INSERT INTO roles (name, description) VALUES
('Admin', 'Pengelola sistem dengan akses penuh'),
('Mahasiswa', 'Pelapor prestasi'),
('Dosen Wali', 'Verifikator prestasi mahasiswa bimbingannya');

-- Insert Permissions
INSERT INTO permissions (name, resource, action, description) VALUES
('achievement:create', 'achievement', 'create', 'Membuat prestasi baru'),
('achievement:read', 'achievement', 'read', 'Membaca data prestasi'),
('achievement:update', 'achievement', 'update', 'Mengupdate data prestasi'),
('achievement:delete', 'achievement', 'delete', 'Menghapus data prestasi'),
('achievement:verify', 'achievement', 'verify', 'Memverifikasi prestasi'),
('user:manage', 'user', 'manage', 'Mengelola pengguna');

-- Insert Role Permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE (r.name = 'Admin' AND p.name IN (
    'achievement:create', 'achievement:read', 'achievement:update', 
    'achievement:delete', 'achievement:verify', 'user:manage'
))
OR (r.name = 'Mahasiswa' AND p.name IN (
    'achievement:create', 'achievement:read', 'achievement:update', 'achievement:delete'
))
OR (r.name = 'Dosen Wali' AND p.name IN (
    'achievement:read', 'achievement:verify'
));

-- Insert Users (Total 7: 1 Admin, 3 Dosen Wali, 3 Mahasiswa)
-- Password untuk semua: admin123

-- User Admin (1)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
SELECT 
    'admin',
    'admin@gmail.com',
    '$2y$10$yQxfaJ6HFY71GnpBAJ6f8Om3m8tTMj1jRXmTwz6HDFLDv4uWAPm7e',
    'Administrator',
    r.id,
    true
FROM roles r
WHERE r.name = 'Admin'
LIMIT 1;

-- Users Dosen Wali (3)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) VALUES
('dosen1', 'dosen1@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Prof. Dr. Ahmad Wijaya, S.T., M.T.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true),
('dosen2', 'dosen2@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Dr. Siti Nurhaliza, S.Kom., M.Kom.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true),
('dosen3', 'dosen3@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Dr. Budi Santoso, S.T., M.Sc.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true);

-- Users Mahasiswa (3)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) VALUES
('mahasiswa1', 'mahasiswa1@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Andi Pratama', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true),
('mahasiswa2', 'mahasiswa2@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Budi Setiawan', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true),
('mahasiswa3', 'mahasiswa3@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Citra Dewi', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true);

-- Insert Lecturers (3 data untuk 3 dosen wali)
INSERT INTO lecturers (user_id, lecturer_id, department)
SELECT 
    u.id,
    CASE u.username
        WHEN 'dosen1' THEN 'DOS001'
        WHEN 'dosen2' THEN 'DOS002'
        WHEN 'dosen3' THEN 'DOS003'
    END,
    'Teknik Informatika'
FROM users u
WHERE u.username IN ('dosen1', 'dosen2', 'dosen3');

-- Insert Students (3 data untuk 3 mahasiswa)
INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
SELECT 
    u.id,
    CASE u.username
        WHEN 'mahasiswa1' THEN '202410001'
        WHEN 'mahasiswa2' THEN '202410002'
        WHEN 'mahasiswa3' THEN '202410003'
    END,
    'Teknik Informatika',
    '2024',
    CASE u.username
        WHEN 'mahasiswa1' THEN (SELECT l.id FROM lecturers l JOIN users u2 ON l.user_id = u2.id WHERE u2.username = 'dosen1' LIMIT 1)
        WHEN 'mahasiswa2' THEN (SELECT l.id FROM lecturers l JOIN users u2 ON l.user_id = u2.id WHERE u2.username = 'dosen2' LIMIT 1)
        WHEN 'mahasiswa3' THEN (SELECT l.id FROM lecturers l JOIN users u2 ON l.user_id = u2.id WHERE u2.username = 'dosen3' LIMIT 1)
    END
FROM users u
WHERE u.username IN ('mahasiswa1', 'mahasiswa2', 'mahasiswa3');

-- Insert Achievement References
-- Sample dengan berbagai status: draft, submitted, verified, rejected
INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
SELECT 
    s.id,
    '507f1f77bcf86cd799439011',
    'draft',
    NULL,
    NULL,
    NULL,
    NULL,
    NOW(),
    NOW()
FROM students s
JOIN users u ON s.user_id = u.id
WHERE u.username = 'mahasiswa1'
LIMIT 1;

INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
SELECT 
    s.id,
    '507f1f77bcf86cd799439012',
    'submitted',
    NOW() - INTERVAL '2 days',
    NULL,
    NULL,
    NULL,
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '2 days'
FROM students s
JOIN users u ON s.user_id = u.id
WHERE u.username = 'mahasiswa1'
LIMIT 1;

INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
SELECT 
    s.id,
    '507f1f77bcf86cd799439013',
    'verified',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '3 days',
    (SELECT u.id FROM users u WHERE u.username = 'dosen1' LIMIT 1),
    NULL,
    NOW() - INTERVAL '6 days',
    NOW() - INTERVAL '3 days'
FROM students s
JOIN users u ON s.user_id = u.id
WHERE u.username = 'mahasiswa2'
LIMIT 1;

INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
SELECT 
    s.id,
    '507f1f77bcf86cd799439014',
    'rejected',
    NOW() - INTERVAL '4 days',
    NULL,
    (SELECT u.id FROM users u WHERE u.username = 'dosen2' LIMIT 1),
    'Dokumen tidak lengkap. Silakan lengkapi dokumen pendukung.',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '1 day'
FROM students s
JOIN users u ON s.user_id = u.id
WHERE u.username = 'mahasiswa2'
LIMIT 1;

INSERT INTO achievement_references (student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
SELECT 
    s.id,
    '507f1f77bcf86cd799439015',
    'submitted',
    NOW() - INTERVAL '1 day',
    NULL,
    NULL,
    NULL,
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '1 day'
FROM students s
JOIN users u ON s.user_id = u.id
WHERE u.username = 'mahasiswa3'
LIMIT 1;`

// RunMigrations menjalankan migrasi PostgreSQL dan MongoDB secara berurutan.
func RunMigrations(postgresDB *sql.DB, mongoDB *mongo.Database) error {
	log.Println("Starting database migrations...")

	if err := runPostgresMigrations(postgresDB); err != nil {
		return fmt.Errorf("postgres migrations failed: %w", err)
	}

	studentIDs, err := fetchStudentIDs(postgresDB)
	if err != nil {
		return fmt.Errorf("fetching student IDs: %w", err)
	}

	if err := runMongoMigrations(mongoDB, studentIDs); err != nil {
		return fmt.Errorf("mongo migrations failed: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func runPostgresMigrations(db *sql.DB) error {
	log.Println("Running PostgreSQL schema and seed migrations...")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if _, err := tx.Exec(postgresSchemaSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("executing schema SQL: %w", err)
	}

	if _, err := tx.Exec(postgresSampleDataSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("executing sample data SQL: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Println("PostgreSQL migrations completed")
	return nil
}

func fetchStudentIDs(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query(`
		SELECT s.id, u.username
		FROM students s
		INNER JOIN users u ON u.id = s.user_id
	`)
	if err != nil {
		return nil, fmt.Errorf("query students: %w", err)
	}
	defer rows.Close()

	studentIDs := make(map[string]string)
	for rows.Next() {
		var id string
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			return nil, fmt.Errorf("scan student row: %w", err)
		}
		studentIDs[username] = id
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student rows: %w", err)
	}

	return studentIDs, nil
}

func runMongoMigrations(db *mongo.Database, studentIDs map[string]string) error {
	log.Println("Running MongoDB migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dropCollectionIfExists(ctx, db, "achievements"); err != nil {
		return err
	}

	if err := createAchievementIndexes(ctx, db); err != nil {
		return err
	}

	if err := seedAchievementData(ctx, db, studentIDs); err != nil {
		return err
	}

	log.Println("MongoDB migrations completed")
	return nil
}

func dropCollectionIfExists(ctx context.Context, db *mongo.Database, collectionName string) error {
	names, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return fmt.Errorf("list collections for %s: %w", collectionName, err)
	}

	if len(names) == 0 {
		return nil
	}

	if err := db.Collection(collectionName).Drop(ctx); err != nil {
		return fmt.Errorf("drop collection %s: %w", collectionName, err)
	}

	log.Printf("Dropped collection: %s", collectionName)
	return nil
}

func createAchievementIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("achievements")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetName("idx_student_id"),
		},
		{
			Keys:    bson.D{{Key: "achievement_type", Value: 1}},
			Options: options.Index().SetName("idx_achievement_type"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "description", Value: "text"},
			},
			Options: options.Index().SetName("idx_text_search"),
		},
	}

	if _, err := collection.Indexes().CreateMany(ctx, indexModels); err != nil {
		return fmt.Errorf("create achievement indexes: %w", err)
	}

	log.Println("Created indexes for achievements collection")
	return nil
}

func seedAchievementData(ctx context.Context, db *mongo.Database, studentIDs map[string]string) error {
	log.Println("Seeding MongoDB achievements collection...")

	collection := db.Collection("achievements")

	studentMahasiswa1, err := studentIDFor(studentIDs, "mahasiswa1")
	if err != nil {
		return err
	}

	studentMahasiswa2, err := studentIDFor(studentIDs, "mahasiswa2")
	if err != nil {
		return err
	}

	studentMahasiswa3, err := studentIDFor(studentIDs, "mahasiswa3")
	if err != nil {
		return err
	}

	id1, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	if err != nil {
		return fmt.Errorf("invalid object id 507f1f77bcf86cd799439011: %w", err)
	}
	id2, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	if err != nil {
		return fmt.Errorf("invalid object id 507f1f77bcf86cd799439012: %w", err)
	}
	id3, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")
	if err != nil {
		return fmt.Errorf("invalid object id 507f1f77bcf86cd799439013: %w", err)
	}
	id4, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439014")
	if err != nil {
		return fmt.Errorf("invalid object id 507f1f77bcf86cd799439014: %w", err)
	}
	id5, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439015")
	if err != nil {
		return fmt.Errorf("invalid object id 507f1f77bcf86cd799439015: %w", err)
	}

	docs := []interface{}{
		bson.M{
			"_id":              id1,
			"student_id":       studentMahasiswa1,
			"achievement_type": "competition",
			"title":            "Juara 1 Lomba Programming Nasional",
			"description":      "Meraih juara 1 dalam kompetisi programming tingkat nasional",
			"details": bson.M{
				"competition_name":  "National Programming Contest 2025",
				"competition_level": "national",
				"rank":              1,
				"medal_type":        "Gold",
				"event_date":        time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
				"location":          "Jakarta",
				"organizer":         "Kementerian Pendidikan",
			},
			"attachments": []bson.M{
				{
					"file_name":   "sertifikat.pdf",
					"file_url":    "/uploads/sertifikat.pdf",
					"file_type":   "application/pdf",
					"uploaded_at": time.Date(2025, time.January, 20, 10, 0, 0, 0, time.UTC),
				},
			},
			"tags":       []string{"programming", "competition", "national"},
			"points":     100,
			"created_at": time.Date(2025, time.January, 20, 10, 0, 0, 0, time.UTC),
			"updated_at": time.Date(2025, time.January, 20, 10, 0, 0, 0, time.UTC),
		},
		bson.M{
			"_id":              id2,
			"student_id":       studentMahasiswa1,
			"achievement_type": "publication",
			"title":            "Paper di International Conference",
			"description":      "Publikasi paper di konferensi internasional",
			"details": bson.M{
				"publication_type":  "conference",
				"publication_title": "Machine Learning Applications in Education",
				"authors":           []string{"John Doe", "Jane Smith"},
				"publisher":         "IEEE",
				"issn":              "1234-5678",
				"event_date":        time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
			},
			"tags":       []string{"publication", "research", "conference"},
			"points":     150,
			"created_at": time.Date(2025, time.February, 15, 10, 0, 0, 0, time.UTC),
			"updated_at": time.Date(2025, time.February, 15, 10, 0, 0, 0, time.UTC),
		},
		bson.M{
			"_id":              id3,
			"student_id":       studentMahasiswa2,
			"achievement_type": "organization",
			"title":            "Ketua Himpunan Mahasiswa",
			"description":      "Menjadi ketua himpunan mahasiswa selama 1 tahun",
			"details": bson.M{
				"organization_name": "Himpunan Mahasiswa Teknik Informatika",
				"position":          "Ketua",
				"period": bson.M{
					"start": time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					"end":   time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
				},
			},
			"tags":       []string{"organization", "leadership"},
			"points":     80,
			"created_at": time.Date(2025, time.January, 5, 10, 0, 0, 0, time.UTC),
			"updated_at": time.Date(2025, time.January, 5, 10, 0, 0, 0, time.UTC),
		},
		bson.M{
			"_id":              id4,
			"student_id":       studentMahasiswa2,
			"achievement_type": "certification",
			"title":            "AWS Certified Solutions Architect",
			"description":      "Sertifikasi AWS Solutions Architect Associate",
			"details": bson.M{
				"certification_name":   "AWS Certified Solutions Architect - Associate",
				"issued_by":            "Amazon Web Services",
				"certification_number": "AWS-123456789",
				"valid_until":          time.Date(2027, time.December, 31, 0, 0, 0, 0, time.UTC),
				"event_date":           time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
			},
			"attachments": []bson.M{
				{
					"file_name":   "aws_certificate.pdf",
					"file_url":    "/uploads/aws_certificate.pdf",
					"file_type":   "application/pdf",
					"uploaded_at": time.Date(2025, time.March, 5, 10, 0, 0, 0, time.UTC),
				},
			},
			"tags":       []string{"certification", "aws", "cloud"},
			"points":     120,
			"created_at": time.Date(2025, time.March, 5, 10, 0, 0, 0, time.UTC),
			"updated_at": time.Date(2025, time.March, 5, 10, 0, 0, 0, time.UTC),
		},
		bson.M{
			"_id":              id5,
			"student_id":       studentMahasiswa3,
			"achievement_type": "academic",
			"title":            "IPK 3.95 Semester 7",
			"description":      "Mencapai IPK 3.95 pada semester 7",
			"details": bson.M{
				"score":      3.95,
				"event_date": time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			"tags":       []string{"academic", "gpa"},
			"points":     50,
			"created_at": time.Date(2025, time.February, 1, 10, 0, 0, 0, time.UTC),
			"updated_at": time.Date(2025, time.February, 1, 10, 0, 0, 0, time.UTC),
		},
	}

	if _, err := collection.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("insert achievements: %w", err)
	}

	log.Println("Seeded MongoDB achievements collection")
	return nil
}

func studentIDFor(studentIDs map[string]string, username string) (string, error) {
	id, ok := studentIDs[username]
	if !ok || id == "" {
		return "", fmt.Errorf("student id for %s not found", username)
	}
	return id, nil
}
