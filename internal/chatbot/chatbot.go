package chatbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/assistant"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	productpg "github.com/Abraxas-365/commerce-chat/pkg/product/pgvector"
)

type Chatbot struct {
	db        *database.Connection
	assistant *assistant.Assistant
}

func New(db *database.Connection, assistant *assistant.Assistant) *Chatbot {
	return &Chatbot{
		db,
		assistant,
	}
}

func (c *Chatbot) ChatAllTheStore(messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()
	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content)
	if err != nil {
		return nil, err
	}
	//trart los mas parecidos de la db

	productdb := productpg.New(c.db.Pool)
	mostSimilarProducts, err := productdb.MostSimilarVectors(ctx, questionEmbedding, 15)
	if err != nil {
		return nil, err
	}
	productosArmados := []string{}
	for _, product := range mostSimilarProducts {
		productAndAttributes := fmt.Sprintf(`Product: %s.`, product.Name)
		productosArmados = append(productosArmados, productAndAttributes)
		fmt.Println(product.Name)
	}

	sytemPrompt := "Catalog of products you know are in stock, this are the only products you know are in stock:\n " +
		strings.Join(productosArmados, "\n")

	chat, err := c.assistant.Help(messages, &sytemPrompt, nil)
	if err != nil {
		return nil, err
	}

	return chat, err
}

func (c *Chatbot) ChatWithProduct(sku string, messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()
	productdb := productpg.New(c.db.Pool)
	attributedb := attributepg.New(c.db.Pool)
	product, err := productdb.GetBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content + " " + product.Name)
	if err != nil {
		return nil, err
	}
	otherProducts, err := attributedb.MostSimilarVectorsExeptProductBySku(ctx, questionEmbedding, 2, sku)
	var otherProductsIds []int
	for _, product := range otherProducts {
		otherProductsIds = append(otherProductsIds, product.Id)
	}
	otherAttributes, err := attributedb.GetByProducts(ctx, otherProductsIds)
	if err != nil {
		return nil, err
	}
	productosArmados := []string{}
	for _, product := range otherProducts {
		fmt.Println(product.Name)
		productAttributes := otherAttributes[product.Id]
		attributes := ""
		for _, productAttribute := range productAttributes {
			attributes += productAttribute.Information + "\n"
		}
		productAndAttributes := fmt.Sprintf(`
Name: %s.
Attributes:
%s
		`, product.Name, attributes)

		productosArmados = append(productosArmados, productAndAttributes)

	}

	attributesArray, err := attributedb.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}
	attributes := []string{}
	for _, attribute := range attributesArray {
		attributes = append(attributes, attribute.Information)
	}

	systemPrompt := fmt.Sprintf("Product in stock that is beeing consulted: %s \n product attribures %s. \n aswer in base of this product", product.Name, strings.Join(attributes, "\n"))
	systemInfoPrompt := fmt.Sprintf("Other Products in stock that you can use to extend your answer: %s \n", strings.Join(productosArmados, "\n"))
	chat, err := c.assistant.Help(messages, &systemInfoPrompt, &systemPrompt)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
