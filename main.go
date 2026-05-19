package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	state := ChatState{
		Model:   "llama3.1",
		BaseURL: "http://localhost:11434",
		History: []Message{},
	}

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

		if isExitCommand(msg) {
			break
		}

		if handleCommand(msg, &state) {
			continue
		}

		result := runAgent(msg, state.History, state.Model, state.BaseURL, 5)
		fmt.Printf("agent> %s\n", result.Answer)

		state.History = append(state.History, result.Messages...)
	}
}
