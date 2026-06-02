package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultServerURL = "http://localhost:8000"

func Ask(prompt string, chartData string) error {
	serverURL := os.Getenv("AI_SERVER_URL")
	if serverURL == "" {
		serverURL = defaultServerURL
	}

	body := map[string]string{
		"prompt":     prompt,
		"chart_data": chartData,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(serverURL+"/ask", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Read chunk by chunk and print directly to stdout
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[:n]))
			os.Stdout.Sync()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read stream: %w", err)
		}
	}

	return nil
}

// AskSilent sends the prompt to the AI server and returns the response as a string
func AskSilent(prompt string, chartData string) (string, error) {
	serverURL := os.Getenv("AI_SERVER_URL")
	if serverURL == "" {
		serverURL = defaultServerURL
	}

	body := map[string]string{
		"prompt":     prompt,
		"chart_data": chartData,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(serverURL+"/ask", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	return buf.String(), nil
}
