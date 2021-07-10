package pool

import (
	"log"
	"math/rand"
	"sync/atomic"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

type Pool struct {
	workerPool chan chan service.IRunnable
	queue      chan service.IRunnable
	running    int64
	resetLimit int
}

func (p *Pool) AddJob(job service.IRunnable) {
	p.queue <- job
}

func NewPool(workerSize, queueSize int) *Pool {
	return &Pool{
		workerPool: make(chan chan service.IRunnable, workerSize),
		queue:      make(chan service.IRunnable, queueSize),
		running:    0,
		resetLimit: 5000,
	}
}

func (p *Pool) Start() {
	log.Printf("info: starting worker pool of size: %v", cap(p.workerPool))
	for i := 0; i < cap(p.workerPool); i++ {
		p.newWorker(p.workerPool, i+1)
	}
	for {
		job, active := <-p.queue
		if !active {
			log.Println("info: pool exiting")
			return
		}
		worker := <-p.workerPool
		worker <- job
	}
}

func (p *Pool) newWorker(workerPool chan chan service.IRunnable, i int) {
	workerChannel := make(chan service.IRunnable)
	workerPool <- workerChannel
	go func() {
		for done := 0; done < p.resetLimit+rand.Intn(p.resetLimit); done++ {
			job, active := <-workerChannel
			if !active {
				log.Printf("info: worker %v exiting", i)
				return
			}
			atomic.AddInt64(&p.running, 1)
			job.Run()
			atomic.AddInt64(&p.running, -1)
			workerPool <- workerChannel
		}
		log.Printf("info: worker %v resetting", i)
		go p.newWorker(workerPool, i)
	}()
}

func (p *Pool) IsFull() bool {
	running := atomic.LoadInt64(&p.running)
	if running >= int64(cap(p.queue)) {
		return true
	}
	return false
}
