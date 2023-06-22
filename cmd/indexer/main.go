package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/internal/indexer"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
)

func main() {

	conn, err := database.NewConnection(os.Getenv("DB_URI"), 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))
	if err != nil {
		log.Fatal(err)
	}

	indexer := indexer.New(conn, openia.New(os.Getenv("OPENAI_API_KEY")))
	if err := indexer.Index("output.csv"); err != nil {
		log.Fatal(err)
	}
}
