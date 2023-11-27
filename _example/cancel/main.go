package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/xuender/sail"
)

func main() {
	output := make(chan string)
	defer close(output)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	pool := sail.New(ctx, func(ctx context.Context, num int) string {
		time.Sleep(time.Second)

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
		}
	}()

	for i := 0; i < 100; i++ {
		pool.Post(i)
	}

	log.Println("stop", "workers", pool.Workers())
	time.Sleep(time.Second)
	log.Println("end", "workers", pool.Workers())
}
