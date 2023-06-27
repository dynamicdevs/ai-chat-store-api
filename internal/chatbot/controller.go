package chatbot

import (
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/assistant"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, conn *database.Connection) {
	r := fiberApp.Group("/api")

	r.Post("/messages", func(c *fiber.Ctx) error {
		assistant := assistant.New(openia.New(os.Getenv("OPENAI_API_KEY")))
		bot := New(conn, assistant)
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		systemPrompt := `You are an ecommerce asystenat of ABCDIN that is goig to help the customer aswering
			their question about products, maybe comparing some products, give charactristic, etc.
			If someone ask something not relaited to retail or the store, aswer with sorry i cant help you.
			Dont mention any other product rather the one in the stock
	`
		bot.assistant.AddSystemPrompt(systemPrompt)
		resp, err := bot.ChatAllTheStore(messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message " + err.Error())
		}

		return c.JSON(resp)

	})

	r.Post("/messages/product/:id", func(c *fiber.Ctx) error {
		assistant := assistant.New(openia.New(os.Getenv("OPENAI_API_KEY")))
		bot := New(conn, assistant)
		var messages []chat.Message
		sku := c.Params("id")
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		resp, err := bot.ChatWithProduct(sku, messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message " + err.Error())
		}

		return c.JSON(resp)

	})
}
