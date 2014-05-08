package space

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

type Files map[string]*File
type FileInfo os.FileInfo
type FileInfos map[string]FileInfo

type Filesystem interface {
	Add(FileInfo, string, string)
	Walk(string) FileInfos
	Files() Files
}

type filesystem struct {
	files map[string]*File
}

func NewFilesystem() Filesystem {
	return &filesystem{make(Files, 0)}
}

func (fs *filesystem) Files() Files {
	return fs.files
}

func (fs *filesystem) Add(fi FileInfo, realpath, basepath string) {
	// relative path
	path, _ := filepath.Rel(basepath, realpath)

	fs.files[path] = &File{
		//fs.files = append(fs.files, &File{
		Page: &Page{
			Target: Path{},
			Source: Path{
				Rel: path,
				Abs: realpath,
			},
		},
		Buffer: &bytes.Buffer{},
		Info:   fi,
		//})
	}
}

func (fs *filesystem) Walk(path string) (a FileInfos) {
	a = make(FileInfos, 0)
	path, _ = filepath.EvalSymlinks(path)
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
