package service

// #1 proses: import library yang diperlukan untuk context, database, errors, dan time
import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"time"
)

// #2 proses: definisikan interface untuk operasi laporan dan statistik
type IReportService interface {
	GetStatistics(ctx context.Context, userID string, roleID string) (map[string]interface{}, error)
	GetStudentReport(ctx context.Context, studentID string) (map[string]interface{}, error)
	GetLecturerReport(ctx context.Context, lecturerID string) (map[string]interface{}, error)
	GetCurrentStudentReport(ctx context.Context, userID string) (map[string]interface{}, error)
	GetCurrentLecturerReport(ctx context.Context, userID string) (map[string]interface{}, error)
}

// #3 proses: struct service untuk laporan dengan dependency achievement MongoDB, achievement reference PostgreSQL, student, user, dan lecturer repository
type ReportService struct {
	achievementRepo    repositorymongo.IAchievementRepository
	achievementRefRepo repositorypostgre.IAchievementReferenceRepository
	studentRepo        repositorypostgre.IStudentRepository
	userRepo           repositorypostgre.IUserRepository
	lecturerRepo       repositorypostgre.ILecturerRepository
}

// #4 proses: constructor untuk membuat instance ReportService baru
func NewReportService(
	achievementRepo repositorymongo.IAchievementRepository,
	achievementRefRepo repositorypostgre.IAchievementReferenceRepository,
	studentRepo repositorypostgre.IStudentRepository,
	userRepo repositorypostgre.IUserRepository,
	lecturerRepo repositorypostgre.ILecturerRepository,
) IReportService {
	return &ReportService{
		achievementRepo:    achievementRepo,
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		userRepo:           userRepo,
		lecturerRepo:       lecturerRepo,
	}
}

