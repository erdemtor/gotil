package main

import (
	"gotil/loadbalancer"
	"log"
	"time"
)

func main() {
	pool := loadbalancer.New(fib)

	for i := 0; i < 15; i++ {
		pool.Submit(i)
		log.Printf("%d is submitted", i)
		if i == 10 {
			time.Sleep(time.Second * 5)
		}
	}
	select {}

}

// Very slow fib calculator
func fib(nIn interface{}) {
	n := nIn.(int)
	time.Sleep(time.Second * time.Duration(n))
}
