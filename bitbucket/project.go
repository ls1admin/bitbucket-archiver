package bitbucket

import (
	"bitbucket_archiver/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
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

func (r Repo) GetSSHRepoUrl() *string {
	return r.extractRepoUrls("ssh")
}

func (r Repo) GetHTTPRepoUrl() *string {
	return r.extractRepoUrls("http")
}

func (r Repo) GetSize() int64 {
	sizeUrl := fmt.Sprintf("%s/projects/%s/repos/%s/sizes", utils.Cfg.BitbucketUrl, r.Project.Key, r.Slug)

	req, err := http.NewRequest("GET", sizeUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth())

	// Send an HTTP GET request to the URL
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error sending GET request")
		return 0
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Error("Error: Unexpected status code:", resp.Status)
		return 0
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("Error reading response body")
		return 0
	}

	sizeInfo := struct {
		RepositorySize int64 `json:"repository"`
	}{}
	err = json.Unmarshal(body, &sizeInfo)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshalling JSON")
	}

	return sizeInfo.RepositorySize
}

func (r Repo) Delete() error {
	bitbucketDeleteUrl := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s", utils.Cfg.BitbucketUrl, r.Project.Key, r.Slug)

	req, err := http.NewRequest("DELETE", bitbucketDeleteUrl, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating request")
	}

	// Set the Authorization header for basic authentication
	req.Header.Add("Authorization", basicAuth())

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
