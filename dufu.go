package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/futurespace/dufu/cmd/build"
)

const APP_VER = "0.0.0"

var app = cli.NewApp()

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Name = "Dufu"
	app.Usage = "A fast, pluggable static site generator"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "build",
			ShortName: "b",
			Usage:     "Build your site",
			Action:    build.Action,
			Flags: []cli.Flag{
				cli.StringFlag{"source, s", "src", "Source directory (defaults to ./src)"},
				cli.StringFlag{"destination, d", "build", "Destination directory (defaults to ./build)"},
				cli.StringFlag{"config, c", "", "Custom configuration file (defaults to config.yml|toml|json)"},
			},
		},
	}
}

func main() {
	app.Run(os.Args)
}
