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

func (o *Openia) Chat(chatHistory []chat.Message) (string, error) {

	chat, err := chat.Chat(o.apiKey, chatHistory)
	if err != nil {
		return "", err
	}

	return chat.Choices[0].Message.Content, err

}
