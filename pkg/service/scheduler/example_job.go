package scheduler

import (
	"time"

	"gitlab.geax.io/demeter/gologger/logger"
)

func addExample() {
	_, err := manager.AddFunc("@every 1s", func() {
		logger.Debugf("[Scheduler] %v", time.Now())
	})
	if err != nil {
		logger.Errorf("[Scheduler] add cronjob failed: %v\n", err)
	}
}
