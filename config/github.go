package config

// GithubHost represents a single host for which train has a configuration
type GithubHost struct {
	Token   string         `yaml:"token"`
	Ensures *GithubEnsures `yaml:"ensures"`
	Ignores *GithubIgnores `yaml:"ignores"`
	Limits  *Limits        `yaml:"limits"`
}

type GithubEnsures struct {
	Repos  []string `yaml:"repos"`
	Topics []string `yaml:"topics"`
}

type GithubIgnores struct {
	Repos  []string `yaml:"repos"`
	Topics []string `yaml:"topics"`
}
