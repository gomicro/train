package config

import (
	"fmt"
	"io"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	confDir  = "/.train"
	confFile = "/config"
)

var defaultConfig = Config{
	ReleaseBranch: "release",
	Github: &GithubHost{
		Limits: &Limits{
			RequestsPerSecond: 10,
			Burst:             25,
		},
		Ignores: &GithubIgnores{},
	},
}

// Config represents the config file for train
type Config struct {
	ReleaseBranch string      `yaml:"release_branch"`
	Github        *GithubHost `yaml:"github.com"`
}

// Limits represents a limits override for the client
type Limits struct {
	RequestsPerSecond int `yaml:"request_per_second"`
	Burst             int `yaml:"burst"`
}

// New takes a token string and creates the most basic config capable of being
// written.
func New(tkn string) *Config {
	return &Config{Github: &GithubHost{Token: tkn}}
}

// WriteFile writes the file to the defined location for the current user, and
// returns any errors encountered doing so.
func (c *Config) WriteFile() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("config: marshal: %v", err.Error())
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("config: get home directory: %v", err.Error())
	}

	err = io.WriteFile(filepath.Join(usr.HomeDir, confDir, confFile), b, 0600)
	if err != nil {
		return fmt.Errorf("config: write file: %v", err.Error())
	}

	return nil
}
