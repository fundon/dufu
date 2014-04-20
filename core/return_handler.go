package space

import (
	"reflect"

	W "github.com/futurespaceio/ware"
)

func defaultReturnHandler() W.ReturnHandler {
	return func(ctx W.Context, vals []reflect.Value) {
		// TODO
	}
}
