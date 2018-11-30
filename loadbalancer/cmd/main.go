package main

import (
	"gotil/loadbalancer"
	"time"
)

func main() {
	pool := loadbalancer.New(fib)

	for i := 1; i < 15; i++ {
		pool.Submit(i)
	}
	select {}

}

// Very slow fib calculator
func fib(n interface{}) {
	var fibRec func(int) int64
	fibRec = func(n int) int64 {
		if n <= 1 {
			return int64(n)
		}
		return fibRec(n-1) + fibRec(n-2)
	}
	time.Sleep(time.Second)
	fibRec(n.(int))
}
