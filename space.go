package main

import (
	//"os"

	"runtime"

	"github.com/futurespaceio/space/core"
	"github.com/futurespaceio/space/plugins/markdown"
	//mw "github.com/futurespaceio/ware"
)

const APP_VER = "0.0.0"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	s := space.Classic()
	// File Processor Middleware
	p := s.Processor
	p.Use(markdown.Render())
	s.Run()
}