// #5 proses: ambil statistik achievement dengan filtering berdasarkan role user
func (s *ReportService) GetStatistics(ctx context.Context, userID string, roleID string) (map[string]interface{}, error) {
	// #5a proses: ambil role name untuk menentukan filter yang sesuai
	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	var studentIDFilter string
	var advisorIDFilter string

	// #5b proses: set filter berdasarkan role, Mahasiswa filter milik sendiri, Dosen Wali filter mahasiswa bimbingan
	if roleName == "Mahasiswa" {
		student, err := s.studentRepo.GetStudentByUserID(ctx, userID)
		if err != nil {
			return nil, errors.New("error mengambil data student: " + err.Error())
		}
		studentIDFilter = student.ID
	} else if roleName == "Dosen Wali" {
		lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali")
			}
			return nil, err
		}
		advisorIDFilter = lecturer.ID
	} else if roleName != "Admin" {
		return nil, errors.New("akses ditolak. Role tidak memiliki akses untuk melihat statistik")
	}

	// #5c proses: ambil statistik per tipe achievement dari MongoDB
	byType, err := s.achievementRepo.GetAchievementsByType(ctx)
	if err != nil {
		return nil, errors.New("error mengambil statistik per tipe: " + err.Error())
	}

	// #5d proses: filter statistik per tipe berdasarkan student ID atau advisor ID
	if studentIDFilter != "" {
		filteredByType := make(map[string]int)
		achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, studentIDFilter)
		if err == nil {
			for _, achievement := range achievements {
				filteredByType[achievement.AchievementType]++
			}
		}
		byType = filteredByType
	} else if advisorIDFilter != "" {
		students, err := s.studentRepo.GetStudentsByAdvisorID(ctx, advisorIDFilter)
		if err == nil {
			filteredByType := make(map[string]int)
			for _, student := range students {
				achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, student.ID)
				if err == nil {
					for _, achievement := range achievements {
						filteredByType[achievement.AchievementType]++
					}
				}
			}
			byType = filteredByType
		}
	}

	// #5e proses: ambil statistik per periode dari 1 tahun terakhir
	startDate := time.Now().AddDate(-1, 0, 0)
	endDate := time.Now()
	byPeriod, err := s.achievementRefRepo.GetAchievementsByPeriod(ctx, startDate, endDate)
	if err != nil {
		return nil, errors.New("error mengambil statistik per periode: " + err.Error())
	}

	// #5f proses: filter statistik per periode berdasarkan student ID atau advisor ID
	if studentIDFilter != "" {
		references, err := s.achievementRefRepo.GetAchievementReferenceByStudentID(ctx, studentIDFilter)
		if err == nil {
			filteredByPeriod := make(map[string]int)
			for _, ref := range references {
				period := ref.CreatedAt.Format("2006-01")
				filteredByPeriod[period]++
			}
			byPeriod = filteredByPeriod
		}
	} else if advisorIDFilter != "" {
		references, err := s.achievementRefRepo.GetAchievementReferencesByAdvisorID(ctx, advisorIDFilter)
		if err == nil {
			filteredByPeriod := make(map[string]int)
			for _, ref := range references {
				period := ref.CreatedAt.Format("2006-01")
				filteredByPeriod[period]++
			}
			byPeriod = filteredByPeriod
		}
	}

	// #5g proses: ambil top students berdasarkan points dengan filtering
	var topStudents []struct {
		StudentID        string `bson:"_id" json:"student_id"`
		TotalPoints      int    `bson:"totalPoints" json:"total_points"`
		AchievementCount int    `bson:"count" json:"achievement_count"`
	}

	if studentIDFilter != "" {
		// #5h proses: jika filter student, ambil ranking student tersebut saja
		allTopStudents, err := s.achievementRepo.GetTopStudentsByPoints(ctx, 1000)
		if err == nil {
			for _, topStudent := range allTopStudents {
				if topStudent.StudentID == studentIDFilter {
					topStudents = []struct {
						StudentID        string `bson:"_id" json:"student_id"`
						TotalPoints      int    `bson:"totalPoints" json:"total_points"`
						AchievementCount int    `bson:"count" json:"achievement_count"`
					}{topStudent}
					break
				}
			}
		}
	} else if advisorIDFilter != "" {
		// #5i proses: jika filter advisor, ambil ranking mahasiswa bimbingan saja
		students, err := s.studentRepo.GetStudentsByAdvisorID(ctx, advisorIDFilter)
		if err == nil {
			studentIDMap := make(map[string]bool)
			for _, student := range students {
				studentIDMap[student.ID] = true
			}
			allTopStudents, err := s.achievementRepo.GetTopStudentsByPoints(ctx, 1000)
			if err == nil {
				for _, topStudent := range allTopStudents {
					if studentIDMap[topStudent.StudentID] {
						topStudents = append(topStudents, topStudent)
					}
				}
			}
		}
	} else {
		// #5j proses: admin ambil top 10 students
		topStudentsResult, err := s.achievementRepo.GetTopStudentsByPoints(ctx, 10)
		if err != nil {
			return nil, errors.New("error mengambil top mahasiswa: " + err.Error())
		}
		for _, ts := range topStudentsResult {
			topStudents = append(topStudents, struct {
				StudentID        string `bson:"_id" json:"student_id"`
				TotalPoints      int    `bson:"totalPoints" json:"total_points"`
				AchievementCount int    `bson:"count" json:"achievement_count"`
			}{
				StudentID:        ts.StudentID,
				TotalPoints:      ts.TotalPoints,
				AchievementCount: ts.AchievementCount,
			})
		}
	}

	// #5k proses: ambil distribusi level kompetisi dari MongoDB
	competitionLevelDist, err := s.achievementRepo.GetCompetitionLevelDistribution(ctx)
	if err != nil {
		return nil, errors.New("error mengambil distribusi tingkat kompetisi: " + err.Error())
	}

	// #5l proses: filter distribusi level kompetisi berdasarkan student ID atau advisor ID
	if studentIDFilter != "" {
		achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, studentIDFilter)
		if err == nil {
			filteredCompetitionDist := make(map[string]int)
			for _, achievement := range achievements {
				if achievement.AchievementType == "competition" && achievement.Details.CompetitionLevel != nil {
					level := *achievement.Details.CompetitionLevel
					filteredCompetitionDist[level]++
				}
			}
			competitionLevelDist = filteredCompetitionDist
		}
	} else if advisorIDFilter != "" {
		students, err := s.studentRepo.GetStudentsByAdvisorID(ctx, advisorIDFilter)
		if err == nil {
			filteredCompetitionDist := make(map[string]int)
			for _, student := range students {
				achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, student.ID)
				if err == nil {
					for _, achievement := range achievements {
						if achievement.AchievementType == "competition" && achievement.Details.CompetitionLevel != nil {
							level := *achievement.Details.CompetitionLevel
							filteredCompetitionDist[level]++
						}
					}
				}
			}
			competitionLevelDist = filteredCompetitionDist
		}
	}

	// #5m proses: enrich top students dengan nama student dari database
	topStudentsWithNames := make([]map[string]interface{}, 0)
	for _, topStudent := range topStudents {
		student, err := s.studentRepo.GetStudentByID(ctx, topStudent.StudentID)
		if err != nil {
			continue
		}

		user, err := s.userRepo.FindUserByID(ctx, student.UserID)
		if err != nil {
			continue
		}

		topStudentsWithNames = append(topStudentsWithNames, map[string]interface{}{
			"student_id":        topStudent.StudentID,
			"student_name":      user.FullName,
			"total_points":      topStudent.TotalPoints,
			"achievement_count": topStudent.AchievementCount,
		})
	}

	// #5n proses: build response dengan semua statistik
	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"byType":                       byType,
			"byPeriod":                     byPeriod,
			"topStudents":                  topStudentsWithNames,
			"competitionLevelDistribution": competitionLevelDist,
		},
	}, nil
}

