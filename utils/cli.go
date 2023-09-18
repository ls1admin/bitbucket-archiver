package utils

import (
	"flag"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GitUsername string `env:"GIT_USERNAME,notEmpty"`
	GitPassword string `env:"GIT_PASSWORD,notEmpty"`

	BitbucketUrl      string `env:"BITBUCKET_URL,notEmpty"`
	BitbucketUsername string `env:"BITBUCKET_USERNAME,notEmpty"`
	BitbucketPassword string `env:"BITBUCKET_PASSWORD,notEmpty"`

	UseSSHCloning bool `env:"USE_SSH_CLONING" envDefault:"false"`
	PagingLimit   int  `env:"PAGING_LIMIT" envDefault:"100"`

	OutputDir string `env:"OUTPUT_DIR" envDefault:"./repos_zipped"`
	CloneDir  string `env:"CLONE_DIR" envDefault:"./repos_cloned"`

	DeleteRepos bool // No env var for this one, just a flag

	ProjectFile string
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
	flag.StringVar(&Cfg.ProjectFile, "project-file", "", "Path to project file")
	flag.Parse()

	log.Debug("Config loaded: ", Cfg)
}
