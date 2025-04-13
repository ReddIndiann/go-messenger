package routes

import (
		"github.com/ReddIndiann/go-messanger/controllers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app fiber.Router) {

	middleware.StartCleanupRoutine()
	// Public routes
	app.Post("/api/register", controllers.RegisterUser())
	app.Post("/api/verify", controllers.VerifyMail())
	app.Post("/api/resend-otp", controllers.ResendOTP())
	app.Post("/api/login", controllers.LoginUser())
	app.Post("/api/refresh-token", controllers.RefreshToken())

	// Protected routes - require JWT authentication
	api := app.Group("/api", middleware.JWTAuthMiddleware())


	// User routes
	api.Get("/users", controllers.GetAllUsers())
	api.Get("/users/:id", controllers.GetUserByID())
	api.Put("/users/:id", controllers.UpdateUser())
	api.Delete("/users/:id", controllers.DeleteUser())
	api.Post("/logout", controllers.LogoutUser())
}