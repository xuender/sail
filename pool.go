package sail

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type pool[I, O any] struct {
	// nolint: containedctx
	ctx      context.Context
	id       any
	yield    func(context.Context, I) O
	input    chan *payload[I, O]
	max      int32
	min      int32
	workers  sync.Pool
	payloads sync.Pool
	waits    sync.Pool
	timers   sync.Pool
	count    atomic.Int32
	idle     time.Duration
	busy     time.Duration
	workID   atomic.Uint32
	isClose  bool
}

func (p *pool[I, O]) post(ctx context.Context, items ...*payload[I, O]) {
	timer, _ := p.timers.Get().(*time.Timer)
	timer.Reset(p.busy)

	for idx, item := range items {
		select {
		case <-ctx.Done():
			p.isClose = true

			timer.Stop()
			p.timers.Put(timer)

			for _, other := range items[idx:] {
				other.wait.Done()
			}

			return
		case p.input <- item:
		case <-timer.C:
			p.up(ctx)
			p.input <- item
			timer.Reset(p.busy)
		}
	}

	timer.Stop()
	p.timers.Put(timer)
}

func (p *pool[I, O]) SingleCtx(ctx context.Context, elem I) (O, error) {
	if p.isClose {
		var zero O

		return zero, ErrClosed
	}

	payload, _ := p.payloads.Get().(*payload[I, O])
	wait, _ := p.waits.Get().(*sync.WaitGroup)

	wait.Add(1)
	payload.wait = wait
	payload.input = elem

	p.post(ctx, payload)

	wait.Wait()

	ret := payload.output

	p.waits.Put(wait)
	p.payloads.Put(payload)

	return ret, nil
}

func (p *pool[I, O]) Single(elem I) (O, error) {
	return p.SingleCtx(p.ctx, elem)
}

func (p *pool[I, O]) ProcessCtx(ctx context.Context, elems []I) ([]O, error) {
	if p.isClose {
		return nil, ErrClosed
	}

	length := len(elems)
	if length == 0 {
		return nil, nil
	}

	payloads := make([]*payload[I, O], length)
	wait, _ := p.waits.Get().(*sync.WaitGroup)

	wait.Add(len(elems))

	for idx, elem := range elems {
		pay, _ := p.payloads.Get().(*payload[I, O])

		pay.wait = wait
		pay.input = elem
		payloads[idx] = pay
	}

	p.post(ctx, payloads...)
	wait.Wait()
	p.waits.Put(wait)

	ret := make([]O, length)
	for idx, pay := range payloads {
		ret[idx] = pay.output
		p.payloads.Put(pay)
	}

	return ret, nil
}

func (p *pool[I, O]) Process(elems []I) ([]O, error) {
	return p.ProcessCtx(p.ctx, elems)
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
	p.isClose = true
	close(p.input)
}

func (p *pool[I, O]) MaxWorkers() int32 {
	return p.max
}

func (p *pool[I, O]) MinWorkers() int32 {
	return p.min
}

func (p *pool[I, O]) newWorker() any {
	return &worker[I, O]{
		id:   p.workID.Add(1),
		pool: p,
	}
}

func (p *pool[I, O]) newPayload() any {
	return &payload[I, O]{}
}

func (p *pool[I, O]) newWait() any {
	return &sync.WaitGroup{}
}

func (p *pool[I, O]) newTimer() any {
	return time.NewTimer(p.busy)
}

func (p *pool[I, O]) up(ctx context.Context) {
	if p.count.Load() > p.max {
		return
	}

	work, _ := p.workers.Get().(*worker[I, O])
	p.count.Add(1)

	go work.run(ctx)
}

func (p *pool[I, O]) down(work *worker[I, O]) bool {
	if p.count.Load() <= p.min {
		return false
	}

	p.stop(work)

	return true
}

func (p *pool[I, O]) stop(work *worker[I, O]) {
	p.workers.Put(work)
	p.count.Add(-1)
}
