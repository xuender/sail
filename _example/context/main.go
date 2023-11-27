package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/xuender/sail"
)

// nolint
func main() {
	count := 0
	sum := 0

	pool := sail.New(func(ctx context.Context, num int) string {
		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)

		log.Println("pool", ctx.Value(sail.PoolID), "worker", ctx.Value(sail.WorkerID), "num", num)
		count++
		sum += num

		return "num:" + strconv.Itoa(num)
	}).
		ChannelSize(1).
		MaxWorkers(10).
		Busy(time.Millisecond * 200).
		Idle(time.Millisecond * 200).
		Pool()

	for i := 0; i < 10; i++ {
		nums := make([]int, 10)
		for f := 0; f < 10; f++ {
			nums[f] = i*10 + f
		}

		log.Println(pool.Process(nums))
	}

	log.Println("workers", pool.Workers(), "count", count, "sum", sum)
	time.Sleep(time.Second)
	log.Println("end", "workers", pool.Workers())
	pool.Close()
	time.Sleep(time.Second)
	log.Println("stop", "workers", pool.Workers())
}
