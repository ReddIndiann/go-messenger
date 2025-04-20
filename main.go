package main

import (
	"fmt"
	"os"

	"github.com/ReddIndiann/go-messanger/database"
	"github.com/ReddIndiann/go-messanger/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	fmt.Println("Hey World")
	database.Connect()

	app := fiber.New(fiber.Config{
		AppName: "School App",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	routes.SetupUserRoutes(app.Group("/auth"))
	routes.SetupSchoolRoutes(app.Group("/school"))
	routes.SetupTeacherRoutes(app.Group("/teacher"))
	routes.SetupStudentRoutes(app.Group("/student"))
	routes.SetupSubjectRoutes(app.Group("/subject"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Server Connected to" + port)
	app.Listen(":" + port)
	fmt.Println("Server Connected to" + port)
}
