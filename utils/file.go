package utils

import (
	"bufio"
	"os"
	"strings"
)

// readFileLinesToArray reads a file line by line and returns an array of strings
func readFileLinesToArray(filePath string) ([]string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// parseProjectAndRepoName parses a Git repository URL and returns the project and repository name
func parseProjectAndRepoName(repoURL string) (string, string) {
	// Split the input string on slashes
	parts := strings.Split(repoURL, "/")

	// Check if there are any parts after splitting
	if len(parts) > 0 {

		// The project name is the second to last part of the URL
		projectName := parts[len(parts)-2]

		// The repo name is the last part of the URL
		repoName := parts[len(parts)-1]

		return projectName, repoName
	}

	// If there are no parts, return an empty string
	return "", ""
}
