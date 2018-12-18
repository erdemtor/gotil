package worker

import (
	"fmt"
	"gotil/loadbalancer/message"
	"sync"
	"sync/atomic"
	"time"
)

var lastID int32

func generateID() string {
	return fmt.Sprintf("worker[%d]", atomic.AddInt32(&lastID, 1))
}

type Worker struct {
	id       string
	mu       sync.RWMutex
	incoming <-chan message.Message
	outgoing chan<- message.Message
	f        func(interface{})
}

func Start(f func(interface{}), incoming, outgoing chan message.Message, count int) {
	for i := 0; i < count; i++ {
		w := &Worker{
			id:       generateID(),
			incoming: incoming,
			outgoing: outgoing,
			f:        f,
		}
		go w.start()
	}
}

func resetTimer() <-chan time.Time {
	return time.After(time.Second * 2)
}

func (w *Worker) start() {
	idleTimer := resetTimer()
	w.outgoing <- message.OfType(message.WorkerStarted)
	for {
		select {
		case <-idleTimer:
			w.outgoing <- message.OfType(message.WorkerDied)
			return
		case msg := <-w.incoming:
			if msg.Type == message.NewTask {
				w.outgoing <- message.OfType(message.TaskStarted)
				w.f(msg.Payload)
				w.outgoing <- message.OfType(message.TaskCompleted)
				idleTimer = resetTimer()
			}

		}
	}
}
