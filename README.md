# Sistem Pelaporan Prestasi Mahasiswa

**Nama:** Andino Ferdiansah  
**NIM:** 434231065  
**Kelas:** C4

## Deskripsi

Sistem Pelaporan Prestasi Mahasiswa adalah aplikasi backend berbasis REST API yang memungkinkan mahasiswa melaporkan prestasi, dosen wali memverifikasi, dan admin mengelola sistem secara keseluruhan. Sistem ini menggunakan arsitektur dual database dengan PostgreSQL untuk data relasional (RBAC) dan MongoDB untuk data prestasi dinamis.

## Fitur Utama

- **Autentikasi & Otorisasi**
  - Login dengan JWT token
  - Role-Based Access Control (RBAC)
  - Permission-based authorization
  - Refresh token dan logout

- **Manajemen Prestasi**
  - Pelaporan prestasi dengan field dinamis
  - Berbagai tipe prestasi (akademik, kompetisi, organisasi, publikasi, sertifikasi)
  - Verifikasi prestasi oleh dosen wali
  - Status workflow (draft, submitted, verified, rejected)

- **Manajemen Pengguna**
  - Multi-role (Admin, Mahasiswa, Dosen Wali)
  - Manajemen permissions per role
  - Profile management

- **Database Migrations**
  - Manual migration via command
  - Schema creation untuk PostgreSQL
  - Data seeding untuk development
  - Support untuk PostgreSQL dan MongoDB

## Teknologi yang Digunakan

- **Framework:** Go Fiber v2
- **Database:**
  - PostgreSQL (data relasional, RBAC)
  - MongoDB (data prestasi dinamis)
- **Authentication:** JWT (JSON Web Token)
- **Password Hashing:** bcrypt
- **Language:** Go 1.21+

## Struktur Proyek

```
sistem-pelaporan-prestasi-mahasiswa/
├── app/
│   ├── model/
│   │   ├── mongo/          # Model untuk MongoDB
│   │   └── postgre/        # Model untuk PostgreSQL
│   ├── repository/
│   │   └── postgre/        # Data access layer
│   └── service/
│       └── postgre/        # Business logic layer
├── config/
│   ├── env.go              # Environment variables loader
│   ├── logger.go           # Logger configuration
│   └── mongo/
│       └── app.go          # Fiber app configuration
├── database/
│   ├── migration.go         # Database migrations
│   ├── mongo.go            # MongoDB connection
│   ├── postgre.go          # PostgreSQL connection
│   ├── mongo_schema.js     # MongoDB schema documentation
│   ├── postgre_schema.sql  # PostgreSQL schema
│   └── postgre_sample_data.sql  # PostgreSQL sample data
├── helper/
│   └── util.go             # Helper functions
├── middleware/
│   ├── logger.go           # Request logging middleware
│   └── postgre/
│       └── auth.go         # JWT & RBAC middleware
├── route/
│   └── postgre/
│       └── user_route.go   # Authentication routes
├── utils/
│   └── postgre/
│       ├── jwt.go          # JWT utilities
│       └── password.go     # Password hashing utilities
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
└── README.md               # Documentation
```

## Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL 12+
- MongoDB 4.4+
- Git

## Setup & Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/andinoferdi/Sistem-Pelaporan-Prestasi-Mahasiswa.git
cd Sistem-Pelaporan-Prestasi-Mahasiswa
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment Variables

Buat file `.env` di root project dengan konfigurasi berikut:

```env
# Application
APP_PORT=3001

# PostgreSQL
DB_DSN=postgres://username:password@localhost:5432/dbname?sslmode=disable

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=sppm_2025

# JWT
JWT_SECRET=your-secret-key-minimum-32-characters-long-for-production-security
```

### 4. Setup Database

**PostgreSQL:**
- Buat database baru
- Jalankan migration untuk membuat schema dan seed data

**MongoDB:**
- Pastikan MongoDB service berjalan
- Database dan collection akan dibuat saat migration dijalankan

### 5. Menjalankan Migration

Jalankan migration untuk membuat schema dan seed data:

```bash
go run cmd/migrate/main.go
```

