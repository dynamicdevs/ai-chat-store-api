package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/internal/indexer"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
)

func main() {

	conn, err := database.NewConnection("localhost", 5432, "myuser", "mypassword", "mydb")
	if err != nil {
		log.Fatal(err)
	}

	indexer := indexer.New(conn, openia.New(os.Getenv("OPENAI_API_KEY")))
	if err := indexer.Index("test.csv"); err != nil {
		log.Fatal(err)
	}
}
