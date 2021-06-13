package scheduler

import (
	"github.com/robfig/cron/v3"
	"gitlab.geax.io/demeter/gologger/logger"
)

var manager *cron.Cron

func Init() {
	defer logger.Debugf("[Scheduler] Initialized.")

	manager = cron.New()
	addExample()
}

func Start() {
	go func() {
		manager.Start()
	}()
}
