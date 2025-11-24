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

INSERT INTO roles (name, description) VALUES
('Admin', 'Pengelola sistem dengan akses penuh'),
('Mahasiswa', 'Pelapor prestasi'),
('Dosen Wali', 'Verifikator prestasi mahasiswa bimbingannya');

INSERT INTO permissions (name, resource, action, description) VALUES
('achievement:create', 'achievement', 'create', 'Membuat prestasi baru'),
('achievement:read', 'achievement', 'read', 'Membaca data prestasi'),
('achievement:update', 'achievement', 'update', 'Mengupdate data prestasi'),
('achievement:delete', 'achievement', 'delete', 'Menghapus data prestasi'),
('achievement:verify', 'achievement', 'verify', 'Memverifikasi prestasi'),
('user:manage', 'user', 'manage', 'Mengelola pengguna');
