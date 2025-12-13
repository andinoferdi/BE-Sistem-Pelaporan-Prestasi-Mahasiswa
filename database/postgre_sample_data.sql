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

-- Insert Users (1 Admin)
-- Password: 123123123

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

