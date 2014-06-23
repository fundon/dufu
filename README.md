# Dufu (WIP) [![Gobuild Download](http://beta.gobuild.io/badge/github.com/futurespace/dufu/download.png)](http://beta.gobuild.io/github.com/futurespace/dufu)

A fast, pluggable static site generator using [ware][] in Golang.



## Usage:

### Commands

```
$ dufu help
```

#### dufu build
```
$ dufu help build
```


## Samples:

```go
package main

import (
	"log"
	"runtime"
	"time"

	"github.com/futurespace/dufu/plugins/drafts"
	"github.com/futurespace/dufu/plugins/markdown"
	"github.com/futurespace/dufu/plugins/permalinks"
	"github.com/futurespace/dufu/plugins/template"
	"github.com/futurespace/dufu/space"
	mw "github.com/futurespace/ware"
)

func main() {
	s := space.Classic()
	s.Use(func(c mw.Context, fs space.Filesystem, log *log.Logger) {
		c.Next()
		log.Printf("Compiled %v files\n", len(fs.Files()))
	})
	// File Processor Middleware
	p := s.Processor
	p.Use(func(c mw.Context, log *log.Logger, f *space.File) {
		start := time.Now()
		log.Printf("File Started %s", f.Info.Name())
		c.Next()
		log.Printf("File Rendered %v \n", time.Since(start))
	})
	p.Use(drafts.Handle())
	p.Use(markdown.Render())
	p.Use(permalinks.Handle("pretty"))
	p.Use(template.Renderer(template.Options{
		Layout: "layout",
	}))
	p.Use(func(f *space.File, r template.Render) {
		r.HTML(0, "post", f.Metadata)
	})
	s.Run()
}
```



## License

MIT

[ware]: https://github.com/futurespace/ware
