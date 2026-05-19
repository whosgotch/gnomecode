package main

import (
	"os"
	"os/exec"
	"strings"
)

type RepositoryInfo struct {
	CurrentDir   string
	RootDir      string
	IsGitRepo    bool
	TrackedFiles []string
}

func inspectRepository() RepositoryInfo {
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = ""
	}

	rootOutput, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return RepositoryInfo{
			CurrentDir: currentDir,
			IsGitRepo:  false,
		}
	}

	rootDir := strings.TrimSpace(string(rootOutput))
	filesOutput, err := exec.Command("git", "-C", rootDir, "ls-files").Output()
	if err != nil {
		return RepositoryInfo{
			CurrentDir: currentDir,
			RootDir:    rootDir,
			IsGitRepo:  true,
		}
	}

	trackedFiles := []string{}
	for _, file := range strings.Split(string(filesOutput), "\n") {
		file = strings.TrimSpace(file)
		if file != "" {
			trackedFiles = append(trackedFiles, file)
		}
	}

	return RepositoryInfo{
		CurrentDir:   currentDir,
		RootDir:      rootDir,
		IsGitRepo:    true,
		TrackedFiles: trackedFiles,
	}
}
