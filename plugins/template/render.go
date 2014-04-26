// Package render is a middleware for Dufu that provides easy HTML template rendering.
//
// Forked from https://github.com/martini-contrib/render/blob/master/render.go
package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/futurespaceio/dufu/space"
	mw "github.com/futurespaceio/ware"
)

type Options struct {
	Directory  string
	Layout     string
	Extensions []string
	Funcs      []template.FuncMap
	Delims     Delims
}

type Delims struct {
	Left  string
	Right string
}

var helperFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
	"current": func() (string, error) {
		return "", nil
	},
	"content": func() (string, error) {
		return "", nil
	},
}

type Render interface {
	HTML(status int, name string, v interface{}, htmlOpt ...HTMLOptions)
}

type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

func Renderer(options ...Options) mw.Handler {
	opt := prepareOptions(options)
	t := compile(opt)

	return func(f *space.File, c mw.Context) {
		var tc *template.Template
		if space.Env == space.Dev {
			// recompile for easy development
			tc = compile(opt)
		} else {
			// use a clone of the initial template
			tc, _ = t.Clone()
		}
		c.MapTo(&renderer{f, tc, opt}, (*Render)(nil))
	}
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl"}
	}

	return opt
}

func compile(options Options) *template.Template {
	dir := options.Directory
	t := template.New(dir)
	t.Delims(options.Delims.Left, options.Delims.Right)
	// parse an initial template in case we don't have any
	template.Must(t.Parse("Dufu"))

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(r)
		for _, extension := range options.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (r[0 : len(r)-len(ext)])
				tmpl := t.New(filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range options.Funcs {
					tmpl.Funcs(funcs)
				}

				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}

		return nil
	})

	return t
}

type renderer struct {
	f   *space.File
	t   *template.Template
	opt Options
}

func (r *renderer) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	opt := r.prepareHTMLOptions(htmlOpt)
	if len(opt.Layout) > 0 {
		r.addYield(name, binding)
		name = opt.Layout
	}

	out, err := r.execute(name, binding)
	if err != nil {
		fmt.Println(err)
		r.f.Status(500)
		return
	}
	r.f.Buffer.Reset()
	io.Copy(r.f.Buffer, out)
}

func (r *renderer) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	return buf, r.t.ExecuteTemplate(buf, name, binding)
}

func (r *renderer) addYield(name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(name, binding)
			// return safe html here since we are rendering our own template
			return template.HTML(buf.String()), err
		},
		"current": func() (string, error) {
			return name, nil
		},
		"content": func() (template.HTML, error) {
			tmpl := r.t.New("content")
			template.Must(tmpl.Funcs(helperFuncs).Parse(r.f.Buffer.String()))
			buf, err := r.execute("content", binding)
			return template.HTML(buf.String()), err
		},
	}
	r.t.Funcs(funcs)
}

func (r *renderer) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.opt.Layout,
	}
}
