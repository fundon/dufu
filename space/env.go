// Forked form https://github.com/go-martini/martini/blob/master/env.go
package space

import (
	"os"
)

// Envs
const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

// Env is the environment that Space is executing in. The SPACE_ENV is read on initialization to set this variable.
var Env = Dev

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("SPACE_ENV"))
}
