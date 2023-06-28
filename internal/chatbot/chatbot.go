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

func (c *Chatbot) ChatRetrieveProductsBasedOnChat(messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()
	latestMessage := messages[len(messages)-1].Content

	questionEmbedding, err := c.assistant.GetQuestionEmbedding(latestMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to get question embedding: %w", err)
	}

	similarProducts, err := c.pservice.GetByEmbedding(ctx, questionEmbedding, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products by embedding: %w", err)
	}

	catalogPrompt := buildProductCatalogPrompt(similarProducts)
	c.assistant.AddSystemPrompt(catalogPrompt)

	return c.assistant.Help(messages)
}

func buildProductCatalogPrompt(products []product.ProdutDetailes) string {
	productDescriptions := make([]string, len(products))
	for i, product := range products {
		productDescriptions[i] = fmt.Sprintf("Product: %s.", product.Name)
	}

	return "Catalog of products you know are in stock, these are the only products you know are in stock:\n" +
		strings.Join(productDescriptions, "\n")
}

func (c *Chatbot) ChatWithRelevantProducts(sku string, messages chat.Messages) (chat.Messages, error) {
	ctx := context.Background()

	mainProduct, err := c.pservice.GetBySku(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve main product: %w", err)
	}

	questionWithProduct := messages[len(messages)-1].Content + " " + mainProduct.Name
	questionEmbedding, err := c.assistant.GetQuestionEmbedding(questionWithProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to get question embedding: %w", err)
	}

	similarProducts, err := c.pservice.OtherSimilars(ctx, sku, questionEmbedding, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve similar products: %w", err)
	}

	systemInfoPrompt := buildProductInfoPrompt(similarProducts)
	systemPrompt := fmt.Sprintf("Product in stock that is being consulted:\n%s\nAnswer based on this product.", buildProductInfo(*mainProduct))

	c.assistant.AddSystemPrompt(systemPrompt)
	c.assistant.AddSystemPrompt(systemInfoPrompt)

	return c.assistant.Help(messages)
}

func buildProductInfo(product product.ProdutDetailes) string {
	productInfoTemplate := `
Name: %s.
Attributes:
%s
`
	attributes := strings.Join(product.Attributes, "\n")

	return fmt.Sprintf(productInfoTemplate, product.Name, attributes)
}

func buildProductInfoPrompt(products []product.ProdutDetailes) string {
	var productInfoList []string
	for _, product := range products {
		productInfoList = append(productInfoList, buildProductInfo(product))
	}

	return fmt.Sprintf("Other Products in stock that you can use to extend your answer:\n%s", strings.Join(productInfoList, "\n"))
}
