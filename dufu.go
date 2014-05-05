package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

const APP_VER = "0.0.0"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "Dufu"
	app.Usage = "A fast, pluggable static site generator"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		CmdBuild,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
