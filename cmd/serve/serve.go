package serve

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

// StaticOptions is a struct for specifying configuration options for the martini.Static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
}

type Static struct {
	Directory http.Dir
	Options   StaticOptions
}

func prepareStaticOptions(options StaticOptions) StaticOptions {
	var opt StaticOptions

	// Defaults
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}
	// Normalize the prefix if provided
	if opt.Prefix != "" {
		// Ensure we have a leading '/'
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		// Remove any trailing '/'
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	return opt
}

func (s *Static) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		return
	}
	file := req.URL.Path
	dir := s.Directory
	opt := s.Options
	// if we have a prefix, filter requests by stripping the prefix
	if opt.Prefix != "" {
		if !strings.HasPrefix(file, opt.Prefix) {
			return
		}
		file = file[len(opt.Prefix):]
		if file != "" && file[0] != '/' {
			return
		}
	}
	f, err := dir.Open(file)
	if err != nil {
		// discard the error?
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(req.URL.Path, "/") {
			http.Redirect(res, req, req.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, opt.IndexFile)
		f, err = dir.Open(file)
		if err != nil {
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return
		}
	}

	if !opt.SkipLogging {
		log.Println("[Static] Serving " + file)
	}

	// Add an Expires header to the static content
	if opt.Expires != nil {
		res.Header().Set("Expires", opt.Expires())
	}

	http.ServeContent(res, req, file, fi.ModTime(), f)
}

func Action(c *cli.Context) {
	port := c.String("port")
	directory := c.GlobalString("destination")
	log.Printf("Server's port is %s.\n", port)
	s := &Static{http.Dir(directory), prepareStaticOptions(StaticOptions{})}
	http.ListenAndServe(":"+port, s)
}
