package sail

import "sync"

type payload[I, O any] struct {
	wait   *sync.WaitGroup
	input  I
	output O
}
