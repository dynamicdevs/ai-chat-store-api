package assistand

import (
	"context"
	"fmt"
	"strings"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
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

func (a *Assistand) HelpWithEveryThing(chat []chat.Message) (string, error) {
	//get embedding of the question
	ctx := context.Background()
	question := chat[len(chat)-1].Content
	chat[len(chat)-1].Content = fmt.Sprintf(`Question:%s`, question)
	embedding, err := a.openia.GenerateEmbedding(question)
	if err != nil {
		return "", err
	}
	attributedb := attributepg.New(a.db.Pool)
	mostSimilarProducts, err := attributedb.MostSimilarVectors(ctx, embedding, 5)
	if err != nil {
		return "", err
	}

	listOfProducts := []int{}
	for _, product := range mostSimilarProducts {
		listOfProducts = append(listOfProducts, product.Id)
	}

	productIdToAttributes, err := attributedb.GetByProducts(ctx, listOfProducts)
	if err != nil {
		return "", err
	}

	productosArmados := []string{}
	for _, product := range mostSimilarProducts {
		productAttributes := productIdToAttributes[product.Id]
		attributes := ""
		for _, productAttribute := range productAttributes {
			attributes += productAttribute.Information + "\n"
		}
		productAndAttributes := fmt.Sprintf(`
Name: %s.
Attributes:
%s
		`, product.Name, attributes)
		fmt.Println(productAndAttributes)

		productosArmados = append(productosArmados, productAndAttributes)

	}

	chat[len(chat)-1].Content = fmt.Sprintf(`Question:%s`, question)

	response, err := a.openia.Chat(chat, strings.Join(productosArmados, "\n"))
	if err != nil {
		return "", err
	}

	return response, nil
}
