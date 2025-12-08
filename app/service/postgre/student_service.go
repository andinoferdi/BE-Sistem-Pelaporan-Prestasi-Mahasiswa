package service

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func GetAllStudentsService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

func GetStudentByIDService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

func GetStudentAchievementsService(c *fiber.Ctx, postgresDB *sql.DB, mongoDB interface{}) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

func UpdateStudentAdvisorService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

