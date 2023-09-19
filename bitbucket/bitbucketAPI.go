package bitbucket

import (
	"bitbucket_archiver/utils"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var client = http.Client{}

// Outer wrapper for the JSON payload
type RepositoryPayload struct {
	Values        []Repo `json:"values"`
	NextPageStart int    `json:"nextPageStart"`
	IsLastPage    bool   `json:"isLastPage"`
}

func parseRepositoryPayloadJSON(jsonBytes []byte) (r RepositoryPayload) {
	err := json.Unmarshal(jsonBytes, &r)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshalling JSON")
	}
	return
}

func GetArchivedRepositories() ([]Repo, error) {
	return getPaginatedAPI(fmt.Sprintf("%s/rest/api/latest/repos?archived=ARCHIVED", utils.Cfg.BitbucketUrl))
}

func GetRepositoriesForProject(projectKey string) ([]Repo, error) {
	return getPaginatedAPI(fmt.Sprintf("%s/rest/api/latest/projects/%s/repos", utils.Cfg.BitbucketUrl, projectKey))
}

func getPaginatedAPI(url string) ([]Repo, error) {
	isLastPage := false
	currentPage := 0
	var repos []Repo

	for !isLastPage {
		payload, err := getRepoPayload(url, currentPage)
		if err != nil {
			log.WithError(err).Fatal("Error getting repo payload")
			return nil, err
		}

		repos = append(repos, payload.Values...)
		isLastPage = payload.IsLastPage
		currentPage = payload.NextPageStart
	}

	return repos, nil
}

func basicAuth() string {
	auth := utils.Cfg.BitbucketUsername + ":" + utils.Cfg.BitbucketPassword
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func getRepoPayload(url string, startPage int) (RepositoryPayload, error) {
	bitbucketUrl := fmt.Sprintf("%s?limit=1000&start=%d", url, startPage)
	log.Debugf("GET %s", bitbucketUrl)
	req, err := http.NewRequest("GET", bitbucketUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth())

	// Send an HTTP GET request to the URL
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error sending GET request")
		return RepositoryPayload{}, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Error("Error: Unexpected status code:", resp.Status)
		return RepositoryPayload{}, errors.New("Status Code: " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("Error reading response body")
		return RepositoryPayload{}, err
	}
	return parseRepositoryPayloadJSON(body), nil
}
