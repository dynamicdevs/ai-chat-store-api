// Package chatbot provides a chatbot service integrated into an e-commerce platform. The chatbot is
// powered by AI natural language processing and supports user interactions for general queries
// as well as product-specific queries.
//
// The Chatbot struct is the core entity in this package which uses an AI assistant for processing
// natural language queries and interacts with the database to fetch product information.
//
// Example Usage:
//
//	// Initialize database connection and assistant
//	db := database.NewConnection(...)
//	aiAssistant := assistant.New(...)
//
//	// Create a Chatbot
//	chatbot := chatbot.New(db, aiAssistant)
//
//	// Interact with the Chatbot
//	messages := chat.Messages{
//	    {Role: "user", Content: "Can you tell me about product XYZ?"},
//	}
//	response, err := chatbot.ChatWithProduct("XYZ", messages)
//
// The package also provides a ControllerFactory for setting up HTTP routes and handling incoming
// HTTP requests. The factory can be used to integrate the chatbot service into a web server.
//
// The available HTTP routes are:
//
//	POST /api/messages
//	    - General purpose chatbot query. The request body should contain an array of chat messages,
//	      with the last message being the user's query. The chatbot responds with relevant information
//	      from the store.
//
//	POST /api/messages/product/:id
//	    - Product-specific query. Similar to the general purpose query, but focused on a specific product.
//	      The ':id' parameter in the URL should be the SKU or ID of the product the user is inquiring about.
//
// The Chatbot's responses are generated based on user queries and information fetched from the
// database about products and their attributes. It can answer general queries about the store
// and specific queries about products.
package chatbot