// #6 proses: ambil laporan lengkap untuk student tertentu
func (s *ReportService) GetStudentReport(ctx context.Context, studentID string) (map[string]interface{}, error) {
	// #6a proses: validasi student ID tidak kosong
	if studentID == "" {
		return nil, errors.New("student ID wajib diisi")
	}

	// #6b proses: ambil student dan user data
	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student tidak ditemukan")
		}
		return nil, errors.New("error mengambil data student: " + err.Error())
	}

	user, err := s.userRepo.FindUserByID(ctx, student.UserID)
	if err != nil {
		return nil, errors.New("error mengambil data user: " + err.Error())
	}

	// #6c proses: ambil achievements dari MongoDB dan references dari PostgreSQL
	achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, studentID)
	if err != nil {
		return nil, errors.New("error mengambil achievements: " + err.Error())
	}

	references, err := s.achievementRefRepo.GetAchievementReferenceByStudentID(ctx, studentID)
	if err != nil {
		return nil, errors.New("error mengambil achievement references: " + err.Error())
	}

	// #6d proses: buat map reference untuk lookup cepat
	referenceMap := make(map[string]modelpostgre.AchievementReference)
	for _, ref := range references {
		referenceMap[ref.MongoAchievementID] = ref
	}

	// #6e proses: hitung statistik dan build achievement details
	totalPoints := 0
	verifiedCount := 0
	byType := make(map[string]int)
	achievementDetails := make([]map[string]interface{}, 0)

	for _, achievement := range achievements {
		ref, exists := referenceMap[achievement.ID.Hex()]
		if !exists {
			continue
		}

		totalPoints += achievement.Points
		byType[achievement.AchievementType]++

		if ref.Status == "verified" {
			verifiedCount++
		}

		achievementDetails = append(achievementDetails, map[string]interface{}{
			"id":              achievement.ID.Hex(),
			"title":           achievement.Title,
			"achievementType": achievement.AchievementType,
			"status":          ref.Status,
			"points":          achievement.Points,
			"createdAt":       achievement.CreatedAt.Format(time.RFC3339),
		})
	}

	// #6f proses: build response dengan data student, statistik, dan list achievements
	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"student": map[string]interface{}{
				"id":         student.ID,
				"name":       user.FullName,
				"student_id": student.StudentID,
			},
			"statistics": map[string]interface{}{
				"total_achievements": len(achievementDetails),
				"total_points":       totalPoints,
				"verified_count":     verifiedCount,
				"by_type":            byType,
			},
			"achievements": achievementDetails,
		},
	}, nil
}

