package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	model := "llama3.1"
	baseURL := "http://localhost:11434"
	messages := []Message{}

	for {
		fmt.Print("you> ")

		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		if msg == ":quit" || msg == "exit" {
			break
		}

		if msg == ":clear" {
			messages = []Message{}
			fmt.Println("conversation cleared")
			continue
		}

		if msg == ":model" {
			fmt.Printf("model: %s\n", model)
			continue
		}

		if msg == ":models" {
			models, err := listModels(baseURL)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			for _, modelName := range models {
				fmt.Println(modelName)
			}
			continue
		}

		if strings.HasPrefix(msg, ":model ") {
			model = strings.TrimSpace(strings.TrimPrefix(msg, ":model "))
			if model == "" {
				fmt.Println("error: model name is empty")
				continue
			}
			fmt.Printf("model: %s\n", model)
			continue
		}

		messages = append(messages, Message{
			Role:    "user",
			Content: msg,
		})

		fmt.Print("agent> ")
		response, err := chatStream(model, messages, baseURL, func(token string) {
			fmt.Print(token)
		})
		if err != nil {
			fmt.Println()
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Println()

		messages = append(messages, Message{
			Role:    "assistant",
			Content: response,
		})
	}
}
