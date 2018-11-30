package loadbalancer

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMaster_Submit(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(20000)

	balancer := New(func(data interface{}) {
		time.Sleep(time.Millisecond * time.Duration(data.(int)))
		wg.Done()
	})

	go func() {
		for range time.Tick(time.Millisecond * 250) {
			log.Printf("(worker, wip, wiq) (%d,%d,%d) \n", atomic.LoadInt32(&balancer.workerCount), atomic.LoadInt32(&balancer.wip), atomic.LoadInt32(&balancer.wiq))
		}
	}()
	for i := 0; i < 3000; i++ {
		balancer.Submit(i)
		if i%100 == 0 {
			time.Sleep(time.Millisecond * 500)
			log.Println("Submit Count:", i)
		}
	}

	wg.Wait()
	log.Println("all tasks are completed")
	time.Sleep(time.Second * 20)

}
