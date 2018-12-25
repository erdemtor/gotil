package loadbalancer

import (
	"fmt"
	"gotil/loadbalancer/message"
	"gotil/loadbalancer/worker"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var lastID int32

func generateID() string {
	return fmt.Sprintf("master[%d]", atomic.AddInt32(&lastID, 1))
}

type Balancer interface {
	Submit(unit interface{})
}

type master struct {
	id          string
	mu          sync.Mutex
	f           func(interface{})
	incoming    chan message.Message
	outgoing    chan message.Message
	wip         int32 // work in process
	wiq         int32 // work in queue
	workerCount int32
}

//Submit allows the client to add more tasks into the queue of tasks and can influence the number of goroutines running
func (m *master) Submit(unit interface{}) {
	atomic.AddInt32(&m.wiq, 1)
	m.outgoing <- message.OfType(message.NewTask).WithPayload(unit)
}

func (m *master) start() {
	for {
		msg := <-m.incoming
		switch msg.Type {
		case message.WorkerStarted:
			atomic.AddInt32(&m.workerCount, 1)
			continue
		case message.TaskCompleted:
			atomic.AddInt32(&m.wip, -1)
		case message.WorkerDied:
			atomic.AddInt32(&m.workerCount, -1)
		case message.TaskStarted:
			atomic.AddInt32(&m.wiq, -1)
			atomic.AddInt32(&m.wip, 1)
		}
		if m.wiq >= m.workerCount || m.workerCount == 0 {
			m.startWorker() // add more workers if the queue length is more than worker count
		}
	}
}

func (m *master) startWorker() {
	w := worker.New(m.f, m.outgoing, m.incoming)
	w.Start()
}

//New creates an initialises a new balancer, ready to submit work units to be run
func New(f func(interface{})) Balancer {
	outgoing := make(chan message.Message, math.MaxInt16)
	incoming := make(chan message.Message, math.MaxInt16)
	m := &master{
		id:       generateID(),
		f:        f,
		incoming: incoming,
		outgoing: outgoing,
	}
	m.startWorker() // to have one worker at least waiting for tasks
	go m.start()
	go func() { // logging purposes
		for range time.Tick(time.Second) {
			log.Printf("(Worker Count, WIP, WIQ) (%d, %d, %d) go: %d", m.workerCount, m.wip, m.wiq, runtime.NumGoroutine())
		}
	}()
	return m
}
