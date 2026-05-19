package main

import (
	"fmt"
	"strings"
)

type ChatState struct {
	Model   string
	BaseURL string
	History []Message
}

func handleCommand(msg string, state *ChatState) bool {
	if msg == ":quit" || msg == "exit" {
		return false
	}

	if msg == ":help" {
		printHelp()
		return true
	}

	if msg == ":clear" {
		state.History = []Message{}
		fmt.Println("conversation cleared")
		return true
	}

	if msg == ":model" {
		fmt.Printf("model: %s\n", state.Model)
		return true
	}

	if msg == ":models" {
		models, err := listModels(state.BaseURL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return true
		}
		for _, modelName := range models {
			fmt.Println(modelName)
		}
		return true
	}

	if msg == ":base-url" {
		fmt.Printf("base URL: %s\n", state.BaseURL)
		return true
	}

	if strings.HasPrefix(msg, ":base-url ") {
		baseURL := strings.TrimSpace(strings.TrimPrefix(msg, ":base-url "))
		if baseURL == "" {
			fmt.Println("error: base URL is empty")
			return true
		}
		state.BaseURL = baseURL
		fmt.Printf("base URL: %s\n", state.BaseURL)
		return true
	}

	if msg == ":repo" {
		repo := inspectRepository()
		fmt.Printf("Current directory: %s\n", repo.CurrentDir)
		if repo.IsGitRepo {
			fmt.Println("Git repository: yes")
			fmt.Printf("Git root: %s\n", repo.RootDir)
			fmt.Printf("Files: %d\n", len(repo.TrackedFiles))
		} else {
			fmt.Println("Git repository: no")
		}
		return true
	}

	if msg == ":files" {
		fmt.Println(listFiles())
		return true
	}

	if strings.HasPrefix(msg, ":read ") {
		path := strings.TrimSpace(strings.TrimPrefix(msg, ":read "))
		fmt.Println(readFile(path))
		return true
	}

	if strings.HasPrefix(msg, ":search ") {
		query := strings.TrimSpace(strings.TrimPrefix(msg, ":search "))
		fmt.Println(search(query))
		return true
	}

	if strings.HasPrefix(msg, ":model ") {
		model := strings.TrimSpace(strings.TrimPrefix(msg, ":model "))
		if model == "" {
			fmt.Println("error: model name is empty")
			return true
		}
		state.Model = model
		fmt.Printf("model: %s\n", state.Model)
		return true
	}

	return false
}

func isExitCommand(msg string) bool {
	return msg == ":quit" || msg == "exit"
}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  :help          Show commands")
	fmt.Println("  :quit          Exit")
	fmt.Println("  exit           Exit")
	fmt.Println("  :clear         Clear conversation history")
	fmt.Println("  :model         Show current model")
	fmt.Println("  :model NAME    Set current model")
	fmt.Println("  :base-url      Show Ollama server URL")
	fmt.Println("  :base-url URL  Set Ollama server URL")
	fmt.Println("  :models        List Ollama models")
	fmt.Println("  :repo          Show repository info")
	fmt.Println("  :files         List tracked files")
	fmt.Println("  :read PATH     Read a tracked file")
	fmt.Println("  :search QUERY  Search tracked files")
}
