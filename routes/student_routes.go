package routes

import (
	"github.com/ReddIndiann/go-messanger/controllers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(app fiber.Router) {
	api := app.Group("/api", middleware.JWTAuthMiddleware())

	// Student routes
	api.Post("/register", controllers.RegisterStudent())
	api.Get("/", controllers.ListStudents())
	api.Get("/:id", controllers.GetStudent())
	api.Put("/:id", controllers.UpdateStudent())
	api.Delete("/:id", controllers.DeleteStudent())
}
