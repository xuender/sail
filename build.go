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
	yield PoolFunc[I, O]
	max   int32
	min   int32
	size  int
	idle  time.Duration
	busy  time.Duration
}

// nolint: revive
func New[I, O any](yield PoolFunc[I, O]) *build[I, O] {
	return &build[I, O]{
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

func (p *build[I, O]) Busy(busy time.Duration) *build[I, O] {
	p.busy = busy

	return p
}

func (p *build[I, O]) Idle(idle time.Duration) *build[I, O] {
	p.idle = idle

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
		ctx:   context.WithValue(context.Background(), PoolID, poolID),
		id:    poolID,
		yield: p.yield,
		max:   p.max,
		min:   p.min,
		idle:  p.idle,
		busy:  p.busy,
		input: make(chan *payload[I, O], p.size),
	}

	ret.workers = sync.Pool{New: ret.newWorker}
	ret.payloads = sync.Pool{New: ret.newPayload}
	ret.waits = sync.Pool{New: ret.newWait}
	ret.timers = sync.Pool{New: ret.newTimer}

	for i := 0; i < int(p.min); i++ {
		ret.up(ret.ctx)
	}

	return ret
}

func (p *build[I, O]) Pool() *pool[I, O] {
	return p.PoolByID(time.Now().UnixMilli())
}
