package scheduler

import (
	"time"

	"gitlab.geax.io/demeter/gologger/logger"
)

func example() {
	logger.Debugf("[Scheduler] %v", time.Now())
}
