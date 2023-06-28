package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/chatbot"
	"github.com/Abraxas-365/commerce-chat/internal/database"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	clientpg "github.com/Abraxas-365/commerce-chat/pkg/client/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	productpg "github.com/Abraxas-365/commerce-chat/pkg/product/pgvector"
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

	conf := chatbot.Config{
		Prepo:  productpg.New(conn.Pool),
		Arepo:  attributepg.New(conn.Pool),
		Crepo:  clientpg.New(conn.Pool),
		Openia: openia.New(os.Getenv("OPENAI_API_KEY")),
	}
	chatbot.ControllerFactory(app, conf)
	// Serve the Swagger UI
	app.Static("/docs", "./dist")
	// Serve the Swagger YAML file
	app.Use("/swagger.yml", func(c *fiber.Ctx) error {
		return c.SendFile("./swagger.yml")
	})

	log.Fatal(app.Listen(":80"))

}
