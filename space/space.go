package space

import (
	"path/filepath"
	"regexp"

	"github.com/codegangsta/inject"
	mw "github.com/futurespace/ware"
)

var reHtml = regexp.MustCompile("(index)?.html$")

var cache map[string]string

type Space struct {
	inject.Injector

	Site      *Site
	Processor *Processor
	fs        Filesystem
	dir       string
}

func New() *Space {
	s := &Space{inject.New(), &Site{}, NewProcessor(), NewFilesystem(), "."}
	return s
}

// Sets a working directory.
func (s *Space) Dir(path string) (err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	s.dir = path
	return
}

// Joins the path and defaults path by key.
func (s *Space) joinPath(path string, defaults string) string {
	if path == "" {
		path = defaults
	}
	return s.Join(path)
}

func (s *Space) Source() string {
	return s.joinPath(s.Site.Source, "src")
}

func (s *Space) Destination() string {
	return s.joinPath(s.Site.Destination, "build")
}

func (s *Space) Join(path ...string) string {
	p := make([]string, len(path)+1)
	p[0] = s.dir
	copy(p[1:], path)
	return filepath.Join(p...)
}

func (s *Space) Paths(path string) FileInfos {
	return s.fs.Walk(path)
}

type ClassicSpace struct {
	*Space
	*mw.Ware
}

func Classic() *ClassicSpace {
	// Create a Ware.
	w := mw.New()
	s := New()
	p := s.Processor
	s.Map(w)
	w.MapTo(s.fs, (*Filesystem)(nil))
	w.Use(Logger())
	w.Use(generate(s))
	w.Use(build(s))
	w.Action(p.Handle)
	return &ClassicSpace{s, w}
}

func generate(s *Space) mw.Handler {
	return func(c mw.Context) {
		c.Next()
		files := s.fs.Files()
		pages := s.Site.Pages
		posts := s.Site.Posts
		for id, _ := range pages {
			files[cache[id]].Write()
		}
		for id, _ := range posts {
			files[cache[id]].Write()
		}
	}
}

func build(s *Space) mw.Handler {
	return func(c mw.Context) {
		source := s.Source()
		fileInfos := s.Paths(source)
		for path, fileInfo := range fileInfos {
			s.fs.Add(fileInfo, path, source)
		}

		c.Next()

		cache = make(map[string]string, 0)
		pages := make(Pages, 0)
		posts := make(Pages, 0)
		s.Site.Pages = pages
		s.Site.Posts = posts

		for k, file := range s.fs.Files() {
			page := file.Page
			page.Url = page.Permalink
			page.Id = filepath.FromSlash(reHtml.ReplaceAllString(page.Permalink, ""))
			page.Target.Rel = page.Permalink
			page.Target.Abs = filepath.Join(s.Destination(), page.Target.Rel)
			if page.Type == "page" {
				pages[page.Id] = page
			} else {
				posts[page.Id] = page
			}
			cache[page.Id] = k
		}
	}
}
