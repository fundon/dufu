package space

import (
	"log"
	"time"

	W "github.com/futurespaceio/ware"
)

func Logger() W.Handler {
	return func(c W.Context, log *log.Logger) {
		start := time.Now()

		c.Next()

		log.Printf("Completed %v \n", time.Since(start))
	}
}
