package openia

import (
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/embedding"
)

type Openia struct {
	apiKey string
}

func New(apiKey string) *Openia {
	return &Openia{
		apiKey: apiKey,
	}
}

func (o *Openia) GenerateEmbedding(text string) ([]float32, error) {
	embedding, err := embedding.GenerateEmbedding(o.apiKey, []string{text})
	if err != nil {
		return nil, err
	}

	return embedding.Data[0].Embedding, nil

}

func (o *Openia) Chat(currentChatState []chat.Message) (string, error) {
	chatHistory := []chat.Message{
		{
			Role: "system",
			Content: `You are an ecommerce asystenat of ABCDIN that is going to help the customer aswering
			their question about products that are in the stock of the store, maybe comparing some products that are in on stock of the store, give charactristic, etc.
			If someone ask something not relaited to retail or the store, aswer with sorry i cant help you.
			Dont menssion products that are not in the catalog, only use the products you know are in stock.
			You only can answer with the products you know are in stock,If you dont know the products, Say you dont have it.
			`,
		},
	}
	chatHistory = append(chatHistory, currentChatState...)

	chat, err := chat.Chat(o.apiKey, chatHistory)
	if err != nil {
		return "", err
	}

	return chat.Choices[0].Message.Content, err

}
