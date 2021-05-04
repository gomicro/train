package config

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	filename = "/.train/config"
)

// Config represents the config file for train
type Config struct {
	Github Host `yaml:"github.com"`
}

// Host represents a single host that train has a configuration for
type Host struct {
	Token string `yaml:"token"`
}

// ParseFromFile reads the train config file from the home directory. It returns
// any errors it encounters with parsing the file.
func ParseFromFile() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Failed getting home directory: %v", err.Error())
	}

	b, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, filename))
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err.Error())
	}

	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config file: %v", err.Error())
	}

	return &conf, nil
}
