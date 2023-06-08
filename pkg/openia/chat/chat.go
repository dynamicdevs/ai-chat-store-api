package chat

import "fmt"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Messages []Message

func (m *Messages) AddSystemMessage(message string) {
	index := len(*m) - 1
	*m = append((*m)[:index], append([]Message{
		{
			Role:    "system",
			Content: message,
		},
	}, (*m)[index:]...)...)
}

func (m *Messages) AddUserMessage(message string) {
	*m = append(*m, Message{
		Role:    "user",
		Content: message,
	})
}

func (m *Messages) AddQuestion(content string) {
	index := len(*m) - 1
	(*m)[index].Content = fmt.Sprintf("Question: %s", content)
}

func (m *Messages) AddAssistant(content string) {
	*m = append(*m, Message{
		Role:    "assistant",
		Content: content,
	})
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}
