package main

import (
	"gotil/loadbalancer"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	pool := loadbalancer.New(consumer)
	for i := 0; i < 150; i++ {
		pool.Submit(i)
		log.Println(i)
	}

	select {}

}

// Very slow consumer calculator
func consumer(_ interface{}) {

	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
}