// #7 proses: ambil laporan lengkap untuk lecturer dengan statistik mahasiswa bimbingan
func (s *ReportService) GetLecturerReport(ctx context.Context, lecturerID string) (map[string]interface{}, error) {
	// #7a proses: validasi lecturer ID tidak kosong
	if lecturerID == "" {
		return nil, errors.New("lecturer ID wajib diisi")
	}

	// #7b proses: ambil lecturer dan user data
	lecturer, err := s.lecturerRepo.GetLecturerByID(ctx, lecturerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lecturer tidak ditemukan")
		}
		return nil, errors.New("error mengambil data lecturer: " + err.Error())
	}

	user, err := s.userRepo.FindUserByID(ctx, lecturer.UserID)
	if err != nil {
		return nil, errors.New("error mengambil data user: " + err.Error())
	}

	// #7c proses: ambil semua mahasiswa bimbingan
	students, err := s.studentRepo.GetStudentsByAdvisorID(ctx, lecturerID)
	if err != nil {
		return nil, errors.New("error mengambil mahasiswa bimbingan: " + err.Error())
	}

	// #7d proses: hitung statistik dari semua mahasiswa bimbingan
	totalPoints := 0
	totalAchievements := 0
	byType := make(map[string]int)
	studentStats := make(map[string]struct {
		TotalPoints      int
		AchievementCount int
	})

	// #7e proses: loop semua mahasiswa bimbingan dan hitung statistik per student
	for _, student := range students {
		achievements, err := s.achievementRepo.GetAchievementsByStudentID(ctx, student.ID)
		if err != nil {
			continue
		}

		references, err := s.achievementRefRepo.GetAchievementReferenceByStudentID(ctx, student.ID)
		if err != nil {
			continue
		}

		// #7f proses: buat map reference, filter yang status bukan deleted
		referenceMap := make(map[string]modelpostgre.AchievementReference)
		for _, ref := range references {
			if ref.Status != "deleted" {
				referenceMap[ref.MongoAchievementID] = ref
			}
		}

		studentPoints := 0
		studentAchievementCount := 0

		// #7g proses: hitung points dan count per student, lalu agregasi ke total
		for _, achievement := range achievements {
			_, exists := referenceMap[achievement.ID.Hex()]
			if !exists {
				continue
			}

			totalPoints += achievement.Points
			studentPoints += achievement.Points
			totalAchievements++
			studentAchievementCount++
			byType[achievement.AchievementType]++
		}

		studentStats[student.ID] = struct {
			TotalPoints      int
			AchievementCount int
		}{
			TotalPoints:      studentPoints,
			AchievementCount: studentAchievementCount,
		}
	}

	// #7h proses: buat list top advisees dari student stats
	type TopAdvisee struct {
		StudentID        string
		TotalPoints      int
		AchievementCount int
	}

	topAdviseesList := make([]TopAdvisee, 0)
	for studentID, stats := range studentStats {
		topAdviseesList = append(topAdviseesList, TopAdvisee{
			StudentID:        studentID,
			TotalPoints:      stats.TotalPoints,
			AchievementCount: stats.AchievementCount,
		})
	}

	// #7i proses: sort top advisees berdasarkan total points descending
	for i := 0; i < len(topAdviseesList)-1; i++ {
		for j := i + 1; j < len(topAdviseesList); j++ {
			if topAdviseesList[i].TotalPoints < topAdviseesList[j].TotalPoints {
				topAdviseesList[i], topAdviseesList[j] = topAdviseesList[j], topAdviseesList[i]
			}
		}
	}

	// #7j proses: ambil top 10 advisees saja
	if len(topAdviseesList) > 10 {
		topAdviseesList = topAdviseesList[:10]
	}

	// #7k proses: enrich top advisees dengan nama student
	topAdviseesWithNames := make([]map[string]interface{}, 0)
	for _, topAdvisee := range topAdviseesList {
		student, err := s.studentRepo.GetStudentByID(ctx, topAdvisee.StudentID)
		if err != nil {
			continue
		}

		studentUser, err := s.userRepo.FindUserByID(ctx, student.UserID)
		if err != nil {
			continue
		}

		topAdviseesWithNames = append(topAdviseesWithNames, map[string]interface{}{
			"student_id":        topAdvisee.StudentID,
			"student_name":      studentUser.FullName,
			"total_points":      topAdvisee.TotalPoints,
			"achievement_count": topAdvisee.AchievementCount,
		})
	}

	// #7l proses: build response dengan data lecturer, statistik, dan top advisees
	return map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"lecturer": map[string]interface{}{
				"id":          lecturer.ID,
				"name":        user.FullName,
				"lecturer_id": lecturer.LecturerID,
				"department":  lecturer.Department,
			},
			"statistics": map[string]interface{}{
				"total_advisees":     len(students),
				"total_achievements": totalAchievements,
				"total_points":       totalPoints,
				"by_type":            byType,
			},
			"topAdvisees": topAdviseesWithNames,
		},
	}, nil
}

// #8 proses: ambil laporan student untuk user yang sedang login
func (s *ReportService) GetCurrentStudentReport(ctx context.Context, userID string) (map[string]interface{}, error) {
	// #8a proses: validasi user ID tidak kosong
	if userID == "" {
		return nil, errors.New("user ID wajib diisi")
	}

	// #8b proses: ambil student berdasarkan user ID, lalu panggil GetStudentReport
	student, err := s.studentRepo.GetStudentByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student profile tidak ditemukan untuk user ini")
		}
		return nil, errors.New("error mengambil data student: " + err.Error())
	}

	return s.GetStudentReport(ctx, student.ID)
}

// #9 proses: ambil laporan lecturer untuk user yang sedang login
func (s *ReportService) GetCurrentLecturerReport(ctx context.Context, userID string) (map[string]interface{}, error) {
	// #9a proses: validasi user ID tidak kosong
	if userID == "" {
		return nil, errors.New("user ID wajib diisi")
	}

	// #9b proses: ambil lecturer berdasarkan user ID, lalu panggil GetLecturerReport
	lecturer, err := s.userRepo.GetLecturerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lecturer profile tidak ditemukan untuk user ini")
		}
		return nil, errors.New("error mengambil data lecturer: " + err.Error())
	}

	return s.GetLecturerReport(ctx, lecturer.ID)
}