Migration akan:
- Drop dan recreate semua tabel PostgreSQL
- Insert roles, permissions, users, lecturers, dan students
- Drop dan recreate collection achievements di MongoDB
- Membuat indexes untuk MongoDB

## Menjalankan Aplikasi

### 1. Jalankan Migration (Pertama Kali)

```bash
go run cmd/migrate/main.go
```

### 2. Development Mode

```bash
go run main.go
```

Server akan berjalan di `http://localhost:3001` (atau sesuai `APP_PORT` di `.env`, default: 3001)

### Build Binary

```bash
go build -o app main.go
./app
```

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required | Permission Required |
|--------|----------|-------------|---------------|---------------------|
| GET | `/api/v1/health` | Health check | No | - |
| POST | `/api/v1/auth/login` | Login user | No | - |
| POST | `/api/v1/auth/refresh` | Refresh JWT token | Yes | - |
| POST | `/api/v1/auth/logout` | Logout user | Yes | - |
| GET | `/api/v1/auth/profile` | Get user profile | Yes | - |

### Achievements

| Method | Endpoint | Description | Auth Required | Permission Required |
|--------|----------|-------------|---------------|---------------------|
| GET | `/api/v1/achievements` | Get all achievements | Yes | - |
| GET | `/api/v1/achievements/:id` | Get achievement by ID | Yes | `achievement:read` |
| POST | `/api/v1/achievements` | Create achievement | Yes | `achievement:create` |
| PUT | `/api/v1/achievements/:id` | Update achievement | Yes | `achievement:update` |
| DELETE | `/api/v1/achievements/:id` | Delete achievement | Yes | `achievement:delete` |
| POST | `/api/v1/achievements/upload` | Upload file | Yes | `achievement:create` |
| POST | `/api/v1/achievements/:id/submit` | Submit achievement | Yes | `achievement:update` |

## Tutorial API dengan Data Asli

### Sample Data yang Tersedia

Setelah migration, sistem akan memiliki data berikut:

**Users:**
- Admin: `admin` / `admin@gmail.com` (password: `12345678`)
- Dosen 1: `dosen1` / `dosen1@gmail.com` (password: `12345678`)
- Dosen 2: `dosen2` / `dosen2@gmail.com` (password: `12345678`)
- Dosen 3: `dosen3` / `dosen3@gmail.com` (password: `12345678`)
- Mahasiswa 1: `mahasiswa1` / `mahasiswa1@gmail.com` (password: `12345678`)
- Mahasiswa 2: `mahasiswa2` / `mahasiswa2@gmail.com` (password: `12345678`)
- Mahasiswa 3: `mahasiswa3` / `mahasiswa3@gmail.com` (password: `12345678`)

**Student IDs:**
- Mahasiswa 1: `202410001`
- Mahasiswa 2: `202410002`
- Mahasiswa 3: `202410003`

### 1. Health Check

**Request:**
```http
GET http://localhost:3001/api/v1/health
```

### 2. Login

**Request:**
```http
POST http://localhost:3001/api/v1/auth/login
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
  "username": "mahasiswa1",
  "password": "12345678"
}
```

Atau bisa menggunakan email:
```json
{
  "username": "mahasiswa1@gmail.com",
  "password": "12345678"
}
```

**Contoh untuk Admin:**
```json
{
  "username": "admin",
  "password": "12345678"
}
```

**Contoh untuk Dosen:**
```json
{
  "username": "dosen1",
  "password": "12345678"
}
```

### 3. Get Profile

**Request:**
```http
GET http://localhost:3001/api/v1/auth/profile
Authorization: Bearer <token-dari-login>
```

### 4. Refresh Token

