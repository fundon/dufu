package drafts

import "github.com/futurespace/dufu/space"

type handler interface{}

func Handle() handler {
	return func(f *space.File) {
		if f.Metadata.Draft == true {
			f.Status(200)
		}
	}
}
