package build

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
	"github.com/futurespace/dufu/plugins/drafts"
	"github.com/futurespace/dufu/plugins/markdown"
	"github.com/futurespace/dufu/plugins/permalinks"
	"github.com/futurespace/dufu/plugins/template"
	"github.com/futurespace/dufu/space"
	"gopkg.in/yaml.v1"

	mw "github.com/futurespace/ware"
)

func Action(c *cli.Context) {
	s := space.Classic()

	config, err := checkConfigFile(c.String("config"))
	var cwd string
	if err == nil {
		s.Site = config
		cwd = path.Dir(c.String("config"))
	} else {
		cwd, _ = os.Getwd()
		if c.String("source") != "" {
			s.Site.Source = c.String("source")
		}
		if c.String("destination") != "" {
			s.Site.Destination = c.String("destination")
		}
	}
	s.Dir(cwd)
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
		layout := f.Page.Layout
		if layout == "" {
			layout = "default"
		}
		locals := make(space.Locals)
		locals["Site"] = s.Site
		locals["Page"] = f.Page
		r.HTML(0, layout, locals)
	})
	s.Run()
}

func checkConfigFile(fpath string) (c *space.Site, err error) {
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
		exts = []string{"yaml", "yml", "toml", "json"}
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
	case "yaml", "yml":
		err = yaml.Unmarshal(bs, &c)
		return c, err
	case "toml":
		_, err = toml.Decode(string(bs), &c)
		return c, err
	case "json":
		err = json.Unmarshal(bs, &c)
		return c, err
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
