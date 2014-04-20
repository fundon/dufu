package space

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Name     string
	Path     string
	Buffer   *bytes.Buffer
	realpath string
	status   int
}

func (f *File) Written() bool {
	return f.status != 0
}

func (f *File) Write() (err error) {
	path, _ := filepath.Split(f.Path)
	ospath := filepath.FromSlash(path)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(f.Path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = f.Buffer.WriteTo(file)
	if err != nil {
		f.status = 200
	}
	return err
}

func (f *File) Read() (err error) {
	data, err := ioutil.ReadFile(f.realpath)
	if err != nil {
		return err
	}

	_, err = f.Buffer.ReadFrom(bytes.NewBuffer(data))
	return err
}
