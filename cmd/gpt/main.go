package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/assistand"
	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
)

func main() {

	//prender server
	// app := fiber.New()
	// app.Use(cors.New())
	// app.Use(logger.New())
	//resivir preguntas
	//enviar respuesta al cliente

	conn, err := database.NewConnection("localhost", 5432, "myuser", "mypassword", "mydb")
	if err != nil {
		log.Fatal(err)
	}

	assistand := assistand.New(conn, openia.New(os.Getenv("OPENAI_API_KEY")))

	resp, err := assistand.HelpWithEveryThing("Soy un estudiante, me gustaria comprar una computadora, para estudiar y poder jugar")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
