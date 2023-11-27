package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/xuender/sail"
)

func main() {
	wait := sync.WaitGroup{}

	output := make(chan string)
	defer close(output)

	pool := sail.New(func(num, _ int) string {
		time.Sleep(time.Second)

		return strconv.Itoa(num)
	}).
		ChannelSize(1).
		MaxWorkers(2).
		Busy(time.Millisecond * 100).
		Idle(time.Millisecond * 100).
		Output(output).
		Pool()
	defer pool.Close()

	go func() {
		for str := range output {
			log.Println("out", "str", str)
			wait.Done()
		}
	}()

	wait.Add(8)
	pool.Post(1, 2, 3, 4, 5, 6, 7, 8)

	wait.Wait()
	log.Println("end", "workers", pool.Workers())
}
