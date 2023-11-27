package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/xuender/sail"
)

func main() {
	wait := sync.WaitGroup{}
	count := 0
	sum := 0

	output := make(chan string)
	defer close(output)

	pool := sail.New(context.Background(), func(ctx context.Context, num int) string {
		time.Sleep(time.Duration(rand.Intn(10)*100) * time.Millisecond)

		log.Println("pool", ctx.Value(sail.PoolID), "worker", ctx.Value(sail.WorkerID), "num", num)
		count++
		sum += num

		return strconv.Itoa(num)
	}).
		ChannelSize(1).
		MaxWorkers(10).
		Busy(time.Millisecond * 200).
		Idle(time.Millisecond * 200).
		Output(output).
		Pool()

	go func() {
		for str := range output {
			log.Println("str", str)
			wait.Done()
		}
	}()

	wait.Add(100)
	for i := 0; i < 100; i++ {
		pool.Post(i)
	}

	wait.Wait()
	log.Println("workers", pool.Workers(), "count", count, "sum", sum)
	time.Sleep(time.Second)
	log.Println("end", "workers", pool.Workers())
	pool.Close()
	time.Sleep(time.Second)
	log.Println("stop", "workers", pool.Workers())
}
