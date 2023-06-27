package assistant

import (
	"fmt"

	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
)

//structura tiene que tener el mensage del systema de cuando se inicia

type Assistant struct {
	openia   *openia.Openia
	messages chat.Messages
}

func New(openia *openia.Openia) *Assistant {
	return &Assistant{
		openia,
		chat.Messages{},
	}
}

func (a *Assistant) AddSystemPrompt(prompt string) {
	a.messages = append(a.messages, chat.Message{
		Role:    "system",
		Content: prompt,
	})
}

func (a *Assistant) Help(messagesHistory chat.Messages) (chat.Messages, error) {

	messages := append(a.messages, messagesHistory...)
	a.formatUserMessageAsQuestion(&messages)

	result, err := a.openia.Chat(messages)
	if err != nil {
		return nil, err
	}

	messages.AddAssistant(result)

	return messages[1:], nil
}

func (a *Assistant) GetQuestionEmbedding(question string) ([]float32, error) {
	return a.openia.GenerateEmbedding(question)
}

func (a *Assistant) formatUserMessageAsQuestion(messages *chat.Messages) {
	userMessage := len(*messages) - 1
	(*messages)[userMessage].Content = fmt.Sprintf("Question: %s", (*messages)[userMessage].Content)
}
