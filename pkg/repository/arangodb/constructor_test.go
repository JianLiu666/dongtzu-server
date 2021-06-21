package arangodb

import (
	"dongtzu/config"

	"github.com/spf13/viper"
	"gitlab.geax.io/demeter/gologger/logger"
)

func initConfig() {
	logger.Init("debug")

	viper.SetConfigFile("./../../../conf.d/env.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("ReadInConfig file failed: %v", err)
	} else {
		logger.Debugf("Using config file: %v", viper.ConfigFileUsed())
	}

	c, err := config.NewFromViper()
	if err != nil {
		logger.Errorf("Init config failed: %v", err)
	}
	config.SetConfig(c)

	Init()
}
