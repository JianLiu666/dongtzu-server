package scheduler

import (
	"dongtzu/pkg/model"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.geax.io/demeter/gologger/logger"
)

var manager *cron.Cron
var timeWheelIdx int
var startTimeWheel [60][]*model.Appointment
var endTimeWheel [60][]*model.Appointment

func Init() {
	defer logger.Debugf("[Scheduler] Initialized.")

	manager = cron.New()
	addJob("example", "@every 1m", example)
	addJob("getAndConfirmAppts", "0 25,55 * * * *", getAndConfirmAppts)
	addJob("processAppts", "0 10 * * * *", processAppts)
}

func Start() {
	go func() {
		defer logger.Debugf("[Scheduler] Start.")
		manager.Start()
	}()
}

func addJob(key, spec string, job func()) {
	_, _ = manager.AddFunc(spec, func() {
		time.Now().UTC().Unix()
		startAt := time.Now().UnixNano()
		logger.Debugf("[Scheduler] %s now is working, and spec = %s", key, spec)
		job()
		spentTime := (time.Now().UnixNano() - startAt) / 1e6
		logger.Debugf("[Scheduler] %s now has done, it spent %d (ms) and spec = %s", key, spentTime, spec)
	})
}
