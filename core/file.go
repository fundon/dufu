package space

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"launchpad.net/goyaml"
)

type File struct {
	Name     string
	Path     string
	Buffer   *bytes.Buffer
	realpath string
	status   int
	Metadata *Metadata
}

func (f *File) Status(i int) {
	f.status = i
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
	fh, err := os.Open(f.realpath)
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
		f.Metadata = &Metadata{}
		// parse yaml
		err = goyaml.Unmarshal(metadata.Bytes(), f.Metadata)
		if err != nil {
			return err
		}
	}

	_, err = f.Buffer.ReadFrom(contents)

	return err
}
