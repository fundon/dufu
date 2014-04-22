package space

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

type FileInfo os.FileInfo
type FileInfos map[string]FileInfo

type Filesystem interface {
	Add(string, FileInfo, string)
	Walk(string) FileInfos
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

func (fs *filesystem) Add(realpath string, fi FileInfo, basepath string) {
	// relative path
	path, _ := filepath.Rel(basepath, realpath)

	fs.files = append(fs.files, &File{
		Path:     path,
		Buffer:   &bytes.Buffer{},
		Info:     fi,
		realpath: realpath,
	})
}

func (fs *filesystem) Walk(path string) (a FileInfos) {
	a = make(FileInfos, 0)
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
			a[name] = fi
		}
		return nil
	}

	filepath.Walk(path, walker)
	return a
}
