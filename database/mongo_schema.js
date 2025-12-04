// MongoDB Schema Documentation untuk Navicat
// Database: sppm_2025
// Collection: achievements

// CATATAN PENTING:
// Data achievements TIDAK di-seed melalui file ini.
// Data achievements dibuat secara dinamis melalui API saat mahasiswa membuat prestasi baru.
// File ini hanya berisi dokumentasi schema dan contoh query untuk referensi.

// Gunakan database
use('sppm_2025');

// Struktur Collection: achievements
// {
//   _id: ObjectId,
//   studentId: String (UUID dari PostgreSQL students.id),
//   achievementType: String, // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
//   title: String,
//   description: String,
//   details: Object (field dinamis berdasarkan achievementType),
//   attachments: Array (optional),
//   tags: Array (optional),
//   points: Number,
//   createdAt: Date,
//   updatedAt: Date
// }

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

// 11. Find achievements by studentId (UUID dari PostgreSQL)
// db.achievements.find({ "studentId": "550e8400-e29b-41d4-a716-446655440000" })

// 12. Find achievements by achievementType
// db.achievements.find({ "achievementType": "competition" })

// 13. Find achievements with text search (requires text index)
// db.achievements.find({ $text: { $search: "programming" } })

// 14. Count total achievements
// db.achievements.countDocuments({})

// 15. Count achievements by type
// db.achievements.countDocuments({ "achievementType": "competition" })

