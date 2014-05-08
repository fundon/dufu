package space

type Site struct {
	Title       string
	Url         string
	Source      string
	Destination string
	Time        string
	Author      Locals
	Plugins     Locals
	Tags        []string
	Categories  []string
	Pages       Pages `yaml:"-"`
	Posts       Pages `yaml:"-"`
}

type Page struct {
	Title      string   `yaml:"title"`
	Date       string   `yaml:"date"`
	Layout     string   `yaml:"layout"`
	Permalink  string   `yaml:"permalink"`
	Draft      bool     `yaml:"draft"`
	Tags       []string `yaml:"tags"`
	Categories []string `yaml:"categories"`
	Next       string   `yaml:"next"`
	Previous   string   `yaml:"previous"`
	Vars       Locals   `yaml:"vars,omitempty"`

	Url          string
	Id           string
	Path         string
	Source       Path
	Target       Path
	NextPage     *Page
	PreviousPage *Page
	//Content      string
}

type Path struct {
	Rel string
	Abs string
}

type Pages map[string]*Page

type Locals map[string]interface{}
