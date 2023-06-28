package chatbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/Abraxas-365/commerce-chat/pkg/assistant"
	"github.com/Abraxas-365/commerce-chat/pkg/attribute"
	"github.com/Abraxas-365/commerce-chat/pkg/client"
	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/Abraxas-365/commerce-chat/pkg/product"
)

type Chatbot struct {
	assistant *assistant.Assistant
	clientdb  client.Repository
	pservice  product.Service
}

type Config struct {
	Prepo  product.Repository
	Arepo  attribute.Repository
	Crepo  client.Repository
	Openia *openia.Openia
}

func New(c Config) *Chatbot {

	pservice := product.NewService(c.Prepo, c.Arepo)
	assistant := assistant.New(c.Openia)
	return &Chatbot{
		assistant,
		c.Crepo,
		pservice,
	}
}

func (c *Chatbot) ChatAllTheStore(messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()
	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content)
	if err != nil {
		return nil, err
	}

	mostSimilarProducts, err := c.pservice.GetByEmbedding(ctx, questionEmbedding, 10)
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

	principalProduct, err := c.pservice.GetBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	questionEmbedding, err := c.assistant.GetQuestionEmbedding(messages[len(messages)-1].Content + " " + principalProduct.Name)
	if err != nil {
		return nil, err
	}

	otherProducts, err := c.pservice.OtherSimilars(ctx, sku, questionEmbedding, 2)
	var otherProductsIds []int
	for _, product := range otherProducts {
		otherProductsIds = append(otherProductsIds, product.Id)
	}

	const productTemplate = `
Name: %s.
Attributes:
%s
`
	var productosArmados []string
	for _, product := range otherProducts {
		fmt.Println(product.Name)

		var attributesBuilder strings.Builder
		for _, productAttribute := range product.Attributes {
			attributesBuilder.WriteString(productAttribute)
			attributesBuilder.WriteString("\n")
		}

		productAndAttributes := fmt.Sprintf(productTemplate, product.Name, attributesBuilder.String())

		productosArmados = append(productosArmados, productAndAttributes)
	}

	systemInfoPrompt := fmt.Sprintf("Other Products in stock that you can use to extend your answer: %s \n", strings.Join(productosArmados, "\n"))
	systemPrompt := fmt.Sprintf("Product in stock that is being consulted: %s \n product attributes %s. \n answer based on this product", principalProduct.Name, strings.Join(principalProduct.Attributes, "\n"))
	c.assistant.AddSystemPrompt(systemPrompt)
	c.assistant.AddSystemPrompt(systemInfoPrompt)
	chat, err := c.assistant.Help(messages)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
