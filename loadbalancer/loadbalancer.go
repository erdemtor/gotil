package loadbalancer

import (
	"fmt"
	"gotil/loadbalancer/message"
	"gotil/loadbalancer/worker"
	"math"
	"sync"
	"sync/atomic"
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
			continue
		case message.TaskCompleted:
			atomic.AddInt32(&m.wip, -1)
		case message.WorkerDied:
			atomic.AddInt32(&m.workerCount, -1)
		case message.TaskStarted:
			atomic.AddInt32(&m.wiq, -1)
			atomic.AddInt32(&m.wip, 1)
		}
		if m.wiq > m.workerCount || m.workerCount == 0 {
			m.StartWorker(int(m.wiq))
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
	m.StartWorker(1)
	go m.start()
	return m
}
