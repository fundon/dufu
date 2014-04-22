package drafts

import "github.com/futurespaceio/space/core"

type handler interface{}

func Drafts() handler {
	return func(f *space.File) {
		metas := f.Metadata
		if metas != nil && metas.Drafts == true {
			f.Status(200)
		}
	}
}
