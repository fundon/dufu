package space

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

type File struct {
	Buffer *bytes.Buffer
	Info   FileInfo
	Page   *Page
	status int
}

func (f *File) Status(i int) {
	f.status = i
}

func (f *File) Written() bool {
	return f.status != 0
}

func (f *File) Write() (err error) {
	abs := f.Page.Target.Abs
	path, _ := filepath.Split(abs)
	ospath := filepath.FromSlash(path)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, f.Buffer)
	if err != nil {
		f.status = 200
	}
	return err
}

func (f *File) Read() (err error) {
	fh, err := os.Open(f.Page.Source.Abs)
	if err != nil {
		return err
	}
	r := bufio.NewReader(fh)
	defer fh.Close()

	// parse front-matter
	contents, metadata, err := FrontMatterParser(r)

	if err != nil {
		return err
	}

	if metadata != nil {
		// parse yaml
		err = yaml.Unmarshal(metadata.Bytes(), f.Page)
		if err != nil {
			return err
		}
	}

	_, err = f.Buffer.ReadFrom(contents)
	return err
}
