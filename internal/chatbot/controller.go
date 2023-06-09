package chatbot

import (
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, app *Chatbot) {
	r := fiberApp.Group("/api")

	r.Post("/messages", func(c *fiber.Ctx) error {
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		resp, err := app.ChatAllTheStore(messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message " + err.Error())
		}

		return c.JSON(resp)

	})

	r.Post("/messages/product/:id", func(c *fiber.Ctx) error {
		var messages []chat.Message
		sku := c.Params("id")
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		resp, err := app.ChatWithProduct(sku, messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message " + err.Error())
		}

		return c.JSON(resp)

	})
}
