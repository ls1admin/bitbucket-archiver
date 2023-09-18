package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Outer wrapper for the JSON payload
type RepositoryPayload struct {
	Values []Repo `json:"values"`
}

// Inner wrapper for the JSON repository representation
type Repo struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
	Links    struct {
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
	} `json:"links"`
}

func ExtractRepoUrls(repositories []Repo, useSSH bool) []string {
	var urls []string
	// Filter for SSH repositories
	for _, repo := range repositories {
		for _, clone := range repo.Links.Clone {
			if clone.Name == "ssh" && useSSH {
				urls = append(urls, clone.Href)
			}
			if clone.Name == "http" && !useSSH {
				urls = append(urls, clone.Href)
			}
		}
	}
	return urls
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

	client := http.Client{}
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
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("Error reading response body")
		return nil, err
	}

	return parseRepositoryPayloadJSON(body).Values, nil
}
