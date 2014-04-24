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
	PERTTY = "/:categories/:year/:month/:day/:title/"
	NONE   = "/:categories/:title.html"
)

// match `/cate0/cate1/2013-01-02-title.md`
var MATCHER = regexp.MustCompile(`^(.+\/)*(\d+-\d+-\d+)-(.*)(\.[^.]+)$`)

// match `/:year/:month/`
var PATTERN_REGEXP = regexp.MustCompile(`:(\w+)`)

var SLUG_REGEXP = regexp.MustCompile(`-{2,}`)

func Handle(opts ...string) interface{} {
	maps := make(map[string]string, 0)
	maps["none"] = NONE
	maps["date"] = DATE
	maps["pertty"] = PERTTY
	var pattern, key string
	if len(opts) == 0 {
		key = "none"
	} else {
		key = opts[0]
	}
	pattern = maps[key]
	if pattern == "" {
		pattern = key
	}
	patternMatchs := PATTERN_REGEXP.FindAllString(pattern, -1)
	now := time.Now()

	return func(f *space.File) {

		var (
			permalink          = f.Metadata.Permalink
			dateStr            = f.Metadata.Date
			title              = f.Info.Name()
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

		if len(matchs) != 0 {
			if f.Metadata.Title != "" {
				title = slugify(f.Metadata.Title)
			} else {
				title = basename(title)
			}
			if len(f.Metadata.Categories) != 0 {
				categories = f.Metadata.Categories
			}
			placeholders := urlPlaceholders(date, title, categories)

			for _, m := range matchs {
				k := m[1:]
				permalink = strings.Replace(permalink, m, placeholders[k], -1)
			}
		}
		noHasExt := strings.HasSuffix(permalink, ".html") == false
		if ishtml == false && noHasExt {
			permalink = filepath.Join(permalink, "index.html")
		} else if ishtml && noHasExt {
			permalink = filepath.Clean(permalink)
			permalink += ".html"
		}

		f.Path = filepath.Clean(permalink)
	}
}

func urlPlaceholders(t time.Time, title string, categories []string) map[string]string {
	opts := make(map[string]string, 0)
	var (
		year  = t.Year()
		month = t.Month()
		day   = t.Day()
	)
	opts["year"] = fmt.Sprintf("%d", year)
	opts["month"] = fmt.Sprintf("%02d", month)
	opts["i_month"] = fmt.Sprintf("%d", month)
	opts["day"] = fmt.Sprintf("%02d", day)
	opts["i_day"] = fmt.Sprintf("%d", day)
	opts["short_month"] = t.Format("Jan")
	opts["y_day"] = fmt.Sprintf("%d", t.YearDay())
	opts["title"] = title
	opts["categories"] = path.Join(categories...)
	return opts
}

func titleizedSlug(slug string) string {
	return strings.Title(strings.Join(strings.Split(slug, "-"), " "))
}

func slugify(title string) string {
	title = strings.ToLower(strings.Replace(title, " ", "-", -1))
	title = SLUG_REGEXP.ReplaceAllString(title, "-")
	return title
}

func basename(p string) string {
	p = strings.Replace(p, path.Ext(p), "", -1)
	return strings.ToLower(p)
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
