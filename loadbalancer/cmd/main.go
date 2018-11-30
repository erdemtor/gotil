package main

import (
	"gotil/loadbalancer"
	"log"
	"time"
)

func main() {
	pool := loadbalancer.New(func(data interface{}) {
		time.Sleep(time.Second * time.Duration(data.(int)))
		log.Println(data)
	})

	for i := 1; i < 15; i++ {
		pool.Submit(i)
	}
	select {}

}
