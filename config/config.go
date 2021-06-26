package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	setOnce      sync.Once
	globalConfig *config
)

type config struct {
	ArangoDB ArangoConfig   `yaml:"arangoDB"`
	Fiber    FiberConfig    `yaml:"fiber"`
	Github   GithubConfig   `yaml:"github"`
	Zoom     ZoomConfig     `yaml:"zoom"`
	NewebPay NewebPayConfig `yaml:"newebPay"`
}

func NewFromViper() (*config, error) {
	var c config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func SetConfig(c *config) {
	setOnce.Do(func() {
		globalConfig = c
	})
}

func GetGlobalConfig() *config {
	return globalConfig
}

type ArangoConfig struct {
	Addr       string `yaml:"addr"`
	DBName     string `yaml:"dbName"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ConnLimit  int    `yaml:"connLimit"`
	RetryCount int    `yaml:"retryCount"`
}

type FiberConfig struct {
	Port string `yaml:"port"`
}

type GithubConfig struct {
	RepoURL  string `yaml:"repoURL"`
	APIToken string `yaml:"apiToken"`
}

type ZoomConfig struct {
	MeetingExtendedTime int `yaml:"meetingExtendedTime"`
}

type NewebPayConfig struct {
	APIUrl          string `yaml:"apiUrl"`
	APIVersion      string `yaml:"apiVersion"`
	MerchantID      string `yaml:"merchantID"`
	MerchantHashKey string `yaml:"merchantHashKey"`
	MerchantIV      string `yaml:"merchantIV"`
}
