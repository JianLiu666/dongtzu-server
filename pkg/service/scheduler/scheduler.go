package scheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var manager *cron.Cron

func Init() {
	defer fmt.Println("Scheduler Initialized.")

	manager = cron.New()
	addExample()
}

func Start() {
	go func() {
		manager.Start()
	}()
}
