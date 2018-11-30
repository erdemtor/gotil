package main

import (
	"gotil/loadbalancer"
	"log"
)

func main() {
	pool := loadbalancer.New(func(data interface{}) {
		fib := fib(data.(int))
		log.Printf("fib(%d) is %d)", data, fib)
	})

	for i := 1; i < 40; i++ {
		pool.Submit(i)
	}
	select {}

}

// Very slow fib calculator
func fib(n int) int64 {
	if n <= 1 {
		return int64(n)
	}
	return fib(n-1) + fib(n-2)
}
