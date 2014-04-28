package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/futurespaceio/dufu/plugins/drafts"
	"github.com/futurespaceio/dufu/plugins/markdown"
	"github.com/futurespaceio/dufu/plugins/permalinks"
	"github.com/futurespaceio/dufu/plugins/template"
	"github.com/futurespaceio/dufu/space"
	mw "github.com/futurespaceio/ware"
)

var CmdBuild = cli.Command{
	Name:        "build",
	Usage:       "Build your site",
	Description: ``,
	Action:      runSpace,
	Flags: []cli.Flag{
		cli.StringFlag{"source, s", "src", "Source directory (defaults to ./src)"},
		cli.StringFlag{"destination, d", "build", "Destination directory (defaults to ./build)"},
	},
}

func runSpace(c *cli.Context) {
	cwd, _ := os.Getwd()
	s := space.Classic()
	s.Dir(cwd)
	s.SetMetadata("source", c.String("source"))
	s.SetMetadata("destination", c.String("destination"))
	s.Use(func(c mw.Context, fs space.Filesystem, log *log.Logger) {
		log.SetPrefix("[dufu]")
		c.Next()
		log.Printf("Compiled %v files\n", len(fs.Files()))
	})
	// File Processor Middleware
	p := s.Processor
	p.Use(func(c mw.Context, f *space.File, log *log.Logger) {
		log.SetPrefix("[dufu]")
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
		layout := f.Metadata.Layout
		if layout == "" {
			layout = "default"
		}
		r.HTML(0, layout, f.Metadata)
	})
	s.Run()
}
