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
)

func FrontMatterParser(r *bufio.Reader) (contents, metedata *bytes.Buffer, err error) {
	var (
		// check frontmatter
		// 0: no fm
		// 1: fm start
		// 2: fm end
		status int
		buf    []byte
	)

	metedata = new(bytes.Buffer)
	contents = new(bytes.Buffer)

	for {
		// read line by line
		buf, err = r.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
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

			if status > 0 {
				metedata.Write(buf)
			}
		} else {
			contents.Write(buf)
		}
	}

	if status == 0 {
		contents.WriteTo(metedata)
		contents = metedata
		metedata = nil
	}

	return contents, metedata, nil
}
