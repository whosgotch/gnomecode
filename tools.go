package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFiles() string {
	repo := inspectRepository()
	if !repo.IsGitRepo {
		return "Error: not inside a git repository"
	}

	return strings.Join(repo.TrackedFiles, "\n")
}

func readFile(path string) string {
	repo := inspectRepository()
	if !repo.IsGitRepo {
		return "Error: not inside a git repository"
	}

	if !isTrackedFile(path, repo.TrackedFiles) {
		return "Error: file is not tracked in this repository"
	}

	fullPath := filepath.Join(repo.RootDir, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Sprintf("Error: could not read file: %v", err)
	}

	if len(content) > 4000 {
		content = content[:4000]
	}

	return string(content)
}

func search(query string) string {
	repo := inspectRepository()
	if !repo.IsGitRepo {
		return "Error: not inside a git repository"
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return "Error: search query is empty"
	}

	queryLower := strings.ToLower(query)
	results := []string{}

	for _, path := range repo.TrackedFiles {
		fullPath := filepath.Join(repo.RootDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		for index, line := range lines {
			if strings.Contains(strings.ToLower(line), queryLower) {
				results = append(results, fmt.Sprintf("%s:%d:%s", path, index+1, strings.TrimSpace(line)))
				if len(results) >= 20 {
					return strings.Join(results, "\n")
				}
			}
		}
	}

	if len(results) == 0 {
		return "No matches found."
	}

	return strings.Join(results, "\n")
}

func isTrackedFile(path string, trackedFiles []string) bool {
	for _, trackedFile := range trackedFiles {
		if path == trackedFile {
			return true
		}
	}
	return false
}
