package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	var apiResponse = []byte(`
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
                "key": "~POM15_JONAS.STICHA",
                "id": 1261,
                "name": "Jonas Sticha",
                "type": "PERSONAL",
                "owner": {
                    "name": "pom15_jonas.sticha",
                    "emailAddress": "jonas.sticha@tum.de",
                    "active": false,
                    "displayName": "Jonas Sticha",
                    "id": 3388,
                    "slug": "pom15_jonas.sticha",
                    "type": "NORMAL",
                    "links": {
                        "self": [
                            {
                                "href": "https://bitbucket.ase.in.tum.de/users/pom15_jonas.sticha"
                            }
                        ]
                    }
                },
                "links": {
                    "self": [
                        {
                            "href": "https://bitbucket.ase.in.tum.de/users/pom15_jonas.sticha"
                        }
                    ]
                }
            },
            "public": false,
            "archived": true,
            "links": {
                "clone": [
                    {
                        "href": "ssh://git@bitbucket.ase.in.tum.de:7999/~pom15_jonas.sticha/00-template.git",
                        "name": "ssh"
                    },
                    {
                        "href": "https://bitbucket.ase.in.tum.de/scm/~pom15_jonas.sticha/00-template.git",
                        "name": "http"
                    }
                ],
                "self": [
                    {
                        "href": "https://bitbucket.ase.in.tum.de/users/pom15_jonas.sticha/repos/00-template/browse"
                    }
                ]
            }
        }`)

	r := parseJSON(apiResponse)
	assert.Equal(t, r.Name, "00-template")
	assert.Equal(t, r.Archived, true)
}
