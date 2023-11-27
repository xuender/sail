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

	"github.com/xuender/sail"
)

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
