package openapi

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	openaiClient *openai.Client
	Request      openai.ChatCompletionRequest
}

func NewClient(token string) *Client {
	chatReq := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Тебя зовут Джейн, девушка 28 лет. отвечай пожалуйста всегда по русски, если только пользователь не попросит обратного.",
			},
		},
	}

	client := openai.NewClient(token)
	return &Client{openaiClient: client, Request: chatReq}
}

func (c *Client) ClearChatContext() {
	c.Request.Messages = c.Request.Messages[:1]
}

func (c *Client) Chat(ctx context.Context, userInput string) (string, error) {
	c.Request.Messages = append(c.Request.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	resp, err := c.openaiClient.CreateChatCompletion(ctx, c.Request)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) CreateText(ctx context.Context, fpath string) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: fpath,
	}

	resp, err := c.openaiClient.CreateTranscription(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}

func (c *Client) CreateVoice(ctx context.Context, mess string) (openai.RawResponse, error) {
	req := openai.CreateSpeechRequest{
		Model: openai.TTSModel1,
		Input: mess,
		Voice: openai.VoiceNova,
	}

	return c.openaiClient.CreateSpeech(ctx, req)
}
