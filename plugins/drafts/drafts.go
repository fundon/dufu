package drafts

import "github.com/futurespaceio/space/core"

type handler interface{}

func Drafts() handler {
	return func(f *space.File) {
		if f.Metadata.Drafts == true {
			f.Status(200)
		}
	}
}
