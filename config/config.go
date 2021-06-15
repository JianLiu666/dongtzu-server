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
	Fiber    FiberConfig   `yaml:"fiber,omitempty"`
	ArangoDB ArangoConfig  `yaml:"arangoDB,omitempty"`
	LineBot  LineBotConfig `yaml:"linebot,omitempty"`
	Zoom     ZoomConfig    `yaml:"zoom,omitempty"`
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
