package main

import (
	"BlanjaInitialBackEndGOFiber/src/configs"
	"BlanjaInitialBackEndGOFiber/src/helpers"
	"BlanjaInitialBackEndGOFiber/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	}))
	configs.InitDB()
	helpers.Migration()
	routes.Router(app)
	app.Listen(":3000")
}
