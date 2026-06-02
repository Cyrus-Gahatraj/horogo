package internal

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

const MODEL = "qwen3:4b"

func Ask(prompt string, chartData string) {
	ctx := context.Background()
	client, err := api.ClientFromEnvironment()

	if err != nil {
		panic(err)
	}

	fullPrompt := fmt.Sprintf("Birth Chart Data:\n%s\n\nUser Question: %s", chartData, prompt)

	req := &api.GenerateRequest{
		Model:  MODEL,
		Prompt: fullPrompt,
		Stream: nil,
	}

	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println()
}
