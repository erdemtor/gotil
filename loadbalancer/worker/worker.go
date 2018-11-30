package worker

import (
	"fmt"
	"gotil/loadbalancer/message"
	"gotil/loadbalancer/messenger"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"
)

var lastID int32

func generateID() string {
	return fmt.Sprintf("worker[%d]", atomic.AddInt32(&lastID, 1))
}

type Worker struct {
	mu       sync.RWMutex
	sender   messenger.Sender
	incoming <-chan message.Message
	f        func(interface{})
}

func Start(f func(interface{}), incoming, outgoing chan message.Message, count int) {
	for i := 0; i < count; i++ {
		w := &Worker{
			sender:   messenger.NewSender(generateID(), outgoing),
			incoming: incoming,
			f:        f,
		}
		w.sender.SetLogger(ioutil.Discard)
		go w.start()
	}

}

func resetTimer() <-chan time.Time {
	return time.After(time.Second * 2)
}

func (w *Worker) start() {
	idleTimer := resetTimer()
	w.sender.Send(message.OfType(message.WorkerStarted))
	for {
		select {
		case <-idleTimer:
			w.sender.Send(message.OfType(message.WorkerDied))
			return

		case msg := <-w.incoming:
			if msg.Type == message.NewTask {
				w.sender.Send(message.OfType(message.TaskStarted))
				w.f(msg.Payload)
				w.sender.Send(message.OfType(message.TaskCompleted))
			}
			idleTimer = resetTimer()
		}
	}
}
