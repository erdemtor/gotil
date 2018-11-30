package loadbalancer

import (
	"gotil/loadbalancer/message"
	"gotil/loadbalancer/worker"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type master struct {
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
	m.outgoing <- message.OfType(message.NewTask).WithData(unit)
}

func (m *master) start() {
	for msg := range m.incoming {

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

func New(f func(interface{})) *master {
	m := &master{
		f:           f,
		incoming:    make(chan message.Message, math.MaxInt16),
		outgoing:    make(chan message.Message, math.MaxInt16),
		workerCount: 3,
	}
	go func() {
		for range time.Tick(time.Millisecond * 10) {
			if m.wiq > 100 {
				worker.Start(m.f, m.outgoing, m.incoming, int(m.wiq/2))
			} else {
				if m.workerCount < 3 && m.wiq > 0 {
					worker.Start(m.f, m.outgoing, m.incoming, int(3-m.workerCount))
				}
			}

		}

	}()

	worker.Start(m.f, m.outgoing, m.incoming, 3)
	go m.start()
	return m
}
