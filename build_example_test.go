package sail_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/xuender/sail"
)

func ExampleNew() {
	wait := sync.WaitGroup{}
	output := make(chan string)
	pool := sail.New(context.Background(), itoa).Output(output).Pool()

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
	pool := sail.New(context.Background(), itoa).MinWorkers(5).Pool()
	fmt.Println(pool.Workers())

	// Output:
	// 5
}

func Example_build_ChannelSize() {
	pool := sail.New(context.Background(), itoa).ChannelSize(13).Pool()
	fmt.Println(pool.Cap())

	// Output:
	// 13
}
