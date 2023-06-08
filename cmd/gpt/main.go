package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/assistand"
	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	conn, err := database.NewConnection("localhost", 5432, "myuser", "mypassword", "mydb")
	if err != nil {
		log.Fatal(err)
	}

	assistandApp := assistand.New(conn, openia.New(os.Getenv("OPENAI_API_KEY")))
	assistand.ControllerFactory(app, assistandApp)

	log.Fatal(app.Listen(":3000"))

}
