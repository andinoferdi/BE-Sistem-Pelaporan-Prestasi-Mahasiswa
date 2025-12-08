package service

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func GetStatisticsService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

func GetStudentReportService(c *fiber.Ctx, db *sql.DB) error {
	return c.Status(501).JSON(fiber.Map{
		"status": "error",
		"data": fiber.Map{
			"message": "Fitur ini belum diimplementasikan.",
		},
	})
}

