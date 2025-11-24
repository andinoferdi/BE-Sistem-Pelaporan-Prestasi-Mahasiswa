// MongoDB Schema dan Sample Data untuk Navicat
// Database: sppm_2025
// Collection: achievements

// Gunakan database
use('sppm_2025');

// Struktur Collection: achievements
// {
//   _id: ObjectId,
//   student_id: String (UUID dari PostgreSQL),
//   achievement_type: String, // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
//   title: String,
//   description: String,
//   details: Object (field dinamis berdasarkan achievement_type),
//   attachments: Array (optional),
//   tags: Array (optional),
//   points: Number,
//   created_at: Date,
//   updated_at: Date
// }

// Sample 1: Competition Achievement
db.achievements.insertOne({
  student_id: "550e8400-e29b-41d4-a716-446655440000",
  achievement_type: "competition",
  title: "Juara 1 Lomba Programming Nasional",
  description: "Meraih juara 1 dalam kompetisi programming tingkat nasional",
  details: {
    competition_name: "National Programming Contest 2025",
    competition_level: "national",
    rank: 1,
    medal_type: "Gold",
    event_date: new Date("2025-01-15"),
    location: "Jakarta",
    organizer: "Kementerian Pendidikan"
  },
  attachments: [
    {
      file_name: "sertifikat.pdf",
      file_url: "/uploads/sertifikat.pdf",
      file_type: "application/pdf",
      uploaded_at: new Date("2025-01-20T10:00:00Z")
    }
  ],
  tags: ["programming", "competition", "national"],
  points: 100,
  created_at: new Date("2025-01-20T10:00:00Z"),
  updated_at: new Date("2025-01-20T10:00:00Z")
});

// Sample 2: Publication Achievement
db.achievements.insertOne({
  student_id: "550e8400-e29b-41d4-a716-446655440000",
  achievement_type: "publication",
  title: "Paper di International Conference",
  description: "Publikasi paper di konferensi internasional",
  details: {
    publication_type: "conference",
    publication_title: "Machine Learning Applications in Education",
    authors: ["John Doe", "Jane Smith"],
    publisher: "IEEE",
    issn: "1234-5678",
    event_date: new Date("2025-02-10")
  },
  tags: ["publication", "research", "conference"],
  points: 150,
  created_at: new Date("2025-02-15T10:00:00Z"),
  updated_at: new Date("2025-02-15T10:00:00Z")
});

// Sample 3: Organization Achievement
db.achievements.insertOne({
  student_id: "550e8400-e29b-41d4-a716-446655440000",
  achievement_type: "organization",
  title: "Ketua Himpunan Mahasiswa",
  description: "Menjadi ketua himpunan mahasiswa selama 1 tahun",
  details: {
    organization_name: "Himpunan Mahasiswa Teknik Informatika",
    position: "Ketua",
    period: {
      start: new Date("2024-01-01"),
      end: new Date("2024-12-31")
    }
  },
  tags: ["organization", "leadership"],
  points: 80,
  created_at: new Date("2025-01-05T10:00:00Z"),
  updated_at: new Date("2025-01-05T10:00:00Z")
});

// Sample 4: Certification Achievement
db.achievements.insertOne({
  student_id: "550e8400-e29b-41d4-a716-446655440000",
  achievement_type: "certification",
  title: "AWS Certified Solutions Architect",
  description: "Sertifikasi AWS Solutions Architect Associate",
  details: {
    certification_name: "AWS Certified Solutions Architect - Associate",
    issued_by: "Amazon Web Services",
    certification_number: "AWS-123456789",
    valid_until: new Date("2027-12-31"),
    event_date: new Date("2025-03-01")
  },
  attachments: [
    {
      file_name: "aws_certificate.pdf",
      file_url: "/uploads/aws_certificate.pdf",
      file_type: "application/pdf",
      uploaded_at: new Date("2025-03-05T10:00:00Z")
    }
  ],
  tags: ["certification", "aws", "cloud"],
  points: 120,
  created_at: new Date("2025-03-05T10:00:00Z"),
  updated_at: new Date("2025-03-05T10:00:00Z")
});

// Sample 5: Academic Achievement
db.achievements.insertOne({
  student_id: "550e8400-e29b-41d4-a716-446655440000",
  achievement_type: "academic",
  title: "IPK 3.95 Semester 7",
  description: "Mencapai IPK 3.95 pada semester 7",
  details: {
    score: 3.95,
    event_date: new Date("2025-01-31")
  },
  tags: ["academic", "gpa"],
  points: 50,
  created_at: new Date("2025-02-01T10:00:00Z"),
  updated_at: new Date("2025-02-01T10:00:00Z")
});

// Query Examples untuk Navicat

// 1. Find all achievements by student_id
// db.achievements.find({ "student_id": "550e8400-e29b-41d4-a716-446655440000" })

// 2. Find achievements by type
// db.achievements.find({ "achievement_type": "competition" })

// 3. Find achievements with text search
// db.achievements.find({ $text: { $search: "programming" } })

// 4. Find achievements by date range
// db.achievements.find({
//   "created_at": {
//     $gte: new Date("2025-01-01T00:00:00Z"),
//     $lte: new Date("2025-12-31T23:59:59Z")
//   }
// })

// 5. Find achievements with specific tags
// db.achievements.find({ "tags": { $in: ["competition", "national"] } })

// 6. Find achievements by competition level
// db.achievements.find({
//   "achievement_type": "competition",
//   "details.competition_level": "national"
// })

// 7. Aggregate: Count achievements by type
// db.achievements.aggregate([
//   { $group: { _id: "$achievement_type", count: { $sum: 1 } } }
// ])

// 8. Aggregate: Sum points by student
// db.achievements.aggregate([
//   { $group: { _id: "$student_id", total_points: { $sum: "$points" } } },
//   { $sort: { total_points: -1 } }
// ])

// 9. Find achievements with attachments
// db.achievements.find({ "attachments": { $exists: true, $ne: [] } })

// 10. Find achievements by points range
// db.achievements.find({ "points": { $gte: 100, $lte: 200 } })

// Verifikasi data setelah insert
print("Total achievements:", db.achievements.countDocuments({}));
print("Competition achievements:", db.achievements.countDocuments({ achievement_type: "competition" }));
print("Publication achievements:", db.achievements.countDocuments({ achievement_type: "publication" }));

