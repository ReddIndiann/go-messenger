package routes

import (
	"github.com/ReddIndiann/go-messanger/controllers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSubjectRoutes(app fiber.Router) {
	api := app.Group("/api", middleware.JWTAuthMiddleware())

	// Subject routes
	api.Post("/register", controllers.RegisterSubject())
	api.Get("/", controllers.ListSubjects())
	api.Get("/:id", controllers.GetSubject())
	api.Put("/:id", controllers.UpdateSubject())
	api.Delete("/:id", controllers.DeleteSubject())
}
