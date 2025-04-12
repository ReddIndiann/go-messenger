package main

import (
	"fmt"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ReddIndiann/go-messanger/database"
	"github.com/ReddIndiann/go-messanger/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	fmt.Println("Hey World")
database.Connect()

app := fiber.New(fiber.Config{
	AppName: "Diary App",
})

app.Use(cors.New(cors.Config{
	AllowOrigins: "*",
	AllowMethods: "GET,POST,PUT,DELETE",
}))

routes.SetupRoutes(app.Group("/auth"))
port := os.Getenv("PORT")
if port ==""{
port = "3000"

}
fmt.Println("Server Connected to"+ port)
app.Listen(":" + port)
fmt.Println("Server Connected to"+ port)
}