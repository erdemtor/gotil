package worker

import (
	"fmt"
	"gotil/loadbalancer/message"
	"gotil/loadbalancer/messenger"
	"sync"
	"sync/atomic"
	"time"
)

var lastID int32

func generateID() string {
	return fmt.Sprintf("worker[%d]", atomic.AddInt32(&lastID, 1))
}

type Worker struct {
	mu sync.RWMutex
	messenger.Messenger
	f func(interface{})
}

func Start(f func(interface{}), incoming, outgoing chan message.Message, count int) {
	for i := 0; i < count; i++ {
		w := &Worker{
			Messenger: messenger.New(generateID(), incoming, outgoing),
			f:         f,
		}
		go w.start()
	}

}

func resetTimer() <-chan time.Time {
	return time.After(time.Second * 2)
}

func (w *Worker) start() {
	idleTimer := resetTimer()
	w.Send(message.OfType(message.WorkerStarted))
	for {
		select {
		case <-idleTimer:
			w.Send(message.OfType(message.WorkerDied))
			return

		case msg := <-w.Receive():
			if msg.Type == message.NewTask {
				w.Send(message.OfType(message.TaskStarted))
				w.f(msg.Payload)
				w.Send(message.OfType(message.TaskCompleted))
			}
			idleTimer = resetTimer()
		}
	}
}
