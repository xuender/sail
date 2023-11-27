// nolint: dupword
package sail_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/xuender/sail"
)

func itoa(_ context.Context, num int) string {
	time.Sleep(time.Millisecond)

	return "num:" + strconv.Itoa(num)
}

func Example_poolPost() {
	pool := sail.New(itoa).
		ChannelSize(1).
		MaxWorkers(2).
		Busy(time.Microsecond * 200).
		Idle(time.Microsecond * 200).
		Pool()

	defer pool.Close()

	fmt.Println(pool.Process([]int{1, 2, 3, 4, 5, 6, 7, 8}))

	fmt.Println(pool.Len())

	// Output:
	// [num:1 num:2 num:3 num:4 num:5 num:6 num:7 num:8] <nil>
	// 0
}
