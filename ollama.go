package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type chatResponse struct {
	Message Message `json:"message"`
}

type chatStreamResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type modelsResponse struct {
	Models []modelInfo `json:"models"`
}

type modelInfo struct {
	Name string `json:"name"`
}

func chat(model string, messages []Message, baseURL string) (string, error) {
	payload, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/api/chat", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("could not connect to Ollama at %s", baseURL)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Ollama returned an error: %s", string(body))
	}

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("Ollama returned an unexpected response")
	}

	if parsed.Message.Content == "" {
		return "", fmt.Errorf("Ollama returned an empty response")
	}

	return parsed.Message.Content, nil
}

func chatStream(model string, messages []Message, baseURL string, onToken func(string)) (string, error) {
	payload, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/api/chat", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("could not connect to Ollama at %s", baseURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Ollama returned an error: %s", string(body))
	}

	var fullResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var chunk chatStreamResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			return "", fmt.Errorf("Ollama returned an unexpected response")
		}

		token := chunk.Message.Content
		if token != "" {
			onToken(token)
			fullResponse += token
		}

		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fullResponse, nil
}

func listModels(baseURL string) ([]string, error) {
	resp, err := http.Get(baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama at %s", baseURL)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Ollama returned an error: %s", string(body))
	}

	var parsed modelsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("Ollama returned an unexpected response")
	}

	names := []string{}
	for _, model := range parsed.Models {
		if model.Name != "" {
			names = append(names, model.Name)
		}
	}

	return names, nil
}
