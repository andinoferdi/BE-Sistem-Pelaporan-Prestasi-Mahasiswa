# Sistem Pelaporan Prestasi Mahasiswa

**Nama:** Andino Ferdiansah  
**NIM:** 434231065  
**Kelas:** C4

## API Endpoints

### 5.1 Authentication

#### POST /api/v1/auth/login

```json
{
  "username": "admin",
  "password": "123123123"
}
```

#### POST /api/v1/auth/refresh

```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /api/v1/auth/logout

#### GET /api/v1/auth/profile

### 5.2 Users (Admin)

#### GET /api/v1/users

#### GET /api/v1/users/:id

#### POST /api/v1/users

**Contoh untuk Mahasiswa:**
```json
{
  "username": "mahasiswa010",
  "email": "mahasiswa010@university.ac.id",
  "password": "123123123",
  "full_name": "Jennie biebier",
  "role_id": "6ff51ff6-c212-4b2d-b2e3-2e8c06059f90",
  "student_id": "M010",
  "program_study": "Teknik Informatika",
  "academic_year": "2024",
  "advisor_id": "8b063da7-9b5b-43b5-8dc3-3e67019d6c81",
  "is_active": true
}
```

**Contoh untuk Dosen Wali:**
```json
{
  "username": "dosen012",
  "email": "dosen0120@university.ac.id",
  "password": "123123123",
  "full_name": "Aciel Willow",
  "role_id": "036285f8-c16b-4a10-9ab6-cab1498cd347",
  "is_active": true,
  "lecturer_id": "dosen012",
  "department": "Teknik Informatika"
}
```

**Contoh untuk Admin:**
```json
{
  "username": "admin2",
  "email": "admin2@university.ac.id",
  "password": "123123123",
  "full_name": "Admin User",
  "role_id": "admin-role-id",
  "is_active": true
}
```

#### PUT /api/v1/users/:id

```json
{
  "username": "mahasiswa010",
  "email": "mahasiswa010@university.ac.id",
  "full_name": "Jennie biebier Updated",
  "role_id": "6ff51ff6-c212-4b2d-b2e3-2e8c06059f90",
  "is_active": true
}
```

#### DELETE /api/v1/users/:id

#### PUT /api/v1/users/:id/role

```json
{
  "role_id": "036285f8-c16b-4a10-9ab6-cab1498cd347"
}
```

### 5.4 Achievements

#### GET /api/v1/achievements

Query params: `page`, `limit`, `status`, `achievementType`, `sortBy`, `sortOrder`

#### GET /api/v1/achievements/:id

#### POST /api/v1/achievements

**Contoh untuk Competition:**
```json
{
  "studentId": "student-uuid",
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

**Contoh untuk Publication:**
```json
{
  "studentId": "student-uuid",
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

**Contoh untuk Organization:**
```json
{
  "studentId": "student-uuid",
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

**Contoh untuk Certification:**
```json
{
  "studentId": "student-uuid",
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

**Contoh untuk Academic:**
```json
{
  "studentId": "student-uuid",
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

**Contoh untuk Other:**
```json
{
  "studentId": "student-uuid",
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

#### PUT /api/v1/achievements/:id

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

#### DELETE /api/v1/achievements/:id

#### POST /api/v1/achievements/:id/submit

#### POST /api/v1/achievements/:id/verify

#### POST /api/v1/achievements/:id/reject

```json
{
  "rejection_note": "Data tidak lengkap"
}
```

#### GET /api/v1/achievements/:id/history

#### POST /api/v1/achievements/:id/attachments

Multipart form-data dengan key `file` (PDF, JPG, PNG, DOC, DOCX, max 10MB)

### 5.5 Students & Lecturers

#### GET /api/v1/students

#### GET /api/v1/students/:id

#### GET /api/v1/students/:id/achievements

#### PUT /api/v1/students/:id/advisor

```json
{
  "advisor_id": "8b063da7-9b5b-43b5-8dc3-3e67019d6c81"
}
```

#### GET /api/v1/lecturers

#### GET /api/v1/lecturers/:id/advisees

### 5.8 Reports & Analytics

#### GET /api/v1/reports/statistics

#### GET /api/v1/reports/student/:id
