package chatbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Abraxas-365/commerce-chat/internal/database"
	"github.com/Abraxas-365/commerce-chat/pkg/assistant"
	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
	attributepg "github.com/Abraxas-365/commerce-chat/pkg/attribute/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/client"
	clientpg "github.com/Abraxas-365/commerce-chat/pkg/client/pgvector"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/Abraxas-365/commerce-chat/pkg/product"
	productpg "github.com/Abraxas-365/commerce-chat/pkg/product/pgvector"
)

type Chatbot struct {
	assistant   *assistant.Assistant
	productdb   product.Repository
	attributedb attribute.Repository
	clientdb    client.Repository
}

type Config struct {
	Db     *database.Connection
	Openia *openia.Openia
}

func New(c Config) *Chatbot {

	assistant := assistant.New(c.Openia)
	productdb := productpg.New(c.Db.Pool)
	attributedb := attributepg.New(c.Db.Pool)
	clientdb := clientpg.New(c.Db.Pool)
	return &Chatbot{
		assistant,
		productdb,
		attributedb,
		clientdb,
	}
}

func (c *Chatbot) ChatAllTheStore(messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()
	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content)
	if err != nil {
		return nil, err
	}
	//trart los mas parecidos de la db

	mostSimilarProducts, err := c.productdb.MostSimilarVectors(ctx, questionEmbedding, 15)
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
	c.assistant.AddSystemPrompt(sytemPrompt)

	chat, err := c.assistant.Help(messages)
	if err != nil {
		return nil, err
	}

	return chat, err
}

func (c *Chatbot) ChatWithProduct(sku string, messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()

	attributesCh := make(chan []string)
	errorCh := make(chan error)

	go func() {
		attributesArray, err := c.attributedb.GetBySKU(ctx, sku)
		if err != nil {
			errorCh <- err
			return
		}
		attributes := []string{}
		for _, attribute := range attributesArray {
			attributes = append(attributes, attribute.Information)
		}
		attributesCh <- attributes
	}()

	product, err := c.productdb.GetBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content + " " + product.Name)
	if err != nil {
		return nil, err
	}

	otherProducts, err := c.productdb.MostSimilarVectorsExeptProductBySku(ctx, questionEmbedding, 2, sku)
	var otherProductsIds []int
	for _, product := range otherProducts {
		otherProductsIds = append(otherProductsIds, product.Id)
	}
	otherAttributes, err := c.attributedb.GetByProducts(ctx, otherProductsIds)
	if err != nil {
		return nil, err
	}

	attributes, ok := <-attributesCh
	if !ok {
		err := <-errorCh
		return nil, err
	}

	productosArmados := []string{}
	for _, product := range otherProducts {
		fmt.Println(product.Name)
		productAttributes := otherAttributes[product.Id]
		attributesStr := ""
		for _, productAttribute := range productAttributes {
			attributesStr += productAttribute.Information + "\n"
		}
		productAndAttributes := fmt.Sprintf(`
Name: %s.
Attributes:
%s
		`, product.Name, attributesStr)

		productosArmados = append(productosArmados, productAndAttributes)
	}

	systemInfoPrompt := fmt.Sprintf("Other Products in stock that you can use to extend your answer: %s \n", strings.Join(productosArmados, "\n"))
	systemPrompt := fmt.Sprintf("Product in stock that is being consulted: %s \n product attributes %s. \n answer based on this product", product.Name, strings.Join(attributes, "\n"))
	c.assistant.AddSystemPrompt(systemPrompt)
	c.assistant.AddSystemPrompt(systemInfoPrompt)
	chat, err := c.assistant.Help(messages)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
