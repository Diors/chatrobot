package api

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/revel/revel"
	"github.com/sashabaranov/go-openai"
)

func stream(message string) string {
	CHATGPT_TOKEN, _ := revel.Config.String("openai.token")
	c := openai.NewClient(CHATGPT_TOKEN)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		panic(err)
	}
	defer stream.Close()

	var replymsg string

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return replymsg
		}

		if err != nil {
			revel.AppLog.Error("Stream error: %v\n", err)
			panic(err)
		}

		return response.Choices[0].Delta.Content
	}
}

func CompletionRequest(message string) string {
	client := openai.NewClient(CHATGPT_TOKEN)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		revel.AppLog.Error("ChatCompletion error: ", err)
		return err.Error()
	}

	return resp.Choices[0].Message.Content
}
