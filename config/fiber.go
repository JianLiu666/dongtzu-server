package config

type FiberConfig struct {
	Port string `yaml:"port,omitempty"`
}

type GithubConfig struct {
	RepoURL  string `yaml:"repoURL,omitempty"`
	APIToken string `yaml:"apiToken,omitempty"`
}
