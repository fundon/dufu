package space

import (
	"path/filepath"

	"github.com/codegangsta/inject"
	mw "github.com/futurespaceio/ware"
)

type Map map[string]interface{}

type Space struct {
	inject.Injector

	metadata  Map
	fs        Filesystem
	Processor *Processor
	dir       string
}

func New() *Space {
	s := &Space{inject.New(), make(Map, 0), NewFilesystem(), NewProcessor(), "."}
	return s
}

func (s *Space) SetMetadata(key string, value interface{}) {
	s.metadata[key] = value
}

func (s *Space) GetMetadata(key string) interface{} {
	value, err := s.metadata[key]
	if err == false {
		return nil
	}
	return value
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
func (s *Space) joinPath(key string, defaults string) string {
	path := s.GetMetadata(key)
	if path == nil {
		path = defaults
	}
	return s.Join(path.(string))
}

func (s *Space) Source() string {
	return s.joinPath("source", "src")
}

func (s *Space) Destination() string {
	return s.joinPath("destination", "build")
}

func (s *Space) Join(path ...string) string {
	p := make([]string, len(path)+1)
	p[0] = s.dir
	copy(p[1:], path[0:])
	return filepath.Join(p...)
}

func (s *Space) Paths(path string) []string {
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
	w.Use(build(s))
	w.Action(p.Handle)
	return &ClassicSpace{s, w}
}

func build(s *Space) mw.Handler {
	return func(c mw.Context) {
		source := s.Source()
		paths := s.Paths(source)

		for _, path := range paths {
			s.fs.Add(path, source)
		}

		c.Next()

		for _, file := range s.fs.Files() {
			path := filepath.Join(s.Destination(), file.Path)
			file.Path = path
			file.Write()
		}
	}
}
