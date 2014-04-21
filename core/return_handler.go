package space

import (
	"reflect"

	mw "github.com/futurespaceio/ware"
)

func defaultReturnHandler() mw.ReturnHandler {
	return func(ctx mw.Context, vals []reflect.Value) {
		// TODO
	}
}
