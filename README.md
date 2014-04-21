# Space in Go

A fast, pluggable static site generator using [ware][] in Golang.



## Samples:

```go
package main

import (
	"github.com/futurespaceio/space/core"
	gfm "github.com/futurespaceio/space/plugins/markdown"
)

func main() {
	s := space.Classic()
	p := s.Processor
	p.Use(func(f *space.File) {
		contents := f.Buffer.Bytes()
		contents = gfm.RenderMarkdown(contents, "")
		f.Buffer.Reset()
		f.Buffer.Write(contents)
	})
	s.Run()
}
```


## License

MIT

[ware]: https://github.com/futurespaceio/ware
