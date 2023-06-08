package assistand

import (
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, app *Assistand) {
	r := fiberApp.Group("/api")

	r.Post("/messages", func(c *fiber.Ctx) error {
		var messages []chat.Message
		if err := c.BodyParser(&messages); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}

		resp, err := app.HelpWithEveryThing(messages)
		if err != nil {
			return c.Status(400).SendString("Failed to get message")
		}

		messages = append(messages, chat.Message{Role: "assistant", Content: resp})
		return c.JSON(messages)

	})

}
