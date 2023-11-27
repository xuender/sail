package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/xuender/sail"
)

func main() {
	wait := sync.WaitGroup{}

	output := make(chan string)
	defer close(output)

	pool := sail.New(func(num, _ int) string { return strconv.Itoa(num) }).Output(output).Pool()
	defer pool.Close()

	go func() {
		for str := range output {
			fmt.Println(str)
			wait.Done()
		}
	}()

	wait.Add(5)
	pool.Post(1, 2, 3, 4, 5)

	wait.Wait()
}
