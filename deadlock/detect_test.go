package deadlock

import (
	"log"
	"os"
	"testing"
)

type Server struct {
	lock RWMutex
	svc  *Service
	c    chan int
}

func (s *Server) dolock() {
	s.lock.Lock()
	go s.svc.dotest()
	<-s.c
}

type Service struct {
	server *Server
}

func (s *Service) dotest() {
	s.server.lock.RLock()
	s.server.c <- 100
}

func TestNormal(t *testing.T) {
	c := make(chan int, 1)
	var lock Mutex
	lock.Lock()
	go func() {
		lock.Lock()
		c <- 1
		lock.Unlock()
	}()
	lock.Unlock()
	lock.Lock()
	lock.Unlock()
	<-c

	var rwlock RWMutex
	rwlock.RLock()
	go func() {
		lock.Lock()
		c <- 1
		lock.Unlock()
	}()
	rwlock.RUnlock()
	lock.Lock()
	lock.Unlock()
	<-c
}

func TestDirectReLock(t *testing.T) {
	var lock Mutex
	Opts.OnDeadlock = func() {
		log.Println("---catch lock deadlock info---")
		lock.Unlock()
	}
	lock.Lock()
	lock.Unlock()

	lock.Lock()

	lock.Lock()
	lock.Unlock()

	var rwlock RWMutex
	Opts.OnDeadlock = func() {
		log.Println("---catch rwlock deadlock info---")
		rwlock.RUnlock()
	}
	rwlock.Lock()
	rwlock.Unlock()

	rwlock.RLock()

	rwlock.Lock()
	rwlock.Unlock()
}

func TestReLock(t *testing.T) {
	Opts.OnDeadlock = func() {
		log.Println("---catch server deadlock info---")
		os.Exit(0)
	}
	s := &Server{
		c: make(chan int),
	}
	s.svc = &Service{server: s}
	s.dolock()
}
