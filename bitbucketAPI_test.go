package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {
	apiResponse := []byte(`
    {
        "size": 50,
        "limit": 50,
        "isLastPage": false,
        "values": [
		  {
            "slug": "00-template",
            "id": 1128,
            "name": "00-template",
            "hierarchyId": "f13df5fb2c96bbc4a1b9",
            "scmId": "git",
            "state": "AVAILABLE",
            "statusMessage": "Available",
            "forkable": true,
            "project": {
                "key": "~POM15",
                "id": 1261,
                "name": "Jonas Sticha",
                "type": "PERSONAL",
                "owner": {
                    "name": "pom15",
                    "emailAddress": "test@tum.de",
                    "active": false,
                    "displayName": "Test",
                    "id": 3388,
                    "slug": "pom15",
                    "type": "NORMAL",
                    "links": {
                        "self": [
                            {
                                "href": "https://bitbucket.example.com/users/pom15"
                            }
                        ]
                    }
                },
                "links": {
                    "self": [
                        {
                            "href": "https://bitbucket.example.com/users/pom15"
                        }
                    ]
                }
            },
            "public": false,
            "archived": true,
            "links": {
                "clone": [
                    {
                        "href": "ssh://git@bitbucket.example.com:7999/~pom15/00-template.git",
                        "name": "ssh"
                    },
                    {
                        "href": "https://bitbucket.example.com/scm/~pom15/00-template.git",
                        "name": "http"
                    }
                ],
                "self": [
                    {
                        "href": "https://bitbucket.example.com/users/pom15/repos/00-template/browse"
                    }
                ]
            }
        }
      ]
    }`)

	r := parseRepositoryPayloadJSON(apiResponse).Values[0]
	assert.Equal(t, r.Name, "00-template")
	assert.Equal(t, r.Archived, true)
	assert.Len(t, r.Links.Clone, 2)
	assert.Equal(t, r.Links.Clone[0].Href, "ssh://git@bitbucket.example.com:7999/~pom15/00-template.git")
	assert.Equal(t, r.Links.Clone[1].Href, "https://bitbucket.example.com/scm/~pom15/00-template.git")
}

func TestUrlExtraction(t *testing.T) {
	apiResponse := []byte(`
    {
        "size": 50,
        "limit": 50,
        "isLastPage": false,
        "values": [
		  {
            "slug": "00-template",
            "id": 1128,
            "name": "00-template",
            "hierarchyId": "f13df5fb2c96bbc4a1b9",
            "scmId": "git",
            "state": "AVAILABLE",
            "statusMessage": "Available",
            "forkable": true,
            "project": {
                "key": "~POM15",
                "id": 1261,
                "name": "Jonas Sticha",
                "type": "PERSONAL",
                "owner": {
                    "name": "pom15",
                    "emailAddress": "test@tum.de",
                    "active": false,
                    "displayName": "Test",
                    "id": 3388,
                    "slug": "pom15",
                    "type": "NORMAL",
                    "links": {
                        "self": [
                            {
                                "href": "https://bitbucket.example.com/users/pom15"
                            }
                        ]
                    }
                },
                "links": {
                    "self": [
                        {
                            "href": "https://bitbucket.example.com/users/pom15"
                        }
                    ]
                }
            },
            "public": false,
            "archived": true,
            "links": {
                "clone": [
                    {
                        "href": "ssh://git@bitbucket.example.com:7999/~pom15/00-template.git",
                        "name": "ssh"
                    },
                    {
                        "href": "https://bitbucket.example.com/scm/~pom15/00-template.git",
                        "name": "http"
                    }
                ],
                "self": [
                    {
                        "href": "https://bitbucket.example.com/users/pom15/repos/00-template/browse"
                    }
                ]
            }
        }
      ]
    }`)

	r := parseRepositoryPayloadJSON(apiResponse).Values[0]
	httpsUrl := extractRepoUrls([]Repo{r}, false)
	assert.Len(t, httpsUrl, 1)
	assert.Equal(t, httpsUrl[0], "https://bitbucket.example.com/scm/~pom15/00-template.git")

	sshUrl := extractRepoUrls([]Repo{r}, true)
	assert.Len(t, sshUrl, 1)
	assert.Equal(t, sshUrl[0], "ssh://git@bitbucket.example.com:7999/~pom15/00-template.git")
}
