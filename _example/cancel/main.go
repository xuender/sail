package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xuender/sail"
)

// nolint
func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	pool := sail.New(func(ctx context.Context, num int) string {
		time.Sleep(time.Second)

		return "num:" + strconv.Itoa(num)
	}).
		ChannelSize(1).
		MaxWorkers(10).
		Busy(time.Millisecond * 200).
		Idle(time.Millisecond * 200).
		Pool()

	for i := 0; i < 100; i++ {
		fmt.Println(pool.SingleCtx(ctx, i))
	}

	log.Println("stop", "workers", pool.Workers())
	time.Sleep(time.Second)
	log.Println("end", "workers", pool.Workers())
}
