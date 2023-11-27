// nolint: dupword
package sail_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/xuender/sail"
)

func itoa(_ context.Context, num int) string {
	time.Sleep(time.Millisecond)

	return strconv.Itoa(num)
}

func Example_poolPost() {
	wait := sync.WaitGroup{}
	output := make(chan string)
	pool := sail.New(context.Background(), itoa).
		ChannelSize(1).
		MaxWorkers(2).
		Busy(time.Microsecond * 200).
		Idle(time.Microsecond * 200).
		Output(output).
		Pool()

	defer pool.Close()
	defer close(output)

	go func() {
		for range output {
			fmt.Println("ok")
			wait.Done()
		}
	}()

	wait.Add(8)
	pool.Post(1, 2, 3, 4, 5, 6, 7, 8)

	wait.Wait()
	fmt.Println(pool.Len())

	// Output:
	// ok
	// ok
	// ok
	// ok
	// ok
	// ok
	// ok
	// ok
	// 0
}
