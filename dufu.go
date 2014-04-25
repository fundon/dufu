package main

import (
	"log"
	"runtime"
	"time"

	"github.com/futurespaceio/dufu/modules/drafts"
	"github.com/futurespaceio/dufu/modules/markdown"
	"github.com/futurespaceio/dufu/modules/permalinks"
	"github.com/futurespaceio/dufu/space"
	mw "github.com/futurespaceio/ware"
)

const APP_VER = "0.0.0"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	s := space.Classic()
	s.Use(func(c mw.Context, fs space.Filesystem, log *log.Logger) {
		c.Next()
		log.Printf("Compiled %v files\n", len(fs.Files()))
	})
	// File Processor Middleware
	p := s.Processor
	p.Use(func(c mw.Context, f *space.File, log *log.Logger) {
		start := time.Now()
		log.Printf("File Started %s", f.Info.Name())
		c.Next()
		log.Printf("File Rendered %v \n", time.Since(start))
	})
	p.Use(drafts.Handle())
	p.Use(markdown.Render())
	p.Use(permalinks.Handle("pretty"))
	s.Run()
}
