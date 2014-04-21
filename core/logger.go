package space

import (
	"log"
	"time"

	mw "github.com/futurespaceio/ware"
)

func Logger() mw.Handler {
	return func(c mw.Context, log *log.Logger) {
		start := time.Now()

		c.Next()

		log.Printf("Completed %v \n", time.Since(start))
	}
}
