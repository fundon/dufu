package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/futurespaceio/dufu/plugins/drafts"
	"github.com/futurespaceio/dufu/plugins/markdown"
	"github.com/futurespaceio/dufu/plugins/permalinks"
	"github.com/futurespaceio/dufu/plugins/template"
	"github.com/futurespaceio/dufu/space"
	"gopkg.in/yaml.v1"

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
		cli.StringFlag{"config, c", "", "Custom configuration file (defaults to config.yml|toml|json)"},
	},
}

func runSpace(c *cli.Context) {
	s := space.Classic()

	config, err := checkConfigFile(c.String("config"))
	if err == nil {
		s.Metadata(config)
	} else {
		cwd, _ := os.Getwd()
		s.Dir(cwd)
		s.SetMetadata("source", c.String("source"))
		s.SetMetadata("destination", c.String("destination"))
	}
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

func checkConfigFile(fpath string) (config space.Map, err error) {
	if fpath == "" {
		fpath = "config.*"
	}
	var (
		de        = path.Ext(fpath)
		name, ext string
	)
	if len(de) > 0 {
		ext = de[1:]
		name = fpath[0 : len(fpath)-len(de)]
	}
	if name == "" || ext == "" {
		return nil, fmt.Errorf("not found")
	}

	var (
		exts = []string{"yml", "yaml", "toml", "json"}
		arr  map[string]string
	)
	arr = make(map[string]string, 0)
	if ext == "*" {
		for _, e := range exts {
			arr[e] = "config." + e
		}
	} else {
		arr[ext] = fpath
	}

	var bs []byte
	for e, a := range arr {
		ext = e
		bs, err = readFile(a)
		if err == nil {
			break
		}
	}

	switch ext {
	case "yml", "yaml":
		err = yaml.Unmarshal(bs, &config)
		return config, err
	case "toml":
		_, err = toml.Decode(string(bs), &config)
		return config, err
	case "json":
		err = json.Unmarshal(bs, &config)
		return config, err
	}
	return nil, fmt.Errorf("not found")
}

func readFile(fpath string) (bs []byte, err error) {
	reg := regexp.MustCompile(`config\.(ya?ml|toml|json)`)
	if reg.MatchString(fpath) == false {
		return nil, fmt.Errorf("not found")
	}

	bs, err = ioutil.ReadFile(fpath)
	return bs, err
}
