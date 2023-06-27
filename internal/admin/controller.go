package admin

import (
	"context"

	"github.com/Abraxas-365/commerce-chat/pkg/client"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, conf Config) {
	r := fiberApp.Group("/api/admin")
	admin := New(conf)

	r.Post("/client", func(c *fiber.Ctx) error {
		ctx := context.Background()
		var newClient client.Client
		if err := c.BodyParser(&newClient); err != nil {
			return c.Status(400).SendString("Failed to parse request")
		}
		id, err := admin.client.Save(ctx, &newClient)
		if err != nil {
			return c.Status(500).SendString("Failed to create user: " + err.Error())
		}

		return c.Status(201).JSON(fiber.Map{"id": id})

	})

}
