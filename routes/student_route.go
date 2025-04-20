package routes

import (
	"github.com/ReddIndiann/go-messanger/controllers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/gofiber/fiber/v2"
)

func StudentRoute(app *fiber.App) {
	api := app.Group("/api", middleware.JWTAuthMiddleware())

	// Teacher routes
	api.Post("/register", controllers.RegisterTeacher())
	api.Get("/", controllers.GetAllTeachers())
	api.Get("/:id", controllers.GetTeacherByID())
	api.Put("/:id", controllers.UpdateTeacher())
	api.Delete("/:id", controllers.DeleteTeacher())
}
