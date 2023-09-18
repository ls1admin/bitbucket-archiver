package main

//TODOs:
// - Zip the repos after cloning

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	log "github.com/sirupsen/logrus"
)

var clone_wg sync.WaitGroup

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

func cloneGitRepo(repoURL, destPath, username, password string) error {
	// Clone the repository from the given URL
	defer clone_wg.Done()

	projectName, repoName := parseProjectAndRepoName(repoURL)

	destPath = destPath + "/" + projectName + "/" + repoName

	log.Debugf("Cloning repo: %s to: %s", repoName, destPath)

	_, err := git.PlainClone(destPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: nil, // Avoid sending progress to stdout and the server
		Auth: &http.BasicAuth{
			Username: username, // Your Git username
			Password: password, // Your Git password or personal access token
		},
	})

	if err != nil {
		log.WithError(err).Errorf("Error cloning repo: %s to: %s ", repoName, destPath)
	}

	return nil
}

// cloneListOfRepos clones a list of Git repositories to the given destination path in parrallel
// This funciton limits the number of parallel clones to 30
func cloneListOfRepos(repos []Repo, destPath, username, password string) error {
	log.Debugf("Cloning %d repos", len(repos))
	// safety check
	if len(repos) > 30 {
		log.Error("too many repos to clone at once")
		return errors.New("too many repos to clone at once")
	}

	for _, repo := range repos {
		clone_wg.Add(1)
		go cloneGitRepo(*repo.GetHTTPRepoUrl(), destPath, username, password)
	}
	clone_wg.Wait()
	return nil
}

// function that creates chunks of length n from a slice of strings
func chunks[S any](s []S, n int) [][]S {
	var chunks [][]S
	for n < len(s) {
		s, chunks = s[n:], append(chunks, s[0:n:n])
	}
	return append(chunks, s)
}

func main() {
	if is_debug := os.Getenv("DEBUG"); is_debug == "true" {
		log.SetLevel(log.DebugLevel)
		log.Warn("DEBUG MODE ENABLED")
	}

	LoadConfig() // Load the config file into a global variable
	repos, err := GetArchivedRepositories()
	if err != nil {
		log.WithError(err).Panic("Error reading the file to Array")
	}
	// Log the whole repos slice with new lines between each repo
	log.Info("Number of repos: ", len(repos))

	destPath := "./repos"

	// Split the list of repos into chunks to limit parallelism
	repoChunks := chunks(repos, 30)
	for _, chunk := range repoChunks {
		cloneListOfRepos(chunk, destPath, Cfg.GitUsername, Cfg.GitPassword)
	}
	log.Info("Starting to zip repos")
	// ZipAllProjects("repos", Cfg.OutputDir)
}
