package chatbot

import (
	"log"

	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/gofiber/fiber/v2"
)

// ControllerFactory initializes chatbot routes and handlers
func ControllerFactory(fiberApp *fiber.App, conf Config) {
	apiGroup := fiberApp.Group("/api")

	// Handler for chatting with the store
	apiGroup.Post("/messages", func(c *fiber.Ctx) error {
		bot := New(conf)
		client, err := bot.clientdb.GetById(c.Context(), 1)
		if err != nil {
			log.Printf("Failed to get client: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input format"})
		}

		bot.assistant.AddSystemPrompt(client.SystemPromt)
		resp, err := bot.ChatAllTheStore(messages)
		if err != nil {
			log.Printf("Failed to get message: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
		}

		return c.JSON(resp)
	})

	// Handler for chatting with a specific product
	apiGroup.Post("/products/:id/messages", func(c *fiber.Ctx) error {
		bot := New(conf)
		client, err := bot.clientdb.GetById(c.Context(), 1)
		if err != nil {
			log.Printf("Failed to get client: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input format"})
		}
		sku := c.Params("id")

		bot.assistant.AddSystemPrompt(client.SystemPromt)
		resp, err := bot.ChatWithProduct(sku, messages)
		if err != nil {
			log.Printf("Failed to get message: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
		}

		return c.JSON(resp)
	})
}
