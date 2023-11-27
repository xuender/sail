package sail

import (
	"time"
)

type worker[I, O any] struct {
	pool *pool[I, O]
	id   int
}

func (p *worker[I, O]) run() {
	timer := time.NewTimer(p.pool.idle)

	for {
		select {
		case <-timer.C:
			timer.Stop()

			if p.pool.down(p) {
				return
			}

			timer.Reset(p.pool.idle)
		case item, has := <-p.pool.input:
			timer.Stop()

			if !has {
				p.pool.down(p)

				return
			}

			if p.pool.output == nil {
				p.pool.yield(item, p.id)
			} else {
				p.pool.output <- p.pool.yield(item, p.id)
			}

			timer.Reset(p.pool.idle)
		}
	}
}
