// Package assistant provides an implementation for an AI chat assistant that integrates
// with the OpenAI GPT models for natural language processing. The Assistant struct and
// related functions are designed to facilitate interaction with the OpenAI models, by
// formatting and processing the messages exchanged between the user and the AI assistant.
//
// The Assistant struct is the primary entity in this package which holds the system
// configuration prompt, an instance of OpenAI client and the chat messages.
//
// Example Usage:
//
//	// Create an OpenIA client
//	openiaClient := openia.New(...)
//
//	// Create an Assistant with system configuration prompt
//	assistant := assistant.New("This is a chatbot that helps with...", openiaClient)
//
//	// Message to be processed
//	messages := chat.Messages{
//	    {Role: "user", Content: "What's the weather like today?"},
//	}
//
//	// Get assistant's response
//	response, err := assistant.Help(messages, nil, nil)
//
// The Assistant’s responses are generated based on the provided system configuration prompt,
// and any additional messages passed. The package also provides functionality for obtaining
// embeddings for the user’s questions, which can be used in similarity searches and other
// NLP tasks.
//
// The package is designed to be integrated into an e-commerce chatbot application.
package assistant
