package sail

import (
	"context"
	"sync"
	"time"
)

const (
	MaxWorkers   int32 = 1_000
	MinWorkers   int32 = 1
	ChannelSize        = 1_000
	IdleDuration       = time.Second
	BusyDuration       = time.Millisecond * 100
)

type PoolFunc[I, O any] func(context.Context, I) O

type build[I, O any] struct {
	// nolint: containedctx
	ctx    context.Context
	yield  PoolFunc[I, O]
	output chan<- O
	max    int32
	min    int32
	size   int
	idle   time.Duration
	busy   time.Duration
}

// nolint: revive
func New[I, O any](ctx context.Context, yield PoolFunc[I, O]) *build[I, O] {
	return &build[I, O]{
		ctx:   ctx,
		yield: yield,
		max:   MaxWorkers,
		min:   MinWorkers,
		size:  ChannelSize,
		idle:  IdleDuration,
		busy:  BusyDuration,
	}
}

func (p *build[I, O]) PoolFunc(yield PoolFunc[I, O]) *build[I, O] {
	p.yield = yield

	return p
}

func (p *build[I, O]) Context(ctx context.Context) *build[I, O] {
	p.ctx = ctx

	return p
}

func (p *build[I, O]) Busy(busy time.Duration) *build[I, O] {
	p.busy = busy

	return p
}

func (p *build[I, O]) Idle(idle time.Duration) *build[I, O] {
	p.idle = idle

	return p
}

func (p *build[I, O]) Output(output chan<- O) *build[I, O] {
	p.output = output

	return p
}

func (p *build[I, O]) ChannelSize(size int) *build[I, O] {
	if size < 0 {
		size = ChannelSize
	}

	p.size = size

	return p
}

func (p *build[I, O]) MinWorkers(min uint32) *build[I, O] {
	num := int32(min)
	if num < MinWorkers {
		num = MinWorkers
	}

	p.min = num

	return p
}

func (p *build[I, O]) MaxWorkers(max uint32) *build[I, O] {
	num := int32(max)
	if num < 1 {
		num = MaxWorkers
	}

	p.max = num

	return p
}

func (p *build[I, O]) PoolByID(poolID any) *pool[I, O] {
	if p.max < p.min {
		p.max = p.min
	}

	ret := &pool[I, O]{
		ctx:    context.WithValue(p.ctx, PoolID, poolID),
		yield:  p.yield,
		max:    p.max,
		min:    p.min,
		output: p.output,
		idle:   p.idle,
		busy:   p.busy,
		input:  make(chan I, p.size),
	}

	ret.workers = &sync.Pool{New: ret.newWorker}

	for i := 0; i < int(p.min); i++ {
		ret.up()
	}

	return ret
}

func (p *build[I, O]) Pool() *pool[I, O] {
	return p.PoolByID(time.Now().UnixMilli())
}
