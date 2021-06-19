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
	addJob("CreateMeetingUrl", "*/10 * * * *", createMeetingUrl)
	addJob("SendMeetingUrl", "*/3 * * * *", sendMeetingUrl)
}

func Start() {
	go func() {
		defer logger.Debugf("[Scheduler] Start.")
		manager.Start()
	}()
}

func addJob(key, spec string, job func()) {
	_, err := manager.AddFunc(spec, func() {
		startAt := time.Now()
		job()
		endAt := time.Now()
		spentTime := (endAt.UnixNano() - startAt.UnixNano()) / 1e6
		logger.Debugf("[Scheduler] %s is start at %v, end at %v, and spent %d ms.", key, startAt.Format("15:04:05.0000"), endAt.Format("15:04:05.0000"), spentTime)
	})

	if err != nil {
		logger.Errorf("[Scheduler] Add job failed: %v", err)
		return
	}
	logger.Debugf("[Scheduler] Add job %v, and spec = %v", key, spec)
}
