package main

import (
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
	Values []Repo `json:"values"`
}

// Inner wrapper for the JSON repository representation
type Repo struct {
	Archived bool    `json:"archived"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Project  Project `json:"project"`
	Links    struct {
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
	} `json:"links"`
}

type Project struct {
	Key string `json:"key"`
}

func (r Repo) GetSSHRepoUrl() *string {
	return r.extractRepoUrls("ssh")
}

func (r Repo) GetHTTPRepoUrl() *string {
	return r.extractRepoUrls("http")
}

func (r Repo) Delete() error {
	// Create a basic authentication header
	auth := Cfg.BitbucketUsername + ":" + Cfg.BitbucketPassword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	bitbucketDeleteUrl := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s", Cfg.BitbucketUrl, r.Project.Key, r.Slug)

	req, err := http.NewRequest("DELETE", bitbucketDeleteUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth)

	// Send an HTTP GET request to the URL
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error sending GET request")
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is 202 Accepted
	if resp.StatusCode != http.StatusAccepted {
		log.Error("Error: Unexpected status code:", resp.Status)
		return errors.New("Status Code: " + resp.Status)
	}

	return nil
}

func (r Repo) extractRepoUrls(protocol string) *string {
	for _, clone := range r.Links.Clone {
		if clone.Name == protocol {
			return &clone.Href
		}
	}
	return nil
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
	auth := Cfg.BitbucketUsername + ":" + Cfg.BitbucketPassword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	bitbucketUrl := fmt.Sprintf("%s/rest/api/latest/repos?archived=ARCHIVED&limit=%d", Cfg.BitbucketUrl, Cfg.PagingLimit)

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
