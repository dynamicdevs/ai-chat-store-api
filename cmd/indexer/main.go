package main

import (
	"log"
	"os"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/internal/indexer"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
)

func main() {

	conn, err := database.NewConnection("chat-store-server.postgres.database.azure.com", 5432, "chatstoreadmin", "Ch4tSt0R34dm1n", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	indexer := indexer.New(conn, openia.New(os.Getenv("OPENAI_API_KEY")))
	if err := indexer.Index("output.csv"); err != nil {
		log.Fatal(err)
	}
}
