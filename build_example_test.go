package sail_test

import (
	"fmt"

	"github.com/xuender/sail"
)

func ExampleNew() {
	pool := sail.New(itoa).Pool()

	defer pool.Close()

	fmt.Println(pool.Process([]int{1, 2, 3, 4, 5}))

	// Output:
	// [num:1 num:2 num:3 num:4 num:5] <nil>
}

func Example_build_MinWorkers() {
	pool := sail.New(itoa).MinWorkers(5).Pool()
	fmt.Println(pool.Workers())

	// Output:
	// 5
}

func Example_build_ChannelSize() {
	pool := sail.New(itoa).ChannelSize(13).Pool()
	fmt.Println(pool.Cap())

	// Output:
	// 13
}
