package usecase

import "log"

type Semaphore interface {
	Acquire(num int)
	Release(num int)
}

type semaphores struct {
	semC chan struct{}
}

func New(maxConcurrency int) Semaphore {
	return &semaphores{
		semC: make(chan struct{}, maxConcurrency),
	}
}

func (s *semaphores) Acquire(num int) {
	s.semC <- struct{}{}

	log.Println(num, "Acquired!")
}

func (s *semaphores) Release(num int) {
	<-s.semC
	log.Println(num, "Released!")
}
