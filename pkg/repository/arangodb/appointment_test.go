package arangodb

import (
	"context"
	"dongtzu/config"
	"dongtzu/pkg/model"
	"testing"

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

func TestCreateAppointment(t *testing.T) {
	initConfig()

	appt := &model.Appointment{
		ID:             "",
		ProviderID:     "1",
		ScheduleID:     "103847",
		ConsumerID:     "1",
		FeedbackID:     "",
		StartTimestamp: 1623776400,
		EndTimestamp:   1623780000,
		Note:           "",
		Status:         0,
	}

	status := CreateAppointment(context.TODO(), appt)
	logger.Debugf("%v", status)
}
