package chatbot

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, conf Config) {
	r := fiberApp.Group("/api")

	r.Post("/messages", func(c *fiber.Ctx) error {
		ctx := context.Background()
		bot := New(conf)
		client, err := bot.clientdb.GetById(ctx, 1)
		if err != nil {
			return c.Status(500).SendString("Failed to get client: " + err.Error())
		}
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		bot.assistant.AddSystemPrompt(client.SystemPromt)
		resp, err := bot.ChatAllTheStore(messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message " + err.Error())
		}

		return c.JSON(resp)

	})

	r.Post("/messages/product/:id", func(c *fiber.Ctx) error {
		bot := New(conf)
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
