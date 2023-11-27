# sail

[![Action][action-svg]][action-url]
[![Report Card][goreport-svg]][goreport-url]
[![Lines of code][lines-svg]][lines-url]
[![godoc][godoc-svg]][godoc-url]
[![License][license-svg]][license-url]

✨ **`github.com/xuender/sail` is a goroutine pool for Go.**

- dynamically expanding goroutine;
- idle goroutine auto released;

## 🚀 Install

```shell
go get github.com/xuender/sail
```

## 💡 Usage

```golang
package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/xuender/sail"
)

func main() {
	wait := sync.WaitGroup{}

	output := make(chan string)
	defer close(output)

	pool := sail.New(context.Background(), func(_ context.Context, num int) string {
		return strconv.Itoa(num)
	}).
		Output(output).
		Pool()
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
```

## 👤 Contributors

![Contributors][contributors-svg]

## 📝 License

© ender, 2023~time.Now

[MIT LICENSE][license-url]

[action-url]: https://github.com/xuender/sail/actions
[action-svg]: https://github.com/xuender/sail/workflows/Go/badge.svg

[goreport-url]: https://goreportcard.com/report/sail
[goreport-svg]: https://goreportcard.com/badge/sail

[godoc-url]: https://godoc.org/sail
[godoc-svg]: https://godoc.org/sail?status.svg

[license-url]: https://github.com/xuender/sail/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg

[contributors-svg]: https://contrib.rocks/image?repo=sail

[lines-svg]: https://sloc.xyz/sail
[lines-url]: https://github.com/boyter/scc
