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
	ArangoDB arangoConfig
	LineBot  lineBotConfig
	Zoom     zoomConfig
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
