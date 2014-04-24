package drafts

import "github.com/futurespaceio/space/core"

type handler interface{}

func Handle() handler {
	return func(f *space.File) {
		if f.Metadata.Draft == true {
			f.Status(200)
		}
	}
}
