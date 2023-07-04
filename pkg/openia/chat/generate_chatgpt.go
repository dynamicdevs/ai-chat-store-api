package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Chat(openiaKey string, messages []Message) (*ChatCompletion, error) {
	url := "https://api.openai.com/v1/chat/completions"

	data := map[string]interface{}{
		"model":       "gpt-3.5-turbo-16k",
		"messages":    messages,
		"temperature": 0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openiaKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result ChatCompletion
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
