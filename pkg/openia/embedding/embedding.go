package embedding

type Embedding struct {
	Data   []data `json:"data"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  usage  `json:"usage"`
}

type usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type data struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object"`
}

type requestBody struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}
