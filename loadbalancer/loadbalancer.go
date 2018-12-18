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
	wip         int32
	wiq         int32
	workerCount int32
}

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
		case message.TaskCompleted:
			atomic.AddInt32(&m.wip, -1)
		case message.WorkerDied:
			atomic.AddInt32(&m.workerCount, -1)
		case message.TaskStarted:
			atomic.AddInt32(&m.wiq, -1)
			atomic.AddInt32(&m.wip, 1)
		}
	}
}

func (m *master) StartWorker(count int) {
	worker.Start(m.f, m.outgoing, m.incoming, count)
}

func New(f func(interface{})) Balancer {
	outgoing := make(chan message.Message, math.MaxInt16)
	incoming := make(chan message.Message, math.MaxInt16)
	m := &master{
		id:       generateID(),
		f:        f,
		incoming: incoming,
		outgoing: outgoing,
	}
	go func() {
		for range time.Tick(time.Millisecond * 10) {
			if m.wiq > m.workerCount || m.workerCount == 0 {
				m.StartWorker(int(m.wiq))
			}
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			log.Printf("(worker, wip, wiq) (%d,%d,%d) - GO: %d\n", atomic.LoadInt32(&m.workerCount), atomic.LoadInt32(&m.wip), atomic.LoadInt32(&m.wiq), runtime.NumGoroutine())
		}
	}()

	m.StartWorker(1)
	go m.start()
	return m
}
