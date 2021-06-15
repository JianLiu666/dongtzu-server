package config

type GithubConfig struct {
	RepoURL  string `yaml:"repoURL,omitempty"`
	APIToken string `yaml:"apiToken,omitempty"`
}
