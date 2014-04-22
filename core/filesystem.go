package space

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

type Filesystem interface {
	Add(string, string)
	Walk(string) []string
	Files() []*File
}

type filesystem struct {
	files []*File
}

func NewFilesystem() Filesystem {
	return &filesystem{}
}

func (fs *filesystem) Files() []*File {
	return fs.files
}

func (fs *filesystem) Add(realpath string, basepath string) {
	// relative path
	path, _ := filepath.Rel(basepath, realpath)
	// file's basename, including extname
	_, name := filepath.Split(realpath)

	fs.files = append(fs.files, &File{
		Name:     name,
		Path:     path,
		Buffer:   &bytes.Buffer{},
		realpath: realpath,
	})
}

func (fs *filesystem) Walk(path string) []string {
	var a []string
	walker := func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Println("Walker: ", "Please check source dir.")
			return nil
		}

		dot := filepath.Base(name)[0]
		isDir := fi.IsDir()
		isSkip := dot == '.'

		if isSkip && isDir {
			return filepath.SkipDir
		}

		if isSkip == isDir {
			a = append(a, name)
		}
		return nil
	}

	filepath.Walk(path, walker)
	return a
}
