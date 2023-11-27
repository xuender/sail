package sail_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuender/sail"
)

func Test_ChannelSize(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pool := sail.New(context.Background(), itoa).
		ChannelSize(3).
		Pool()

	ass.Equal(3, pool.Cap())
	pool.Close()

	pool = sail.New(context.Background(), itoa).
		ChannelSize(-1).
		Pool()

	ass.Equal(sail.ChannelSize, pool.Cap())
	pool.Close()
}

func Test_MaxWorkers(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pool := sail.New(context.Background(), itoa).
		MaxWorkers(3).
		Pool()

	ass.Equal(int32(3), pool.MaxWorkers())
	pool.Close()

	pool = sail.New(context.Background(), itoa).
		MaxWorkers(0).
		Pool()

	ass.Equal(sail.MaxWorkers, pool.MaxWorkers())
	pool.Close()
}

func Test_MinWorkers(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	pool := sail.New(context.Background(), itoa).
		MinWorkers(3).
		Pool()

	ass.Equal(int32(3), pool.MinWorkers())
	pool.Close()

	pool = sail.New(context.Background(), itoa).
		MinWorkers(0).
		Pool()

	ass.Equal(sail.MinWorkers, pool.MinWorkers())
	pool.Close()
}

func Test_MinMax(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)

	pool := sail.New(context.Background(), itoa).
		MinWorkers(3).
		MaxWorkers(2).
		Pool()
	defer pool.Close()

	ass.Equal(int32(3), pool.MinWorkers())
	ass.Equal(int32(3), pool.MaxWorkers())
}
