package bitbucket

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"bitbucket_archiver/utils"

	log "github.com/sirupsen/logrus"
)

var client = http.Client{}

// Outer wrapper for the JSON payload
type RepositoryPayload struct {
	Values []Repo `json:"values"`
}

func parseRepositoryPayloadJSON(jsonBytes []byte) (r RepositoryPayload) {
	err := json.Unmarshal(jsonBytes, &r)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshalling JSON")
	}
	return
}

func GetArchivedRepositories() ([]Repo, error) {
	// Create a basic authentication header
	auth := utils.Cfg.BitbucketUsername + ":" + utils.Cfg.BitbucketPassword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	bitbucketUrl := fmt.Sprintf("%s/rest/api/latest/repos?archived=ARCHIVED&limit=%d", utils.Cfg.BitbucketUrl, utils.Cfg.PagingLimit)

	req, err := http.NewRequest("GET", bitbucketUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth)

	// Send an HTTP GET request to the URL
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error sending GET request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Error("Error: Unexpected status code:", resp.Status)
		return nil, errors.New("Status Code: " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("Error reading response body")
		return nil, err
	}

	return parseRepositoryPayloadJSON(body).Values, nil
}

func GetRepositoriesForProject(projectKey string) ([]Repo, error) {
	// Create a basic authentication header
	auth := utils.Cfg.BitbucketUsername + ":" + utils.Cfg.BitbucketPassword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	bitbucketUrl := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos?limit=%d", utils.Cfg.BitbucketUrl, projectKey, utils.Cfg.PagingLimit)

	req, err := http.NewRequest("GET", bitbucketUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth)

	// Send an HTTP GET request to the URL
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error sending GET request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Error("Error: Unexpected status code:", resp.Status)
		return nil, errors.New("Status Code: " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("Error reading response body")
		return nil, err
	}

	return parseRepositoryPayloadJSON(body).Values, nil
}
