package main

import (
	"log"

	"github.com/Abraxas-365/commerce-chat/internal/admin"
	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

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

	conf := admin.Config{
		Db: conn,
	}
	admin.ControllerFactory(app, conf)

	log.Fatal(app.Listen(":80"))

}
