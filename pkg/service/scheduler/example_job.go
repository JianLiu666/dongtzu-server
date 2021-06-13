package scheduler

import (
	"fmt"
	"time"
)

func addExample() {
	_, err := manager.AddFunc("@every 1s", func() {
		fmt.Println(time.Now().Clock())
	})
	if err != nil {
		fmt.Printf("add cronjob failed: %v\n", err)
	}
}
