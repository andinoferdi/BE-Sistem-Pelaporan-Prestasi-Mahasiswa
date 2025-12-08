package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const postgresSchemaSQL = `DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;

DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS achievement_references CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS lecturers CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

DROP TYPE IF EXISTS achievement_status CASCADE;
DROP TYPE IF EXISTS notification_type CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE achievement_status AS ENUM ('draft', 'submitted', 'verified', 'rejected', 'deleted');

CREATE TYPE notification_type AS ENUM ('achievement_rejected', 'achievement_submitted');

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

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    achievement_id UUID REFERENCES achievement_references(id) ON DELETE CASCADE,
    mongo_achievement_id VARCHAR(24),
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_achievement_id ON notifications(achievement_id);
CREATE INDEX idx_notifications_mongo_achievement_id ON notifications(mongo_achievement_id);

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
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`

const postgresSampleDataSQL = `-- Sample Data untuk PostgreSQL
-- Jalankan file ini setelah menjalankan postgre_schema.sql

-- Hapus data yang sudah ada (jika ada)
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
-- Password untuk semua: 12345678

-- User Admin (1)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
SELECT 
    'admin',
    'admin@gmail.com',
    '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja',
    'Administrator',
    r.id,
    true
FROM roles r
WHERE r.name = 'Admin'
LIMIT 1;

-- Users Dosen Wali (3)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) VALUES
('dosen1', 'dosen1@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Prof. Dr. Ahmad Wijaya, S.T., M.T.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true),
('dosen2', 'dosen2@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Dr. Siti Nurhaliza, S.Kom., M.Kom.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true),
('dosen3', 'dosen3@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Dr. Budi Santoso, S.T., M.Sc.', (SELECT id FROM roles WHERE name = 'Dosen Wali'), true);

-- Users Mahasiswa (3)
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) VALUES
('mahasiswa1', 'mahasiswa1@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Andi Pratama', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true),
('mahasiswa2', 'mahasiswa2@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Budi Setiawan', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true),
('mahasiswa3', 'mahasiswa3@gmail.com', '$2a$12$iix7znEDxwTFySv47.9.2u6Uh3LYNBh/TcNRbBfqK0Sg24wWmdyja', 'Citra Dewi', (SELECT id FROM roles WHERE name = 'Mahasiswa'), true);

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
WHERE u.username IN ('mahasiswa1', 'mahasiswa2', 'mahasiswa3');`

// RunMigrations menjalankan migrasi PostgreSQL dan MongoDB secara berurutan.
func RunMigrations(postgresDB *sql.DB, mongoDB *mongo.Database) error {
	log.Println("Starting database migrations...")

	if err := runPostgresMigrations(postgresDB); err != nil {
		return fmt.Errorf("postgres migrations failed: %w", err)
	}

	if err := runMongoMigrations(mongoDB); err != nil {
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

func runMongoMigrations(db *mongo.Database) error {
	log.Println("Running MongoDB migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dropCollectionIfExists(ctx, db, "achievements"); err != nil {
		return err
	}

	if err := createAchievementIndexes(ctx, db); err != nil {
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
			Keys:    bson.D{{Key: "studentId", Value: 1}},
			Options: options.Index().SetName("idx_student_id"),
		},
		{
			Keys:    bson.D{{Key: "achievementType", Value: 1}},
			Options: options.Index().SetName("idx_achievement_type"),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
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
