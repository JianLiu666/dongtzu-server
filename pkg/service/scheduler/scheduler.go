package scheduler

import (
	"dongtzu/pkg/model"

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
	_, _ = manager.AddFunc("@every 1m", example)
	_, _ = manager.AddFunc("0 25,55 * * * *", getAndConfirmAppts)
	_, _ = manager.AddFunc("0 10 * * * *", processAppts)
}

func Start() {
	go func() {
		manager.Start()
	}()
}
