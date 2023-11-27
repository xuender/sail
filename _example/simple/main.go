package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/xuender/sail"
)

// nolint
func main() {
	pool := sail.New(func(_ context.Context, num int) string {
		return "num:" + strconv.Itoa(num)
	}).
		Pool()
	defer pool.Close()

	fmt.Println(pool.Process([]int{1, 2, 3, 4, 5}))

	// Output:
	// [num:1 num:2 num:3 num:4 num:5] <nil>
}
