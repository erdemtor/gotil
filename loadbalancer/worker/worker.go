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

//Worker is a representation of a goroutine communicating with the master.
//After it's started it waits for work units and kill itself when being left idle.
type Worker struct {
	id       string
	mu       sync.RWMutex
	incoming <-chan message.Message
	outgoing chan<- message.Message
	f        func(interface{})
}

//New creates a new *Worker by setting the given parameters
func New(f func(interface{}), incoming, outgoing chan message.Message) *Worker {
	return &Worker{
		id:       generateID(),
		incoming: incoming,
		outgoing: outgoing,
		f:        f,
	}
}

func newIdleWaiter() <-chan time.Time {
	return time.After(time.Second * 2)
}

//Start starts a goroutine and starts waiting for tasks to be submitted from the master
//If nothing is submitted for 2 seconds it'll finalise the routine.
func (w *Worker) Start() {
	go func() {
		idleTimer := newIdleWaiter()
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
					idleTimer = newIdleWaiter()
				}

			}
		}
	}()
}
