-- Sample Data untuk PostgreSQL
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
WHERE u.username IN ('mahasiswa1', 'mahasiswa2', 'mahasiswa3');

