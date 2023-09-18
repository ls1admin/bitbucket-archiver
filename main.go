package main

//TODOs:
// - Zip the repos after cloning

import (
	"bitbucket_archiver/bitbucket"
	"bitbucket_archiver/utils"
	"path"

	"errors"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	log "github.com/sirupsen/logrus"
)

var clone_wg sync.WaitGroup

func cloneGitRepo(repo bitbucket.Repo, username, password string) error {
	// Clone the repository from the given URL
	defer clone_wg.Done()

	destPath := path.Join(utils.Cfg.CloneDir, repo.Project.Name, repo.Name)

	log.Debugf("Cloning repo: %s to: %s", repo.Name, destPath)

	_, err := git.PlainClone(destPath, false, &git.CloneOptions{
		URL:      *repo.GetHTTPRepoUrl(),
		Progress: nil, // Avoid sending progress to stdout and the server
		Auth: &http.BasicAuth{
			Username: username, // Your Git username
			Password: password, // Your Git password or personal access token
		},
	})

	if err != nil {
		log.WithError(err).Errorf("Error cloning repo: %s to: %s ", repo.Name, destPath)
	}

	return nil
}

// cloneListOfRepos clones a list of Git repositories to the given destination path in parrallel
// This funciton limits the number of parallel clones to 30
func cloneListOfRepos(repos []bitbucket.Repo, username, password string) error {
	log.Debugf("Cloning %d repos", len(repos))
	// safety check
	if len(repos) > 20 {
		log.Error("too many repos to clone at once")
		return errors.New("too many repos to clone at once")
	}

	for _, repo := range repos {
		clone_wg.Add(1)
		go cloneGitRepo(repo, username, password)
	}
	clone_wg.Wait()
	return nil
}

func main() {
	if is_debug := os.Getenv("DEBUG"); is_debug == "true" {
		log.SetLevel(log.DebugLevel)
		log.Warn("DEBUG MODE ENABLED")
	}

	utils.LoadConfig() // Load the config file into a global variable

	var repos []bitbucket.Repo

	// Distinguish between API based archiving and file based archiving
	if utils.Cfg.ProjectFile != "" {
		// Project file is defined -> use file based archiving
		log.Info("Starting file based archiving")

		projects, err := utils.ReadFileLinesToArray(utils.Cfg.ProjectFile)
		if err != nil {
			log.WithError(err).Fatal("Error reading the file to Array")
		}

		// Iterate over all projects and gather all repos
		for _, project := range projects {
			log.Debug("Append repos from Project: ", project)

			projectRepos, err := bitbucket.GetRepositoriesForProject(project)

			log.Debugf("Found %d repos in %s: ", len(projectRepos), project)
			if err != nil {
				log.WithError(err).Panic("Error getting repos for project")
			}

			repos = append(repos, projectRepos...)
		}

	} else {
		// Get repositories marked as archived from Bitbucket
		archivedRepos, err := bitbucket.GetArchivedRepositories()
		if err != nil {
			log.WithError(err).Panic("Error reading the file to Array")
		}
		repos = append(repos, archivedRepos...)
	}

	// Log the whole repos slice with new lines between each repo
	log.Info("Number of repos: ", len(repos))

	// Split the list of repos into chunks to limit parallelism
	repoChunks := utils.Chunks(repos, 20)
	for _, chunk := range repoChunks {
		cloneListOfRepos(chunk, utils.Cfg.GitUsername, utils.Cfg.GitPassword)
	}
	log.Info("Starting to zip repos")
	// ZipAllProjects("repos", Cfg.OutputDir)
}
