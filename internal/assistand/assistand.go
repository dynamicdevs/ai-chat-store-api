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

func (a *Assistand) HelpWithEveryThing(messages chat.Messages) (chat.Messages, error) {
	//get embedding of the question
	ctx := context.Background()
	question := messages[len(messages)-1].Content
	messages[len(messages)-1].Content = fmt.Sprintf(`Question:%s`, question)
	embedding, err := a.openia.GenerateEmbedding(question)
	if err != nil {
		return nil, err
	}
	attributedb := attributepg.New(a.db.Pool)
	mostSimilarProducts, err := attributedb.MostSimilarVectors(ctx, embedding, 3)
	if err != nil {
		return nil, err
	}

	listOfProducts := []int{}
	for _, product := range mostSimilarProducts {
		listOfProducts = append(listOfProducts, product.Id)
	}

	productIdToAttributes, err := attributedb.GetByProducts(ctx, listOfProducts)
	if err != nil {
		return nil, err
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
		fmt.Println(product.Name)

		productosArmados = append(productosArmados, productAndAttributes)

	}

	messages.AddSystemMessage(
		"Catalog of products you know are in stock, this are the only products you know are in stock:\n " +
			strings.Join(productosArmados, "\n"))
	messages.AddQuestion(question)
	response, err := a.openia.Chat(messages)
	if err != nil {
		return nil, err
	}

	messages.AddAssistant(response)

	return messages, nil
}
