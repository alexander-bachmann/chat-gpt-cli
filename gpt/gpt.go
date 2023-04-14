package gpt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var chatEndpoint = "https://api.openai.com/v1/chat/completions"

type GPT struct {
	Messages []map[string]string
	client   *http.Client
	apiKey   string
	model    string
}

func New(apiKey string, model string) *GPT {
	return &GPT{
		client: &http.Client{},
		apiKey: apiKey,
		model:  model,
		Messages: []map[string]string{
			// the system role helps set the behavior of ChatGPT
			{"role": "system", "content": "You are a helpful assistant."},
		},
	}
}

func (g *GPT) Chat(prompt string) (string, error) {
	g.Messages = append(g.Messages, map[string]string{"role": "user", "content": prompt})
	err := g.sendChat()
	if err != nil {
		return "", err
	}

	// reverse iterate over messages to find index of last message from user
	lastUserMessageIdx := -1
	for i := len(g.Messages) - 1; i >= 0; i-- {
		if g.Messages[i]["role"] == "user" {
			lastUserMessageIdx = i
			break
		}
	}

	// build string containing all messages from assistant since last user message
	var sb strings.Builder
	for i := lastUserMessageIdx + 1; i < len(g.Messages); i++ {
		if g.Messages[i]["role"] == "assistant" {
			sb.WriteString(g.Messages[i]["content"])
		}
	}

	return sb.String(), nil
}

func (g *GPT) sendChat() error {
	data := map[string]any{
		"messages": g.Messages,
		"model":    g.model,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v\n", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", chatEndpoint, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("Error creating request: %v\n", err)
	}
	req.Header.Add("Authorization", "Bearer "+g.apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v\n", err)
	}
	defer resp.Body.Close()
	var respBody chatResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v\n", err)
	}
	for _, choice := range respBody.Choices {
		g.Messages = append(g.Messages, map[string]string{"role": "assistant", "content": choice.Message.Content})
	}
	return nil
}

type chatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   usage    `json:"usage"`
	Choices []choice `json:"choices"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type choice struct {
	Message message `json:"message"`
}

type message struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
}
