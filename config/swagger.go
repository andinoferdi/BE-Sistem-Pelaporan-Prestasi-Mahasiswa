package config

import (
	"encoding/json"
	"strings"

	_ "sistem-pelaporan-prestasi-mahasiswa/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/swaggo/swag"
)

func CustomSwaggerHandler() fiber.Handler {
	handler := fiberSwagger.FiberWrapHandler(
		fiberSwagger.DeepLinking(true),
		fiberSwagger.DocExpansion("none"),
	)

	return func(c *fiber.Ctx) error {
		path := c.Path()

		if strings.HasSuffix(path, "swagger.json") || strings.HasSuffix(path, "doc.json") {
			swaggerInfo := swag.GetSwagger("swagger")
			if swaggerInfo != nil {
				specBytes := []byte(swaggerInfo.ReadDoc())

				var specJSON map[string]interface{}
				if err := json.Unmarshal(specBytes, &specJSON); err == nil {
					tags := []map[string]interface{}{
						{"name": "Authentication", "description": "Authentication endpoints"},
						{"name": "Users", "description": "User management endpoints"},
						{"name": "Achievements", "description": "Achievement management endpoints"},
						{"name": "Students", "description": "Student management endpoints"},
						{"name": "Lecturers", "description": "Lecturer management endpoints"},
						{"name": "Reports", "description": "Reporting and analytics endpoints"},
						{"name": "Notifications", "description": "Notification management endpoints"},
						{"name": "System", "description": "System endpoints"},
					}
					specJSON["tags"] = tags

					if jsonBytes, err := json.Marshal(specJSON); err == nil {
						return c.Type("application/json").Send(jsonBytes)
					}
				}
			}
		}

		if path == "/swagger" || path == "/swagger/" {
			return c.Redirect("/swagger/index.html", 301)
		}

		return handler(c)
	}
}
