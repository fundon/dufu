// Front-matter: http://jekyllrb.com/docs/frontmatter/
package space

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

var (
	FRONT_MATTER = regexp.MustCompile(`---\s*`)
	EOL          = byte('\n')
)

func FrontMatterParser(r *bufio.Reader) (contents, metedata *bytes.Buffer, err error) {
	var (
		// frontmatter: opened/closed
		opened, closed int
		buf            []byte
	)

	metedata = new(bytes.Buffer)
	contents = new(bytes.Buffer)

	for {
		// read line by line
		buf, err = r.ReadBytes(EOL)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}

		if opened&closed == 0 {
			if FRONT_MATTER.Match(buf) {
				if opened == 0 {
					opened = 1
				} else {
					closed = 1
				}
			}

			if opened|closed == 1 {
				metedata.Write(buf)
			}
		} else {
			contents.Write(buf)
		}
	}

	if opened&closed == 0 {
		contents.WriteTo(metedata)
		contents = metedata
		metedata = nil
	}

	return contents, metedata, nil
}
