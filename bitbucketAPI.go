package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type repo struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

func parseJSON(jsonBytes []byte) (r repo) {
	err := json.Unmarshal(jsonBytes, &r)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshalling JSON")
	}
	return
}
