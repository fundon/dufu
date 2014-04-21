# Space in Go

A fast, pluggable static site generator using [ware][] in Golang.



## Samples:

```go
package main

import (
	"github.com/futurespaceio/space/core"
	"github.com/futurespaceio/space/plugins/markdown"
)

func main() {
	s := space.Classic()
	// File Processor Middleware
	p := s.Processor
	p.Use(markdown.Render())
	s.Run()
}
```


## License

MIT

[ware]: https://github.com/futurespaceio/ware
