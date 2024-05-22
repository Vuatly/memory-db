package utils

import "sync"

type Semaphore struct {
	current int
	max     int
	cond    *sync.Cond
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		max:  max,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	if s.current >= s.max {
		s.cond.Wait()
	}

	s.current++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.current--
	s.cond.Signal()
}
