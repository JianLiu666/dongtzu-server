package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.geax.io/demeter/gologger/logger"
)

var manager *cron.Cron

func Init() {
	defer logger.Debugf("[Scheduler] Initialized.")

	manager = cron.New()
	addJob("updateScheduleAndCreateMeetingUrl", "*/1 * * * *", updateScheduleAndCreateMeetingUrl)
	addJob("processReadyStartAppts", "25,55 * * * *", processReadyStartAppts)
	addJob("processReadyDismissAppts", "25,55 * * * *", processReadyDismissAppts)
}

func Start() {
	go func() {
		defer logger.Debugf("[Scheduler] Start.")
		manager.Start()
	}()
}

func addJob(key, spec string, job func()) {
	_, err := manager.AddFunc(spec, func() {
		time.Now().UTC().Unix()
		startAt := time.Now().UnixNano()
		logger.Debugf("[Scheduler] %s now is working, and spec = %s", key, spec)
		job()
		spentTime := (time.Now().UnixNano() - startAt) / 1e6
		logger.Debugf("[Scheduler] %s now has done, it spent %d (ms) and spec = %s", key, spentTime, spec)
	})

	if err != nil {
		logger.Errorf("[Scheduler] add job failed: %v", err)
		return
	}
	logger.Debugf("[Scheduler] add job %v, and spec = %v", key, spec)
}
