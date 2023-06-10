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

	conn, err := database.NewConnection(os.Getenv("DB_URI"), 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))
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

	log.Fatal(app.Listen(":3000"))

}
