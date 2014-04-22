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
		status int
		buf    []byte
	)

	metedata = new(bytes.Buffer)
	contents = new(bytes.Buffer)

	for {
		// read line by line
		buf, err = r.ReadBytes(EOL)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		if status < 2 {
			if FRONT_MATTER.Match(buf) {
				if status == 0 {
					status = 1
				} else {
					status = 2
				}
			}

		}

		if status%3 == 0 {
			contents.Write(buf)
		} else {
			metedata.Write(buf)
			if status == 2 {
				status = 3
			}
		}
	}

	if status == 1 {
		metedata.WriteTo(contents)
		metedata = nil
	}

	return contents, metedata, nil
}
