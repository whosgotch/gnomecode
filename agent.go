package main

import (
	"fmt"
	"strings"
)

type ParsedAgentResponse struct {
	Kind  string
	Name  string
	Input string
	Text  string
}

func buildAgentMessages(task string) []Message {
	systemPrompt := `You are GnomeCode, a read-only coding agent.

You can answer normal questions directly.

For questions about repository code, use tools before answering. Do not guess.

Available tools:
- list_files
- read_file
- search

Use exactly one of these formats:

ACTION: list_files

ACTION: read_file
INPUT: path/to/file.go

ACTION: search
INPUT: query

FINAL:
answer

Rules:
- Output exactly one command per response.
- Do not output ACTION and FINAL in the same response.
- If you output ACTION, stop immediately after the optional INPUT line.
- If you output FINAL, do not include placeholders.
- Never say that a tool result "will be displayed"; use ACTION to request the tool.
- If the user names a specific file, read that file first.
- If the user asks to explain a specific file, read that file and then answer.
- After read_file returns the requested file content, answer the user directly.
- Do not call list_files after successfully reading a user-specified file.
- Do not search for unrelated topics unless the user asks.
- If the user asks to show file contents or code, use read_file and return exact code from tool results.
- Do not invent code.
- After receiving TOOL RESULT, either call another tool or return FINAL.
- You are read-only. Do not write files.`

	return []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: task},
	}
}

func parseAgentResponse(response string) ParsedAgentResponse {
	lines := strings.Split(response, "\n")
	action := ""
	input := ""
	final := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "FINAL:") {
			final = strings.TrimSpace(strings.TrimPrefix(line, "FINAL:"))
		}
		if strings.HasPrefix(line, "ACTION:") {
			action = strings.TrimSpace(strings.TrimPrefix(line, "ACTION:"))
		}
		if strings.HasPrefix(line, "INPUT:") {
			input = strings.TrimSpace(strings.TrimPrefix(line, "INPUT:"))
		}
	}

	if action != "" {
		return ParsedAgentResponse{Kind: "action", Name: action, Input: input}
	}

	if final != "" {
		return ParsedAgentResponse{Kind: "final", Text: final}
	}

	if strings.Contains(response, "FINAL:") {
		return ParsedAgentResponse{Kind: "final", Text: strings.TrimSpace(strings.SplitN(response, "FINAL:", 2)[1])}
	}

	return ParsedAgentResponse{Kind: "final", Text: response}
}

func runTool(name string, input string) string {
	switch name {
	case "list_files":
		return listFiles()
	case "read_file":
		return readFile(input)
	case "search":
		return search(input)
	default:
		return fmt.Sprintf("Error: unknown tool %s", name)
	}
}

func runAgent(task string, model string, baseURL string, maxSteps int) string {
	messages := buildAgentMessages(task)

	for step := 0; step < maxSteps; step++ {
		response, err := chat(model, messages, baseURL)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}

		parsed := parseAgentResponse(response)
		if parsed.Kind == "final" {
			return parsed.Text
		}

		if parsed.Kind == "action" {
			if parsed.Input != "" {
				fmt.Printf("[tool] %s %s\n", parsed.Name, parsed.Input)
			} else {
				fmt.Printf("[tool] %s\n", parsed.Name)
			}

			result := runTool(parsed.Name, parsed.Input)
			toolLabel := parsed.Name
			if parsed.Input != "" {
				toolLabel += " " + parsed.Input
			}
			messages = append(messages, Message{Role: "assistant", Content: response})
			messages = append(messages, Message{
				Role:    "user",
				Content: fmt.Sprintf("TOOL RESULT FROM %s:\n%s", toolLabel, result),
			})
		}
	}

	return "Error: agent reached max steps"
}
