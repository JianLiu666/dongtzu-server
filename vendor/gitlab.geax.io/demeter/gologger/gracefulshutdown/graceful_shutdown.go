package gracefulshutdown

import (
	"context"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
)

var ctxs []context.Context
var wgs []*sync.WaitGroup
var shutdownChan chan os.Signal

// Start return 的 chan os.Signal 是 for test 使用的，正常使用不需要拿來處理。
func Start() chan os.Signal {
	shutdownChan = make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	cancels := make([]context.CancelFunc, math.MaxUint8)
	ctxs = make([]context.Context, math.MaxUint8)
	wgs = make([]*sync.WaitGroup, math.MaxUint8)

	for i := 0; i < math.MaxUint8; i++ {
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
		wgs[i] = new(sync.WaitGroup)
	}
	go func() {
		<-shutdownChan
		signal.Stop(shutdownChan)
		for i := 0; i < math.MaxUint8; i++ {
			cancels[i]()
			wgs[i].Wait()
		}
		log.Fatal("GracefulShutdown Done")
	}()
	return shutdownChan
}

// GetContext
// 取得系統通知準備關係的 context
// 取得要通知可以關閉的 chan
func GetContext(level uint8) (context.Context, chan struct{}) {
	c := make(chan struct{})
	wgs[level].Add(1)
	go func(inc chan struct{}, inwg *sync.WaitGroup) {
		<-inc
		inwg.Done()
	}(c, wgs[level])
	return ctxs[level], c
}

func Shutdown() {
	shutdownChan <- os.Interrupt
}
