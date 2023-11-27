package sail

import (
	"context"
	"time"
)

type worker[I, O any] struct {
	// nolint: containedctx
	ctx  context.Context
	pool *pool[I, O]
}

func (p *worker[I, O]) run() {
	timer := time.NewTimer(p.pool.idle)

	for {
		select {
		case <-p.ctx.Done():
			timer.Stop()
			p.pool.stop(p)

			return
		case <-timer.C:
			timer.Stop()

			if p.pool.down(p) {
				return
			}

			timer.Reset(p.pool.idle)
		case item, has := <-p.pool.input:
			timer.Stop()

			if !has {
				p.pool.stop(p)

				return
			}

			if p.pool.output == nil {
				p.pool.yield(p.ctx, item)
			} else {
				p.pool.output <- p.pool.yield(p.ctx, item)
			}

			timer.Reset(p.pool.idle)
		}
	}
}
