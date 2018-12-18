package main

import (
	"gotil/loadbalancer"
	"log"
	"time"
)

func main() {
	pool := loadbalancer.New(consumer)

	for i := 0; i < 15; i++ {
		pool.Submit(i)
		log.Printf("%d is submitted", i)
		if i == 10 {
			time.Sleep(time.Second * 5)
		}
	}

	select {}

}

//
//go func() {
//	for range time.Tick(time.Second) {
//		log.Printf("(worker, wip, wiq) (%d,%d,%d) - GO: %d\n", atomic.LoadInt32(&m.workerCount), atomic.LoadInt32(&m.wip), atomic.LoadInt32(&m.wiq), runtime.NumGoroutine())
//	}
//}()

// Very slow consumer calculator
func consumer(nIn interface{}) {
	n := nIn.(int)
	time.Sleep(time.Second * time.Duration(n))
}
