package llm

import (
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	api *openai.Client
}

func NewClient(apiKey string) *Client {
	return &Client{api: openai.NewClient(apiKey)}
}

func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := c.api.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "Eres un agente comercial de Kavak. Responde con base en la información proporcionada."},
			{Role: "user", Content: prompt},
		},
		MaxTokens: 300,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Client) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo",
		Messages:  messages,
		MaxTokens: 300,
	}

	resp, err := c.api.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
