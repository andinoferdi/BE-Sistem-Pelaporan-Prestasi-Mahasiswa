package service

import (
	"database/sql"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"sistem-pelaporan-prestasi-mahasiswa/helper"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAchievementService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya mahasiswa yang dapat membuat prestasi.",
			},
		})
	}

	studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Data mahasiswa tidak ditemukan. Pastikan user memiliki profil mahasiswa.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data mahasiswa dari database. Detail: " + err.Error(),
			},
		})
	}

	var req modelmongo.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
			},
		})
	}

	if req.AchievementType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Achievement type wajib diisi.",
			},
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Title wajib diisi.",
			},
		})
	}

	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Description wajib diisi.",
			},
		})
	}

	validTypes := map[string]bool{
		"academic":      true,
		"competition":   true,
		"organization":  true,
		"publication":   true,
		"certification": true,
		"other":         true,
	}

	if !validTypes[req.AchievementType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Achievement type tidak valid. Gunakan: academic, competition, organization, publication, certification, atau other.",
			},
		})
	}

	req.StudentID = studentID

	achievement := modelmongo.Achievement{
		StudentID:       req.StudentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Attachments:     req.Attachments,
		Tags:            req.Tags,
		Points:          req.Points,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	createdAchievement, err := repositorymongo.CreateAchievement(mongoDB, achievement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menyimpan prestasi ke database. Detail: " + err.Error(),
			},
		})
	}

	refReq := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          studentID,
		MongoAchievementID: createdAchievement.ID.Hex(),
		Status:             modelpostgre.AchievementStatusDraft,
	}

	_, err = repositorypostgre.CreateAchievementReference(postgresDB, refReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error membuat reference prestasi. Detail: " + err.Error(),
			},
		})
	}

	response := modelmongo.CreateAchievementResponse{
		Status: "success",
		Data:   *createdAchievement,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func SubmitAchievementService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya mahasiswa yang dapat submit prestasi.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Prestasi hanya dapat di-submit jika status adalah draft.",
			},
		})
	}

	studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data mahasiswa. Detail: " + err.Error(),
			},
		})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat submit prestasi milik Anda sendiri.",
			},
		})
	}

	now := time.Now()
	err = repositorypostgre.UpdateAchievementReferenceStatus(postgresDB, ref.ID, modelpostgre.AchievementStatusSubmitted, &now)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengupdate status prestasi. Detail: " + err.Error(),
			},
		})
	}

	updatedRef, err := repositorypostgre.GetAchievementReferenceByID(postgresDB, ref.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi yang diupdate. Detail: " + err.Error(),
			},
		})
	}

	response := modelpostgre.UpdateAchievementReferenceResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func VerifyAchievementService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya dosen wali yang dapat memverifikasi prestasi.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Prestasi hanya dapat diverifikasi jika status adalah submitted.",
			},
		})
	}

	lecturer, err := repositorypostgre.GetLecturerByUserID(postgresDB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data dosen wali. Detail: " + err.Error(),
			},
		})
	}

	var student modelpostgre.Student
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id = $1`
	err = postgresDB.QueryRow(query, ref.StudentID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data student. Detail: " + err.Error(),
			},
		})
	}

	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat memverifikasi prestasi mahasiswa bimbingan Anda.",
			},
		})
	}

	err = repositorypostgre.UpdateAchievementReferenceVerify(postgresDB, ref.ID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error memverifikasi prestasi. Detail: " + err.Error(),
			},
		})
	}

	updatedRef, err := repositorypostgre.GetAchievementReferenceByID(postgresDB, ref.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi yang diupdate. Detail: " + err.Error(),
			},
		})
	}

	response := modelpostgre.VerifyAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func RejectAchievementService(c *fiber.Ctx, postgresDB *sql.DB) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya dosen wali yang dapat menolak prestasi.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	var req modelpostgre.RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
			},
		})
	}

	if req.RejectionNote == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Rejection note wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if ref.Status != modelpostgre.AchievementStatusSubmitted {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Prestasi hanya dapat ditolak jika status adalah submitted.",
			},
		})
	}

	lecturer, err := repositorypostgre.GetLecturerByUserID(postgresDB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data dosen wali. Detail: " + err.Error(),
			},
		})
	}

	var student modelpostgre.Student
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id = $1`
	err = postgresDB.QueryRow(query, ref.StudentID).Scan(
		&student.ID, &student.UserID, &student.StudentID,
		&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data student. Detail: " + err.Error(),
			},
		})
	}

	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat menolak prestasi mahasiswa bimbingan Anda.",
			},
		})
	}

	err = repositorypostgre.UpdateAchievementReferenceReject(postgresDB, ref.ID, userID, req.RejectionNote)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menolak prestasi. Detail: " + err.Error(),
			},
		})
	}

	updatedRef, err := repositorypostgre.GetAchievementReferenceByID(postgresDB, ref.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi yang diupdate. Detail: " + err.Error(),
			},
		})
	}

	response := modelpostgre.RejectAchievementResponse{
		Status: "success",
		Data:   *updatedRef,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteAchievementService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya mahasiswa yang dapat menghapus prestasi.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Prestasi hanya dapat dihapus jika status adalah draft.",
			},
		})
	}

	studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data mahasiswa. Detail: " + err.Error(),
			},
		})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat menghapus prestasi milik Anda sendiri.",
			},
		})
	}

	err = repositorymongo.DeleteAchievement(mongoDB, mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menghapus prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	err = repositorypostgre.UpdateAchievementReferenceStatus(postgresDB, ref.ID, modelpostgre.AchievementStatusDeleted, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengupdate status prestasi menjadi deleted. Detail: " + err.Error(),
			},
		})
	}

	response := modelmongo.DeleteAchievementResponse{
		Status: "success",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func UploadFileService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName == "Mahasiswa" {
		studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil data mahasiswa. Detail: " + err.Error(),
				},
			})
		}

		if ref.StudentID != studentID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Akses ditolak. Anda hanya dapat mengupload file untuk prestasi milik Anda sendiri.",
				},
			})
		}

		if ref.Status != modelpostgre.AchievementStatusDraft {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "File hanya dapat diupload jika status prestasi adalah draft.",
				},
			})
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "File tidak ditemukan. Pastikan field name adalah 'file'.",
			},
		})
	}

	allowedTypes := map[string]bool{
		"application/pdf":       true,
		"image/jpeg":            true,
		"image/jpg":            true,
		"image/png":             true,
		"application/msword":    true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".doc":  true,
		".docx": true,
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[fileExt] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Tipe file tidak diizinkan. Gunakan PDF, JPG, PNG, DOC, atau DOCX.",
			},
		})
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Ukuran file terlalu besar. Maksimal 10MB.",
			},
		})
	}

	fileType := mime.TypeByExtension(fileExt)
	if fileType == "" {
		fileType = file.Header.Get("Content-Type")
	}

	if !allowedTypes[fileType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Tipe file tidak diizinkan. Gunakan PDF, JPG, PNG, DOC, atau DOCX.",
			},
		})
	}

	timestamp := time.Now().Unix()
	safeFilename := fmt.Sprintf("%d-%s", timestamp, file.Filename)
	safeFilename = strings.ReplaceAll(safeFilename, " ", "_")

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error membuat folder uploads. Detail: " + err.Error(),
			},
		})
	}

	filePath := filepath.Join(uploadDir, safeFilename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menyimpan file. Detail: " + err.Error(),
			},
		})
	}

	fileURL := fmt.Sprintf("/uploads/%s", safeFilename)
	attachment := modelmongo.Attachment{
		FileName:   file.Filename,
		FileURL:    fileURL,
		FileType:   fileType,
		UploadedAt: time.Now(),
	}

	_, err = repositorymongo.AddAttachmentToAchievement(mongoDB, mongoID, attachment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error menambahkan attachment ke prestasi. Detail: " + err.Error(),
			},
		})
	}

	response := fiber.Map{
		"status": "success",
		"data":   attachment,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetAchievementsService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	page := helper.GetQueryInt(c, "page", 1)
	limit := helper.GetQueryInt(c, "limit", 10)
	page, limit = helper.ValidatePagination(page, limit)

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	var references []modelpostgre.AchievementReference
	var total int

	if roleName == "Mahasiswa" {
		student, err := repositorypostgre.GetStudentByUserID(postgresDB, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil data student. Detail: " + err.Error(),
				},
			})
		}

		references, total, err = repositorypostgre.GetAchievementReferenceByStudentIDPaginated(postgresDB, student.ID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil achievement references. Detail: " + err.Error(),
				},
			})
		}
	} else if roleName == "Dosen Wali" {
		lecturer, err := repositorypostgre.GetLecturerByUserID(postgresDB, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"data": fiber.Map{
						"message": "Data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali.",
					},
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil data dosen wali. Detail: " + err.Error(),
				},
			})
		}

		references, total, err = repositorypostgre.GetAchievementReferencesByAdvisorIDPaginated(postgresDB, lecturer.ID, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil achievement references. Detail: " + err.Error(),
				},
			})
		}
	} else if roleName == "Admin" {
		references, total, err = repositorypostgre.GetAllAchievementReferencesPaginated(postgresDB, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil achievement references. Detail: " + err.Error(),
				},
			})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Role tidak memiliki akses untuk melihat prestasi.",
			},
		})
	}

	if len(references) == 0 {
		totalPages := 0
		if total > 0 {
			totalPages = (total + limit - 1) / limit
		}
		response := fiber.Map{
			"status": "success",
			"data":   []modelmongo.Achievement{},
			"pagination": fiber.Map{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"total_pages": totalPages,
			},
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}

	var mongoIDs []string
	for _, ref := range references {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}

	achievements, err := repositorymongo.GetAchievementsByIDs(mongoDB, mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil achievements dari MongoDB. Detail: " + err.Error(),
			},
		})
	}

	referenceMap := make(map[string]modelpostgre.AchievementReference)
	for _, ref := range references {
		referenceMap[ref.MongoAchievementID] = ref
	}

	var result []fiber.Map
	for _, achievement := range achievements {
		ref, exists := referenceMap[achievement.ID.Hex()]
		if !exists {
			continue
		}

		result = append(result, fiber.Map{
			"id":              achievement.ID.Hex(),
			"studentId":        achievement.StudentID,
			"achievementType":  achievement.AchievementType,
			"title":            achievement.Title,
			"description":      achievement.Description,
			"details":          achievement.Details,
			"attachments":      achievement.Attachments,
			"tags":             achievement.Tags,
			"points":           achievement.Points,
			"createdAt":        achievement.CreatedAt.Format(time.RFC3339),
			"updatedAt":        achievement.UpdatedAt.Format(time.RFC3339),
			"status":           ref.Status,
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	response := fiber.Map{
		"status": "success",
		"data":   result,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": totalPages,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetAchievementByIDService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if roleName == "Mahasiswa" {
	studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data mahasiswa. Detail: " + err.Error(),
			},
		})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat melihat prestasi milik Anda sendiri.",
				},
			})
		}
	} else if roleName == "Dosen Wali" {
		lecturer, err := repositorypostgre.GetLecturerByUserID(postgresDB, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"data": fiber.Map{
						"message": "Data dosen wali tidak ditemukan. Pastikan user memiliki profil dosen wali.",
					},
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil data dosen wali. Detail: " + err.Error(),
				},
			})
		}

		var student modelpostgre.Student
		query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id = $1`
		err = postgresDB.QueryRow(query, ref.StudentID).Scan(
			&student.ID, &student.UserID, &student.StudentID,
			&student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
			&student.CreatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Error mengambil data student. Detail: " + err.Error(),
				},
			})
		}

		if student.AdvisorID != lecturer.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Akses ditolak. Anda hanya dapat melihat prestasi mahasiswa bimbingan Anda.",
				},
			})
		}
	} else if roleName != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Role tidak memiliki akses untuk melihat prestasi.",
			},
		})
	}

	achievement, err := repositorymongo.GetAchievementByID(mongoDB, mongoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil achievement dari database. Detail: " + err.Error(),
			},
		})
	}

	result := fiber.Map{
		"id":              achievement.ID.Hex(),
		"studentId":       achievement.StudentID,
		"achievementType": achievement.AchievementType,
		"title":           achievement.Title,
		"description":     achievement.Description,
		"details":         achievement.Details,
		"attachments":     achievement.Attachments,
		"tags":            achievement.Tags,
		"points":          achievement.Points,
		"createdAt":       achievement.CreatedAt.Format(time.RFC3339),
		"updatedAt":       achievement.UpdatedAt.Format(time.RFC3339),
		"status":          ref.Status,
	}

	if ref.SubmittedAt != nil {
		result["submitted_at"] = ref.SubmittedAt.Format(time.RFC3339)
	}

	if ref.VerifiedAt != nil {
		result["verified_at"] = ref.VerifiedAt.Format(time.RFC3339)
	}

	if ref.VerifiedBy != nil {
		verifiedByUser, err := repositorypostgre.GetUserByID(postgresDB, *ref.VerifiedBy)
		if err == nil && verifiedByUser != nil {
			result["verified_by"] = verifiedByUser.FullName
		} else {
			result["verified_by"] = *ref.VerifiedBy
		}
	}

	if ref.RejectionNote != nil {
		result["rejection_note"] = *ref.RejectionNote
	}

	responseData := fiber.Map{
		"status": "success",
		"data":   result,
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func UpdateAchievementService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB *mongo.Database) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "User ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleID, ok := c.Locals("role_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Role ID tidak ditemukan. Silakan login ulang.",
			},
		})
	}

	roleName, err := repositorypostgre.GetRoleName(postgresDB, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil role name. Detail: " + err.Error(),
			},
		})
	}

	if roleName != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Hanya mahasiswa yang dapat mengupdate prestasi.",
			},
		})
	}

	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "ID prestasi wajib diisi.",
			},
		})
	}

	ref, err := repositorypostgre.GetAchievementReferenceByMongoID(postgresDB, mongoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Prestasi tidak ditemukan.",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data prestasi dari database. Detail: " + err.Error(),
			},
		})
	}

	if ref.Status != modelpostgre.AchievementStatusDraft {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Prestasi hanya dapat diupdate jika status adalah draft.",
			},
		})
	}

	studentID, err := repositorypostgre.GetStudentIDByUserID(postgresDB, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil data mahasiswa. Detail: " + err.Error(),
			},
		})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Akses ditolak. Anda hanya dapat mengupdate prestasi milik Anda sendiri.",
			},
		})
	}

	var req modelmongo.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Format request body tidak valid. Pastikan JSON format benar. Detail: " + err.Error(),
			},
		})
	}

	if req.AchievementType != "" {
		validTypes := map[string]bool{
			"academic":      true,
			"competition":   true,
			"organization":  true,
			"publication":   true,
			"certification": true,
			"other":         true,
		}

		if !validTypes[req.AchievementType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"data": fiber.Map{
					"message": "Achievement type tidak valid. Gunakan: academic, competition, organization, publication, certification, atau other.",
				},
			})
		}
	}

	updatedAchievement, err := repositorymongo.UpdateAchievement(mongoDB, mongoID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengupdate prestasi di database. Detail: " + err.Error(),
			},
		})
	}

	result := fiber.Map{
		"id":              updatedAchievement.ID.Hex(),
		"studentId":       updatedAchievement.StudentID,
		"achievementType": updatedAchievement.AchievementType,
		"title":           updatedAchievement.Title,
		"description":     updatedAchievement.Description,
		"details":         updatedAchievement.Details,
		"attachments":     updatedAchievement.Attachments,
		"tags":            updatedAchievement.Tags,
		"points":          updatedAchievement.Points,
		"createdAt":       updatedAchievement.CreatedAt.Format(time.RFC3339),
		"updatedAt":       updatedAchievement.UpdatedAt.Format(time.RFC3339),
		"status":          ref.Status,
	}

	responseData := fiber.Map{
		"status": "success",
		"data":   result,
	}

	return c.Status(fiber.StatusOK).JSON(responseData)
}

func GetAchievementStatsService(c *fiber.Ctx, postgresDB *sql.DB) error {
	total, verified, err := repositorypostgre.GetAchievementStats(postgresDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data": fiber.Map{
				"message": "Error mengambil statistik prestasi. Detail: " + err.Error(),
			},
		})
	}

	percentage := 0
	if total > 0 {
		percentage = int((float64(verified) / float64(total)) * 100)
	}

	response := fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total":      total,
			"verified":   verified,
			"percentage": percentage,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetAchievementHistoryService(c *fiber.Ctx, postgresDB *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

