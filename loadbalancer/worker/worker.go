package worker

import (
	"gotil/loadbalancer/message"
	"sync"
	"time"
)

type Worker struct {
	mu       sync.RWMutex
	incoming chan message.Message
	outgoing chan message.Message
	f        func(interface{})
}

func Start(f func(interface{}), incoming, outgoing chan message.Message, count int) {
	for i := 0; i < count; i++ {
		go (&Worker{
			incoming: incoming,
			outgoing: outgoing,
			f:        f,
		}).start()
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
			idleTimer = resetTimer()
			if msg.Type == message.NewTask {
				w.outgoing <- message.OfType(message.TaskStarted)
				w.f(msg.Data)
				w.outgoing <- message.OfType(message.TaskCompleted)
			}
		}
	}
}
