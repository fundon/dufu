package permalinks

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/futurespaceio/space/core"
)

const (
	TIME_LAYOUT = time.RFC3339
	// http://jekyll.com/docs/permalinks/
	DATE   = "/:categories/:year/:month/:day/:title.html"
	PRETTY = "/:categories/:year/:month/:day/:title/"
	NONE   = "/:categories/:title.html"
	HTML   = ".html"
	INDEX  = "index" + HTML
)

type store map[string]string

// match `/cate0/cate1/2013-01-02-title.md`
var MATCHER = regexp.MustCompile(`^(.+\/)*(\d+-\d+-\d+)-(.*)(\.[^.]+)$`)

// match `/:year/:month/`
var PATTERN_REGEXP = regexp.MustCompile(`:(\w+)`)

var SLUG_REGEXP = regexp.MustCompile(`-{2,}`)

func Handle(ags ...string) interface{} {
	var (
		pattern, key string
		patterns     = store{}
		now          = time.Now()
	)
	patterns["none"] = NONE
	patterns["date"] = DATE
	patterns["pretty"] = PRETTY
	if len(ags) == 0 {
		key = "none"
	} else {
		key = ags[0]
	}
	pattern = patterns[key]
	if pattern == "" {
		pattern = key
	}
	patternMatchs := PATTERN_REGEXP.FindAllString(pattern, -1)

	return func(f *space.File) {

		var (
			permalink          = f.Metadata.Permalink
			dateStr            = f.Metadata.Date
			title              = strings.ToLower(f.Info.Name())
			ishtml             = isHTML(path.Ext(title))
			matchs, categories []string
			date               time.Time
			err                error
		)

		if dateStr == "" {
			res := MATCHER.FindStringSubmatch(title)
			l := len(res)
			if l == 0 {
				date = now
			} else {
				title = res[l-2]
				dateStr = res[l-3]
				if l == 5 {
					categories = strings.Split(res[l-4], "/")
				}
				date = createDate(dateStr, now)
			}
		} else {
			date, err = time.Parse(TIME_LAYOUT, dateStr)
			if err != nil {
				date = now
			}
		}

		if permalink == "" {
			permalink = pattern
			matchs = patternMatchs
		} else {
			matchs = PATTERN_REGEXP.FindAllString(permalink, -1)
		}

		if len(matchs) > 0 {
			if f.Metadata.Title == "" {
				title = basename(title)
			} else {
				title = strings.ToLower(slugify(f.Metadata.Title))
			}
			if len(f.Metadata.Categories) > 0 {
				categories = f.Metadata.Categories
			}

			placeholders := urlPlaceholders(date, title, categories)

			for _, m := range matchs {
				k := m[1:]
				permalink = strings.Replace(permalink, m, placeholders[k], -1)
			}
		}
		if strings.HasSuffix(permalink, HTML) == false {
			if ishtml {
				permalink = permalink[:len(permalink)-1] + HTML
			} else {
				permalink = filepath.Join(permalink, INDEX)
			}
		}

		f.Path = filepath.Clean(permalink)
	}
}

func urlPlaceholders(t time.Time, title string, categories []string) store {
	var (
		year  = t.Year()
		month = t.Month()
		day   = t.Day()
		ph    = store{}
	)
	ph["year"] = fmt.Sprintf("%d", year)
	ph["month"] = fmt.Sprintf("%02d", month)
	ph["i_month"] = fmt.Sprintf("%d", month)
	ph["day"] = fmt.Sprintf("%02d", day)
	ph["i_day"] = fmt.Sprintf("%d", day)
	ph["short_month"] = t.Format("Jan")
	ph["y_day"] = fmt.Sprintf("%d", t.YearDay())
	ph["title"] = title
	ph["categories"] = path.Join(categories...)
	return ph
}

func titleizedSlug(slug string) string {
	return strings.Title(strings.Join(strings.Split(slug, "-"), " "))
}

func slugify(title string) string {
	title = strings.Replace(title, " ", "-", -1)
	title = SLUG_REGEXP.ReplaceAllString(title, "-")
	return title
}

func basename(p string) string {
	return strings.Replace(p, path.Ext(p), "", -1)
}

func isHTML(p string) bool {
	return p == ".html"
}

func createDate(str string, now time.Time) time.Time {
	arr := strings.Split(str, "-")
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])
	return now.AddDate(year-now.Year(), month-int(now.Month()), day-now.Day())
}
