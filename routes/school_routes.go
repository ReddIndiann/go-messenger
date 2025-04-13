package routes

import (
	"github.com/ReddIndiann/go-messanger/controllers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSchoolRoutes(app fiber.Router) {

	middleware.StartCleanupRoutine()

	api := app.Group("/api", middleware.JWTAuthMiddleware())
	api.Post("/register", controllers.RegisterSchool())
	api.Get("/", controllers.GetAllSchool())
	api.Get("/:id", controllers.GetSchoolByID())
	api.Put("/:id", controllers.UpdateSchool())
	api.Delete("/:id", controllers.DeleteSchool())

}
