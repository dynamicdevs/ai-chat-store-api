package embedding

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GenerateEmbedding(apiKey string, text []string) (*Embedding, error) {
	// Create the request body
	reqBody := &requestBody{
		Input: text,
		Model: "text-embedding-ada-002",
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// Define the request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, _ := ioutil.ReadAll(resp.Body)
	var result Embedding
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
