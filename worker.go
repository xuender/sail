package sail

import (
	"context"
	"time"
)

type worker[I, O any] struct {
	id   uint32
	pool *pool[I, O]
}

func (p *worker[I, O]) run(ctx context.Context) {
	timer, _ := p.pool.timers.Get().(*time.Timer)
	defer p.pool.timers.Put(timer)

	timer.Reset(p.pool.idle)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			p.pool.stop(p)

			return
		case <-timer.C:
			timer.Stop()

			if p.pool.down(p) {
				return
			}

			timer.Reset(p.pool.idle)
		case item, open := <-p.pool.input:
			timer.Stop()

			if !open {
				p.pool.stop(p)

				return
			}

			item.output = p.pool.yield(context.WithValue(ctx, WorkerID, p.id), item.input)

			item.wait.Done()
			timer.Reset(p.pool.idle)
		}
	}
}
