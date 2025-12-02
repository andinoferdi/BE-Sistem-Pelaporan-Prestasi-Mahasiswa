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
  - Automatic schema creation
  - Automatic data seeding
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
- Database akan otomatis dibuat schema dan di-seed saat aplikasi pertama kali dijalankan

**MongoDB:**
- Pastikan MongoDB service berjalan
- Database dan collection akan otomatis dibuat saat aplikasi pertama kali dijalankan

## Menjalankan Aplikasi

### Development Mode

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

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/login` | Login user | No |
| POST | `/api/v1/auth/refresh` | Refresh JWT token | Yes |
| POST | `/api/v1/auth/logout` | Logout user | Yes |
| GET | `/api/v1/auth/profile` | Get user profile | Yes |

### Request/Response Examples

#### Login

**Request:**
```json
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@gmail.com",
  "password": "admin123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil.",
  "data": {
    "user": {
      "id": "uuid",
      "username": "admin",
      "email": "admin@gmail.com",
      "full_name": "Administrator",
      "role_id": "uuid",
      "is_active": true
    },
    "token": "jwt-token-here"
  }
}
```

#### Get Profile

**Request:**
```http
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Data profile berhasil diambil.",
  "data": {
    "user_id": "uuid",
    "username": "admin",
    "email": "admin@gmail.com",
    "full_name": "Administrator",
    "role_id": "uuid"
  }
}
```

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

Sistem secara otomatis melakukan seeding data saat pertama kali dijalankan:

- **Roles:** Admin, Mahasiswa, Dosen Wali
- **Users:** 7 users (1 admin, 3 dosen, 3 mahasiswa)
- **Default Password:** `admin123` (untuk semua user)
- **Achievements:** 5 sample achievements dengan berbagai tipe

## Development

### Menjalankan Migrations

Migrations berjalan otomatis saat aplikasi dijalankan. Untuk development, Anda bisa:

1. Hapus data di database
2. Restart aplikasi
3. Migrations akan otomatis membuat ulang schema dan seed data

### Logging

Logs ditulis ke `logs/app.log` (jika dikonfigurasi) dan console output.

## License

Proyek ini dibuat untuk keperluan akademik.

## Author

**Andino Ferdiansah**  
NIM: 434231065  
Kelas: C4
