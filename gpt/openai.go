package gpt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type AdmissionDecision = struct {
	Admitted bool   `json:"admitted"`
	Reason   string `json:"reason"`
}

const systemPrompt = `You are a Kubernetes admission controller
and get to randomly decide whether a given resource is admitted
to a cluster. Please decide and give a snarky reason.
Return a JSON object with a boolean field "admitted" and a string
field "reason" for the reason.`

func Decide(ctx context.Context, client *openai.Client, object string) (AdmissionDecision, error) {
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo0125,
		MaxTokens: 512,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: object,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return AdmissionDecision{}, fmt.Errorf("failed to create chat completion: %w", err)
	}

	var decision AdmissionDecision
	json.Unmarshal([]byte(resp.Choices[0].Message.Content), &decision)
	if err != nil {
		return AdmissionDecision{}, fmt.Errorf("failed to parse chat response: %w", err)
	}

	return decision, nil
}
