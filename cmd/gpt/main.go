package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/chatbot"
	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/assistant"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	conn, err := database.NewConnection("chat-store-server.postgres.database.azure.com", 5432, "chatstoreadmin", "Ch4tSt0R34dm1n", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	systemPrompt := `You are an ecommerce asystenat of ABCDIN that is goig to help the customer aswering
			their question about products, maybe comparing some products, give charactristic, etc.
			If someone ask something not relaited to retail or the store, aswer with sorry i cant help you.
			Dont mention any other product rather the one in the stock
	`
	assistantApp := assistant.New(systemPrompt, openia.New(os.Getenv("OPENAI_API_KEY")))
	chatbotApp := chatbot.New(conn, assistantApp)
	chatbot.ControllerFactory(app, chatbotApp)
	// Serve the Swagger UI
	app.Static("/docs", "./dist")
	// Serve the Swagger YAML file
	app.Use("/swagger.yml", func(c *fiber.Ctx) error {
		return c.SendFile("./swagger.yml")
	})

	log.Fatal(app.Listen(":80"))

}
