package openia

import (
	"fmt"
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
	fmt.Println("Attribute", text)
	if err != nil {
		return nil, err
	}

	return embedding.Data[0].Embedding, nil

}

func (o *Openia) Chat(currentChatState []chat.Message) (string, error) {
	chatHistory := []chat.Message{
		{
			Role: "system",
			Content: `You are an ecommerce asystenat of ABCDIN that is goig to help the customer aswering
			their question about products, maybe comparing some products, give charactristic, etc.
			If someone ask something not relaited to retail or the store, aswer with sorry i cant help you`,
		},
	}
	chatHistory = append(chatHistory, currentChatState...)
	fmt.Println("Chat History:", chatHistory)
	chat, err := chat.Chat(o.apiKey, chatHistory)
	if err != nil {
		return "", err
	}

	return chat.Choices[0].Message.Content, err

}
