package main

import (
	"gotil/loadbalancer"
	"time"
)

func main() {
	pool := loadbalancer.New(consumer)

	//go func() {
	//	for range time.Tick(time.Second) {
	//		log.Printf("go-routine count: %d", runtime.NumGoroutine())
	//	}
	//}()

	for i := 0; i < 15; i++ {
		pool.Submit(i)
		time.Sleep(time.Second)
	}

	select {}

}

// Very slow consumer calculator
func consumer(nIn interface{}) {
	n := nIn.(int)
	time.Sleep(time.Second * time.Duration(n))
}
