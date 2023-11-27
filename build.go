package sail

import (
	"sync"
	"time"
)

const (
	MaxWorkers   = 1_000
	MinWorkers   = 1
	ChannelSize  = 1_000
	IdleDuration = time.Second
	BusyDuration = time.Millisecond * 100
)

type build[I, O any] struct {
	yield  func(I, int) O
	output chan<- O
	max    int32
	min    int32
	size   uint
	idle   time.Duration
	busy   time.Duration
}

// nolint: revive
func New[I, O any](yield func(I, int) O) *build[I, O] {
	return &build[I, O]{
		yield: yield,
		max:   MaxWorkers,
		min:   MinWorkers,
		size:  ChannelSize,
		idle:  IdleDuration,
		busy:  BusyDuration,
	}
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

func (p *build[I, O]) ChannelSize(size uint) *build[I, O] {
	if size < 1 {
		size = ChannelSize
	}

	p.size = size

	return p
}

func (p *build[I, O]) MinWorkers(min uint32) *build[I, O] {
	if min < MinWorkers {
		min = MinWorkers
	}

	p.min = int32(min)

	return p
}

func (p *build[I, O]) MaxWorkers(max uint32) *build[I, O] {
	if max < 1 {
		max = MaxWorkers
	}

	p.max = int32(max)

	return p
}

func (p *build[I, O]) Pool() *pool[I, O] {
	if p.max < p.min {
		p.max = p.min
	}

	ret := &pool[I, O]{
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
