package sail

import (
	"sync"
	"sync/atomic"
	"time"
)

type pool[I, O any] struct {
	yield   func(I, int) O
	input   chan I
	output  chan<- O
	max     int32
	min     int32
	workers *sync.Pool
	count   atomic.Int32
	idle    time.Duration
	busy    time.Duration
	workID  int
}

func (p *pool[I, O]) Post(elems ...I) {
	timer := time.NewTimer(p.busy)

	for _, elem := range elems {
		select {
		case p.input <- elem:
		case <-timer.C:
			p.up()
			p.input <- elem
		}

		timer.Reset(p.busy)
	}
}

func (p *pool[I, O]) Workers() int32 {
	return p.count.Load()
}

func (p *pool[I, O]) Cap() int {
	return cap(p.input)
}

func (p *pool[I, O]) Len() int {
	return len(p.input)
}

func (p *pool[I, O]) Close() {
	close(p.input)
}

func (p *pool[I, O]) newWorker() any {
	p.workID++

	return &worker[I, O]{pool: p, id: p.workID}
}

func (p *pool[I, O]) up() {
	if p.count.Load() > p.max {
		return
	}

	work, _ := p.workers.Get().(*worker[I, O])
	p.count.Add(1)

	go work.run()
}

func (p *pool[I, O]) down(work *worker[I, O]) bool {
	if p.count.Load() <= p.min {
		return false
	}

	p.workers.Put(work)
	p.count.Add(-1)

	return true
}
