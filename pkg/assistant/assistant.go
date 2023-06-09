package assistant

import (
	"fmt"

	"github.com/Abraxas-365/commerce-chat/pkg/openia"
	"github.com/Abraxas-365/commerce-chat/pkg/openia/chat"
)

//structura tiene que tener el mensage del systema de cuando se inicia

type Assistant struct {
	systemConfPrompt string
	openia           *openia.Openia
	messages         chat.Messages
}

func New(systemConfPrompt string, openia *openia.Openia) *Assistant {
	return &Assistant{
		systemConfPrompt,
		openia,
		chat.Messages{chat.Message{
			Role:    "system",
			Content: systemConfPrompt,
		}},
	}
}

func (a *Assistant) Help(messages chat.Messages, systemInfo *string, addToSystemPrompt *string) (chat.Messages, error) {
	assitantConfigMessages := a.messages
	if addToSystemPrompt != nil {
		assitantConfigMessages[0].Content = fmt.Sprintf("%s \n %s", assitantConfigMessages[0].Content, *addToSystemPrompt)
	}

	messages = append(assitantConfigMessages, messages...)
	a.formatUserMessageAsQuestion(&messages)
	if systemInfo != nil {
		messages.AddSystemMessage(*systemInfo)
	}

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
