package semaphore

type Semaphore interface {
	Lock()
	UnLock()
}
type semaphore struct {
	c chan struct{}
}

func (s *semaphore) Lock() {
	s.c <- struct{}{}
}

func (s *semaphore) UnLock() {
	<-s.c
}

func New(n int) Semaphore {
	return &semaphore{
		c: make(chan struct{}, n),
	}
}