**Request:**
```http
POST http://localhost:3001/api/v1/auth/refresh
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. Logout

**Request:**
```http
POST http://localhost:3001/api/v1/auth/logout
Authorization: Bearer <token>
```

### 6. Get All Achievements

**Request:**
```http
GET http://localhost:3001/api/v1/achievements
Authorization: Bearer <token-mahasiswa>
```

### 7. Get Achievement by ID

**Request:**
```http
GET http://localhost:3001/api/v1/achievements/<mongo-object-id>
Authorization: Bearer <token-mahasiswa>
```

### 8. Create Achievement

**Request:**
```http
POST http://localhost:3001/api/v1/achievements
Authorization: Bearer <token-mahasiswa>
Content-Type: application/json
```

**Body untuk Competition (raw JSON):**
```json
{
  "achievementType": "competition",
  "title": "Juara 1 Lomba Programming Nasional",
  "description": "Meraih juara 1 dalam National Programming Contest 2024",
  "details": {
    "competitionName": "National Programming Contest",
    "competitionLevel": "national",
    "rank": 1,
    "medalType": "gold",
    "eventDate": "2024-01-15T00:00:00Z",
    "location": "Jakarta",
    "organizer": "Kementerian Pendidikan"
  },
  "attachments": [],
  "tags": ["programming", "competition", "national"],
  "points": 100
}
```

**Body untuk Publication (raw JSON):**
```json
{
  "achievementType": "publication",
  "title": "Paper di Journal Internasional",
  "description": "Publikasi paper tentang Machine Learning di journal internasional",
  "details": {
    "publicationType": "journal",
    "publicationTitle": "Advanced Machine Learning Techniques",
    "authors": ["Andi Pratama", "Dr. Ahmad Wijaya"],
    "publisher": "IEEE",
    "issn": "1234-5678"
  },
  "attachments": [],
  "tags": ["publication", "journal", "machine-learning"],
  "points": 150
}
```

**Body untuk Organization (raw JSON):**
```json
{
  "achievementType": "organization",
  "title": "Ketua Himpunan Mahasiswa",
  "description": "Menjadi ketua himpunan mahasiswa teknik informatika",
  "details": {
    "organizationName": "Himpunan Mahasiswa Teknik Informatika",
    "position": "Ketua",
    "period": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-12-31T00:00:00Z"
    }
  },
  "attachments": [],
  "tags": ["organization", "leadership"],
  "points": 80
}
```

**Body untuk Certification (raw JSON):**
```json
{
  "achievementType": "certification",
  "title": "Sertifikasi AWS Solutions Architect",
  "description": "Mendapatkan sertifikasi AWS Solutions Architect Associate",
  "details": {
    "certificationName": "AWS Solutions Architect Associate",
    "issuedBy": "Amazon Web Services",
    "certificationNumber": "AWS-123456",
    "validUntil": "2026-01-15T00:00:00Z"
  },
  "attachments": [],
  "tags": ["certification", "aws", "cloud"],
  "points": 120
}
```

**Body untuk Academic (raw JSON):**
```json
{
  "achievementType": "academic",
  "title": "IPK 3.95 Semester 7",
  "description": "Mencapai IPK 3.95 pada semester 7",
  "details": {
    "score": 3.95,
    "eventDate": "2024-01-15T00:00:00Z"
  },
  "attachments": [],
  "tags": ["academic", "gpa"],
  "points": 50
}
```

**Body untuk Other (raw JSON):**
```json
{
  "achievementType": "other",
  "title": "Prestasi Lainnya",
  "description": "Deskripsi prestasi lainnya",
  "details": {
    "customFields": {
      "field1": "value1",
      "field2": "value2"
    }
  },
  "attachments": [],
  "tags": ["other"],
  "points": 30
}
```

### 9. Update Achievement

**Request:**
```http
PUT http://localhost:3001/api/v1/achievements/<mongo-object-id>
Authorization: Bearer <token-mahasiswa>
Content-Type: application/json
```

**Body (raw JSON) - hanya field yang ingin diupdate:**
```json
{
  "title": "Juara 1 Lomba Programming Internasional",
  "description": "Meraih juara 1 dalam International Programming Contest 2024",
  "details": {
    "competitionLevel": "international"
  },
  "points": 150
}
```

### 10. Upload File

**Request:**
```http
POST http://localhost:3001/api/v1/achievements/upload
Authorization: Bearer <token-mahasiswa>
Content-Type: multipart/form-data
```

**Body (form-data):**
- Key: `file`
- Type: File
- Value: Pilih file (PDF, JPG, PNG, DOC, DOCX, max 10MB)

**Cara menggunakan file yang diupload:**
Setelah upload, copy object attachment dari response dan masukkan ke array `attachments` saat create/update achievement:

```json
{
  "achievementType": "competition",
  "title": "Juara 1 Lomba Programming",
  "description": "Deskripsi prestasi",
  "details": {},
  "attachments": [
    {
      "fileName": "sertifikat.pdf",
      "fileUrl": "/uploads/1705312200-sertifikat.pdf",
      "fileType": "application/pdf",
      "uploadedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "tags": [],
  "points": 100
}
```

### 11. Submit Achievement

**Request:**
```http
POST http://localhost:3001/api/v1/achievements/<mongo-object-id>/submit
Authorization: Bearer <token-mahasiswa>
```

### 12. Delete Achievement

**Request:**
```http
DELETE http://localhost:3001/api/v1/achievements/<mongo-object-id>
Authorization: Bearer <token-mahasiswa>
```

**Catatan:** Achievement hanya bisa dihapus jika status masih `draft`.

## Catatan Penting

### Workflow Achievement

1. **Draft** - Prestasi baru dibuat, bisa di-edit dan dihapus
2. **Submitted** - Prestasi sudah di-submit, tidak bisa di-edit atau dihapus
3. **Verified** - Prestasi sudah diverifikasi dosen wali
4. **Rejected** - Prestasi ditolak oleh dosen wali

### Aturan Akses

- **Mahasiswa:**
  - Hanya bisa melihat prestasi miliknya sendiri
  - Hanya bisa create, update, delete prestasi miliknya
  - Hanya bisa submit prestasi miliknya
  - Update/delete hanya jika status `draft`

- **Dosen Wali:**
  - Bisa melihat prestasi mahasiswa bimbingannya
  - Bisa verify/reject prestasi (endpoint belum tersedia)

- **Admin:**
  - Akses penuh ke semua fitur

### Permission yang Tersedia

- `achievement:create` - Membuat prestasi baru
- `achievement:read` - Membaca data prestasi
- `achievement:update` - Mengupdate data prestasi
- `achievement:delete` - Menghapus data prestasi
- `achievement:verify` - Memverifikasi prestasi
- `user:manage` - Mengelola pengguna

### Tipe Achievement

- `academic` - Prestasi akademik (IPK, nilai, dll)
- `competition` - Prestasi kompetisi/lomba
- `organization` - Prestasi organisasi/kepanitiaan
- `publication` - Prestasi publikasi (jurnal, paper, dll)
- `certification` - Prestasi sertifikasi
- `other` - Prestasi lainnya


## Database Schema

### PostgreSQL Tables

- `roles` - Role definitions (Admin, Mahasiswa, Dosen Wali)
- `users` - User accounts
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mapping
- `lecturers` - Lecturer information
- `students` - Student information
- `achievement_references` - Achievement status tracking

### MongoDB Collections

- `achievements` - Dynamic achievement data dengan berbagai tipe:
  - Competition
  - Publication
  - Organization
  - Certification
  - Academic
  - Other

## Sample Data

Sistem melakukan seeding data saat migration dijalankan:

- **Roles:** Admin, Mahasiswa, Dosen Wali
- **Users:** 7 users (1 admin, 3 dosen, 3 mahasiswa)
- **Default Password:** `12345678` (untuk semua user)
- **Lecturers:** 3 dosen dengan ID DOS001, DOS002, DOS003
- **Students:** 3 mahasiswa dengan student ID 202410001, 202410002, 202410003

## Development

### Menjalankan Migrations

Migration tidak berjalan otomatis. Untuk menjalankan migration:

```bash
go run cmd/migrate/main.go
```

Migration akan:
- Drop semua tabel dan recreate schema PostgreSQL
- Seed data roles, permissions, users, lecturers, students
- Drop dan recreate collection achievements di MongoDB

**Peringatan:** Migration akan menghapus semua data yang ada. Gunakan dengan hati-hati di production.

### Logging

Logs ditulis ke console output dengan format:
```
[2024-01-15 10:30:00] GET /api/v1/achievements
[2024-01-15 10:30:00] GET /api/v1/achievements - 200 - 45ms
```

### File Upload

File yang diupload akan disimpan di folder `./uploads` dan dapat diakses melalui:
```
http://localhost:3001/uploads/<filename>
```

## License

Proyek ini dibuat untuk keperluan akademik.

## Author

**Andino Ferdiansah**  
NIM: 434231065  
Kelas: C4
