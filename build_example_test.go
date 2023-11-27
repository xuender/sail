package sail_test

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/xuender/sail"
)

func ExampleNew() {
	wait := sync.WaitGroup{}
	output := make(chan string)
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

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func Example_build_MinWorkers() {
	pool := sail.New(func(num, _ int) string { return strconv.Itoa(num) }).MinWorkers(5).Pool()
	fmt.Println(pool.Workers())

	// Output:
	// 5
}

func Example_build_ChannelSize() {
	pool := sail.New(func(num, _ int) string { return strconv.Itoa(num) }).ChannelSize(13).Pool()
	fmt.Println(pool.Cap())

	// Output:
	// 13
}
