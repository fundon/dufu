package space

import (
	"github.com/codegangsta/inject"
	W "github.com/futurespaceio/ware"
)

type Processor struct {
	inject.Injector
	*W.Ware
	index int
}

func (p *Processor) Handle(fs Filesystem) {
	for i, file := range fs.Files() {
		// Override ware.Run() method
		file.Read()
		c := p.CreateContext()
		c.Out(file)
		c.Map(file)
		c.MapTo(fs, (*Filesystem)(nil))
		c.Run()
		p.index = i
	}
}

func NewProcessor() *Processor {
	w := W.New()
	p := &Processor{inject.New(), w, 0}
	p.Map(w)
	return p
}
