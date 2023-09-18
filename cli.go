package main

import (
	"flag"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GitUsername string `env:"GIT_USERNAME,notEmpty"`
	GitPassword string `env:"GIT_PASSWORD,notEmpty"`

	OutputDir string `env:"OUTPUT_DIR,notEmpty"`

	DeleteRepos bool // No env var for this one, just a flag
}

var Cfg Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Warn("Error loading .env file")
	}

	err = env.Parse(&Cfg)
	if err != nil {
		log.WithError(err).Fatal("Error parsing environment variables")
	}

	flag.BoolVar(&Cfg.DeleteRepos, "execute-delete", false, "Delete repos after cloning")
	flag.Parse()

	log.Debug("Config loaded: ", Cfg)
}
