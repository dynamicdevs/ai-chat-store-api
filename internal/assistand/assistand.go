package assistand

import (
	"context"
	"fmt"
	"strings"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
)

type Assistand struct {
	db     *database.Connection
	openia *openia.Openia
}

func New(db *database.Connection, openia *openia.Openia) *Assistand {
	return &Assistand{
		db,
		openia,
	}
}

func (a *Assistand) HelpWithEveryThing(question string) (string, error) {
	//get embedding of the question
	ctx := context.Background()
	embedding, err := a.openia.GenerateEmbedding(question)
	if err != nil {
		return "", err
	}
	attributedb := attributepg.New(a.db.Pool)

	productToAttribute := make(map[int][]string)

	mostSimilarAttributes, mostSimilarProducts, err := attributedb.MostSimilarVectors(ctx, embedding, 5)
	if err != nil {
		return "", err
	}

	for _, attribute := range mostSimilarAttributes {
		productToAttribute[attribute.Product] = append(productToAttribute[attribute.Product], attribute.Information)
	}

	productosArmados := []string{}
	for _, product := range mostSimilarProducts {
		attributes := strings.Join(productToAttribute[product.Id], "\n")
		productAndAttributes := fmt.Sprintf(`Name: %s.
			Attributes:
			%s
		`, product.Name, attributes)

		productosArmados = append(productosArmados, productAndAttributes)

	}

	prompt := fmt.Sprintf(`Question:%s`, question)

	response, err := a.openia.Chat(prompt, strings.Join(productosArmados, "\n"))
	if err != nil {
		return "", err
	}

	return response, nil
}
