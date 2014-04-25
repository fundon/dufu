package markdown

import "github.com/futurespaceio/dufu/space"

type handler interface{}

func Render() handler {
	return func(f *space.File) {
		contents := f.Buffer.Bytes()
		contents = RenderMarkdown(contents, "")
		f.Buffer.Reset()
		f.Buffer.Write(contents)
	}
}
