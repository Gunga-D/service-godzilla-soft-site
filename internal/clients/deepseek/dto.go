package deepseek

type CompletionsRequest struct {
	Model    string       `json:"model"`
	Messages []MessageDTO `json:"messages"`
	Stream   bool         `json:"stream"`
}

type MessageDTO struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionsResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []ChoiceDTO `json:"choices"`
}

type ChoiceDTO struct {
	Index   int        `json:"index"`
	Message MessageDTO `json:"message"`
}
